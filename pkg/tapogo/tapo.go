package tapogo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/achetronic/tapogo/api/types"
)

type TapoOptions struct {
	// HandshakeDelay represents number of seconds to wait after a handshake operation is done
	// Higher amounts are more reliable as device is incredible slow performing authorization internally
	HandshakeDelayDuration time.Duration
}

type Tapo struct {

	// Fields related to user's information
	ip       net.IP
	email    string
	password string

	//
	options *TapoOptions

	// Fields related to the device
	httpClient    *http.Client
	handshakeData *HandshakeData
}

func NewTapo(ip, email, password string, options *TapoOptions) (*Tapo, error) {
	d := &Tapo{
		ip:       net.ParseIP(ip),
		email:    email,
		password: password,

		httpClient: &http.Client{
			Timeout: 2 * time.Second,
		},

		options: options,
	}

	// start session
	if err := d.Handshake(); err != nil {
		return d, err
	}

	return d, nil
}

// PerformRequest TODO
func (d *Tapo) PerformRequest(request *types.RequestSpec) (response *types.ResponseSpec, err error) {

	jsonBytes, err := json.Marshal(*request)
	if err != nil {
		return response, err
	}

	//log.Printf("Request: %s", string(jsonBytes)) TODO: Show in debug mode only

	// 'KLAP' forces to encrypt the entire message, not only some parts inside as done in the past by 'securePassthrough'
	encryptedPayload, seq := d.handshakeData.Session.encrypt(string(jsonBytes))

	// Endpoint for KLAP is a bit different
	u, err := url.Parse(fmt.Sprintf("http://%s/app/request?seq=%d", d.ip, seq))
	if err != nil {
		return response, err
	}

	httpRequest, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(encryptedPayload))
	if err != nil {
		return response, err
	}

	httpRequest.Header.Set("Content-Type", "application/json")

	// Forward existing cookies from handshake
	for _, cookie := range d.handshakeData.Cookies {
		httpRequest.AddCookie(cookie)
	}

	httpResponse, err := d.httpClient.Do(httpRequest)
	if err != nil {
		return response, err
	}

	defer httpResponse.Body.Close()
	httpResponseBody, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return response, err
	}

	// Check request status
	if httpResponse.StatusCode != 200 {
		return response, errors.New(fmt.Sprintf("request exited with failed status: %d", httpResponse.StatusCode))
	}

	// Decrypt the payload to process it
	httpResponseBodyString := d.handshakeData.Session.decrypt(httpResponseBody)

	//log.Printf("Response (from device): %s", httpResponseBodyString) TODO: Show in debug mode only
	err = json.Unmarshal([]byte(httpResponseBodyString), &response)

	return response, err
}

// TurnOn TODO
func (d *Tapo) TurnOn() (response *types.ResponseSpec, err error) {
	request := types.RequestSpec{
		Method:          "set_device_info",
		RequestTimeMils: int(time.Now().Unix()),
		Params:          &types.RequestParamsSpec{DeviceOn: true},
	}

	response, err = d.PerformRequest(&request)
	if err != nil {
		return response, err
	}

	return response, nil
}

// TurnOff TODO
func (d *Tapo) TurnOff() (response *types.ResponseSpec, err error) {
	request := types.RequestSpec{
		Method:          "set_device_info",
		RequestTimeMils: int(time.Now().Unix()),
		Params:          &types.RequestParamsSpec{DeviceOn: false},
	}

	response, err = d.PerformRequest(&request)
	if err != nil {
		return response, err
	}

	return response, nil
}

// GetEnergyUsage TODO
func (d *Tapo) GetEnergyUsage() (response *types.ResponseSpec, err error) {
	request := types.RequestSpec{
		Method:          "get_energy_usage",
		RequestTimeMils: int(time.Now().Unix()),
	}

	response, err = d.PerformRequest(&request)
	if err != nil {
		return response, err
	}

	return response, nil
}

// DeviceInfo TODO
func (d *Tapo) DeviceInfo() (response *types.ResponseSpec, err error) {
	request := types.RequestSpec{
		Method:          "get_device_info",
		RequestTimeMils: int(time.Now().Unix()),
	}

	response, err = d.PerformRequest(&request)
	if err != nil {
		return response, err
	}

	return response, nil
}
