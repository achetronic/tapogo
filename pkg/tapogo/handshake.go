package tapogo

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// HandshakeData represents TODO
type HandshakeData struct {
	LocalSeed          []byte // OK
	RemoteSeed         []byte // OK
	AuthHash           []byte // OK
	RemoteSeedAuthHash []byte // OK

	Cookies []*http.Cookie         // OK
	Session *KlapEncryptionSession // OK

	EncodedCredentialsLocalSeed []byte // TODO REVIEW
}

// The protocol works by doing a two stage handshake to obtain
// and encryption key and session id cookie.
//
// Authentication uses an auth_hash which is a combination
// of username and password, hashed together. The way how this
// hash is calculated has changed several times. In the following lines,
// those algorithms are shown by version as <version>: <algorithm>
// v1: md5(md5(username)md5(password))
// v2: sha256(sha1(username)sha1(password))
//
// handshake1: client sends a random 16 byte local_seed to the
// device and receives a random 16 bytes remote_seed, followed
// by sha256(local_seed + auth_hash).  It also returns a
// TP_SESSIONID in the cookie header.  This implementation WILL
// then check this value against the possible auth_hashes
// described above (user cloud, kasa hardcoded, blank).  If it
// finds a match it moves onto handshake2
//
// handshake2: client sends sha25(remote_seed + auth_hash) to
// the device along with the TP_SESSIONID.  Device responds with
// 200 if succesful.  It generally will be because this
// implemenation checks the auth_hash it recevied during handshake1
//
// encryption: local_seed, remote_seed and auth_hash are now used
// for encryption.  The last 4 bytes of the initialisation vector
// are used as a sequence number that increments every time the
// client calls encrypt and this sequence number is sent as an
// url parameter to the device along with the encrypted payload
//
// References:
// https://github.com/python-kasa/python-kasa/blob/master/kasa/klaptransport.py
// https://gist.github.com/chriswheeldon/3b17d974db3817613c69191c0480fe55
// https://github.com/insomniacslk/tapo

// GenerateAuthHashV2 TODO
func (d *Tapo) GenerateAuthHashV2() []byte {
	emailHash := sha1.New()
	passwordHash := sha1.New()

	emailHash.Write([]byte(d.email))
	emailHashBytes := emailHash.Sum(nil)

	passwordHash.Write([]byte(d.password))
	passwordHashBytes := passwordHash.Sum(nil)

	mixedHashBytes := append(emailHashBytes, passwordHashBytes...)
	finalHashBytes := sha256.Sum256(mixedHashBytes)

	return finalHashBytes[:]
}

// GenerateSeedAuthHash TODO
func (d *Tapo) GenerateSeedAuthHash(localSeed []byte, remoteSeed []byte, authHash []byte, handshakeStage int) []byte {

	var finalHashContentBytes []byte

	switch handshakeStage {
	case 1:
		finalHashContentBytes = append(localSeed, remoteSeed...)
	case 2:
		finalHashContentBytes = append(remoteSeed, localSeed...)
	}

	finalHashContentBytes = append(finalHashContentBytes, authHash...)

	finalHashBytes := sha256.Sum256(finalHashContentBytes)
	return finalHashBytes[:]
}

// Handshake1 TODO
func (d *Tapo) Handshake1() (handshakeData HandshakeData, err error) {
	handshakeData.LocalSeed = make([]byte, 16)
	handshakeData.RemoteSeed = make([]byte, 16)
	handshakeData.EncodedCredentialsLocalSeed = make([]byte, 0)

	_, err = rand.Read(handshakeData.LocalSeed)
	if err != nil {
		return handshakeData, fmt.Errorf("error while generating random string: %s", err)
	}

	u, err := url.Parse(fmt.Sprintf("http://%s/app/handshake1", d.ip))
	if err != nil {
		return handshakeData, err
	}

	bodyBytesReader := bytes.NewBuffer(handshakeData.LocalSeed)
	request, err := http.NewRequest(http.MethodPost, u.String(), bodyBytesReader)
	if err != nil {
		return handshakeData, fmt.Errorf("error creating HTTP request: %s", err)
	}

	//log.Printf("REQUEST: %v", request) TODO: show this in debug mode only

	response, err := d.httpClient.Do(request)
	if err != nil {
		return handshakeData, fmt.Errorf("error making HTTP request: %s", err)
	}

	defer response.Body.Close()

	// Check request status
	if response.StatusCode != 200 {
		return handshakeData, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	handshakeData.Cookies = response.Cookies()

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return handshakeData, fmt.Errorf("error reading response body: %s", err)
	}

	// Recover results from server
	handshakeData.RemoteSeed = bodyBytes[0:16]
	handshakeData.EncodedCredentialsLocalSeed = bodyBytes[16:]

	return handshakeData, nil
}

// Handshake2 TODO
func (d *Tapo) Handshake2(handshakeData *HandshakeData) error {
	// Generate AuthHash
	authHash := d.GenerateAuthHashV2()
	handshakeData.AuthHash = authHash
	//log.Printf("AuthHash: %x", authHash) TODO: Show in debug mode only

	// Generate SeedAuthHash
	remoteSeedAuthHash := d.GenerateSeedAuthHash(handshakeData.LocalSeed, handshakeData.RemoteSeed, authHash, 2)
	handshakeData.RemoteSeedAuthHash = remoteSeedAuthHash
	//log.Printf("SeedAuthHash: %x", remoteSeedAuthHash) TODO: Show in debug mode only

	// Create URL for Handshake2
	u, err := url.Parse(fmt.Sprintf("http://%s/app/handshake2", d.ip))
	if err != nil {
		return err
	}

	// Create HTTP request
	request, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(remoteSeedAuthHash))
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %s", err)
	}

	// Forward cookies from Handshake1 stage
	for _, cookie := range handshakeData.Cookies {
		request.AddCookie(cookie)
		// log.Printf("Value: %v", cookie.Value) TODO: show this in debug mode only
		// log.Printf("Unparsed: %v", cookie.Unparsed) TODO: show this in debug mode only
	}

	// Perform the HTTP request
	response, err := d.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("error making HTTP request: %s", err)
	}

	defer response.Body.Close()

	// Check request status
	if response.StatusCode != 200 {
		return fmt.Errorf("handshake 2 failed with status code: %d", response.StatusCode)
	}

	// Create a KLAP encryption session to manage the messages between device and this library
	handshakeData.Session = NewKlapEncryptionSession(
		string(handshakeData.LocalSeed),
		string(handshakeData.RemoteSeed),
		string(handshakeData.AuthHash))

	return nil
}

// Handshake TODO
func (d *Tapo) Handshake() error {

	// Perform first stage of handshake phase
	// The mission here is to get a remote seed and cookies
	handshakeData, err := d.Handshake1()
	if err != nil {
		return err
	}

	// Not waiting ends in failures WTF?!
	time.Sleep(d.options.HandshakeDelayDuration / 2)

	// Perform second stage of handshake phase
	// The mission here is to get a KLAP encryption session
	err = d.Handshake2(&handshakeData)
	if err != nil {
		return err
	}

	// Not waiting ends in failures WTF?!
	time.Sleep(d.options.HandshakeDelayDuration / 2)

	d.handshakeData = &handshakeData
	return nil
}
