package tapogo

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"
)

// KlapEncryptionSession representa una sesión de cifrado y su estado interno.
type KlapEncryptionSession struct {
	localSeed  []byte
	remoteSeed []byte
	userHash   []byte
	key        []byte
	iv         []byte
	seq        int32
	sig        []byte
}

// NewKlapEncryptionSession crea una nueva instancia de KlapEncryptionSession.
func NewKlapEncryptionSession(localSeed, remoteSeed, userHash string) *KlapEncryptionSession {
	session := &KlapEncryptionSession{
		localSeed:  []byte(localSeed),
		remoteSeed: []byte(remoteSeed),
		userHash:   []byte(userHash),
	}

	session.key = session.keyDerive()
	session.iv, session.seq = session.ivDerive()
	session.sig = session.sigDerive()

	return session
}

func (s *KlapEncryptionSession) keyDerive() []byte {
	payload := append([]byte("lsk"), append(append(s.localSeed, s.remoteSeed...), s.userHash...)...)
	hash := sha256.Sum256(payload)
	return hash[:16]
}

func (s *KlapEncryptionSession) ivDerive() ([]byte, int32) {
	payload := append([]byte("iv"), append(append(s.localSeed, s.remoteSeed...), s.userHash...)...)
	fullIV := sha256.Sum256(payload)
	seq := int32(binary.BigEndian.Uint32(fullIV[12:]))
	return fullIV[:12], seq
}

func (s *KlapEncryptionSession) sigDerive() []byte {
	payload := append([]byte("ldk"), append(append(s.localSeed, s.remoteSeed...), s.userHash...)...)
	hash := sha256.Sum256(payload)
	return hash[:28]
}

func (s *KlapEncryptionSession) ivSeq() []byte {
	seqBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(seqBytes, uint32(s.seq))
	return append(s.iv, seqBytes...)
}

func (s *KlapEncryptionSession) encrypt(msg string) ([]byte, int32) {
	s.seq++
	msgBytes := []byte(msg)

	block, err := aes.NewCipher(s.key)
	if err != nil {
		fmt.Println("Error al crear el cifrador AES:", err)
		return nil, 0
	}

	cbc := cipher.NewCBCEncrypter(block, s.ivSeq())
	paddedData := PKCS7Padding(msgBytes)
	ciphertext := make([]byte, len(paddedData))
	cbc.CryptBlocks(ciphertext, paddedData)

	hash := sha256.New()
	hash.Write(append(append(s.sig, seqToBytes(s.seq)...), ciphertext...))
	signature := hash.Sum(nil)

	return append(signature, ciphertext...), s.seq
}

func (s *KlapEncryptionSession) decrypt(msg []byte) string {
	block, err := aes.NewCipher(s.key)
	if err != nil {
		fmt.Println("Error al crear el cifrador AES:", err)
		return ""
	}

	cbc := cipher.NewCBCDecrypter(block, s.ivSeq())
	plaintext := make([]byte, len(msg)-32)
	cbc.CryptBlocks(plaintext, msg[32:])

	unpaddedData, err := PKCS7Unpadding(plaintext)
	if err != nil {
		fmt.Println("Error al desempaquetar PKCS7:", err)
		return ""
	}

	return string(unpaddedData)
}

func seqToBytes(seq int32) []byte {
	seqBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(seqBytes, uint32(seq))
	return seqBytes
}

func PKCS7Padding(data []byte) []byte {
	padding := aes.BlockSize - (len(data) % aes.BlockSize)
	padText := strings.Repeat(string(padding), padding)
	return append(data, []byte(padText)...)
}

func PKCS7Unpadding(data []byte) ([]byte, error) {
	length := len(data)
	unpadding := int(data[length-1])
	if unpadding > length {
		return nil, fmt.Errorf("Invalid PKCS7 padding")
	}
	return data[:(length - unpadding)], nil
}

// usageExample is a function just to expose how to use Tapogo's KLAP library
func usageExample() {
	localSeed := "local_seed"
	remoteSeed := "remote_seed"
	userHash := "user_hash"

	session := NewKlapEncryptionSession(localSeed, remoteSeed, userHash)

	message := "Hello, world!"

	encryptedMsg, seq := session.encrypt(message)
	fmt.Printf("Encrypted message: %s, Sequence number: %d\n", hex.EncodeToString(encryptedMsg), seq)

	decryptedMsg := session.decrypt(encryptedMsg)
	fmt.Printf("Decrypted message: %s\n", decryptedMsg)
}
