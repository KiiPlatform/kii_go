// Provides APIs to access to Kii Cloud and
// Thing Interaction Framework (thing-if).
package kii

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// Represents Application in Kii Cloud.
type App struct {
	AppID       string
	AppKey      string
	AppLocation string
}

// Obtain Host name of the Application endpoint.
func (ka *App) HostName() string {
	lowerLoc := strings.ToLower(ka.AppLocation)
	switch lowerLoc {
	case "jp":
		return "api-jp.kii.com"
	case "us":
		return "api.kii.com"
	case "cn":
		return "api-cn3.kii.com"
	case "sg":
		return "api-sg.kii.com"
	default:
		return lowerLoc
	}
}

// Obtain thing-if endpoint base url.
func (ka *App) ThingIFBaseUrl() string {
	return fmt.Sprintf("https://%s/thing-if/apps/%s", ka.HostName(), ka.AppID)
}

// Obtain Kii Cloud endpoint base url.
func (ka *App) KiiCloudBaseUrl() string {
	return fmt.Sprintf("https://%s/api/apps/%s", ka.HostName(), ka.AppID)
}

// Layout position of the Thing
type LayoutPosition int

const (
	ENDNODE LayoutPosition = iota
	STANDALONE
	GATEWAY
)

// Obtain Layout postion of the Thing in string.
func (lp LayoutPosition) String() string {
	switch lp {
	case ENDNODE:
		return "END_NODE"
	case STANDALONE:
		return "STANDALONE"
	case GATEWAY:
		return "GATEWAY"
	default:
		log.Fatal("never reache here")
		return "invalid layout"
	}
}

// Struct for requesting Gateway Onboard.
type OnboardGatewayRequest struct {
	VendorThingID   string                 `json:"vendorThingID"`
	ThingPassword   string                 `json:"thingPassword"`
	ThingType       string                 `json:"thingType"`
	LayoutPosition  string                 `json:"layoutPosition"`
	ThingProperties map[string]interface{} `json:"thingProperties"`
}

// Struct for receiving response of Gateway Onboard.
type OnboardGatewayResponse struct {
	ThingID      string       `json:"thingID"`
	AccessToken  string       `json:"accessToken"`
	MqttEndpoint MqttEndpoint `json:"mqttEndpoint"`
}

// Struct represents MQTT endpoint.
type MqttEndpoint struct {
	InstallationID string `json:"installationID"`
	Host           string `json:"host"`
	MqttTopic      string `json:"mqttTopic"`
	Username       string `json:"userName"`
	Password       string `json:"password"`
	PortSSL        int    `json:"portSSL"`
	PortTCP        int    `json:"portTCP"`
}

// Struct represents API author.
type APIAuthor struct {
	Token string
	ID    string
	App   App
}

// Struct represents Gateway.
type Gateway struct {
	Token string
	ID    string
	App   App
}

// Struct for requesting end node token
type GenerateEndNodeTokenRequest struct {
	ExpiresIn string `json:"expires_in,omitempty"`
}

// Struct for receiving response of end node token
type GenerateEndNodeTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	ThingID      string `json:"id"`
	RefreshToken string `json:"refresh_token"`
}

// Struct for requesting Thing Registration.
type ThingRegisterRequest struct {
	VendorThingID   string                 `json:"_vendorThingID"`
	ThingPassword   string                 `json:"_password"`
	ThingType       string                 `json:"_thingType"`
	LayoutPosition  string                 `json:"_layoutPosition"`
	ThingProperties map[string]interface{} `json:"thingProperties"`
}

// Struct for receiving response of end node token
type ThingRegisterResponse struct {
	ThingID        string `json:"_thingID"`
	VendorThingID  string `json:"_vendorThingID"`
	ThingType      string `json:"_thingType"`
	LayoutPosition string `json:"_layoutPosition"`
	Created        int    `json:"_created"`
	Disabled       bool   `json:"_disabled"`
}

// Login as Anonymous user.
// When there's no error, Token is updated with Anonymous token.
func (au *APIAuthor) AnonymousLogin() error {
	type AnonymousLoginRequest struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		GrantType    string `json:"grant_type"`
	}
	type AnonymousLoginResponse struct {
		ID           string `json:"id"`
		Access_token string `json:"access_token"`
		Expires_in   int    `json:"expires_in"`
		Token_type   string `json:"token_type"`
	}
	reqObj := AnonymousLoginRequest{
		ClientID:     au.App.AppID,
		ClientSecret: au.App.AppKey,
		GrantType:    "client_credentials",
	}
	reqJson, err := json.Marshal(reqObj)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/oauth2/token", au.App.KiiCloudBaseUrl())
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqJson))
	if err != nil {
		return err
	}
	req.Header.Set("content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
	}

	bodyStr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Println("body: " + string(bodyStr))

	var respObj AnonymousLoginResponse
	err = json.Unmarshal(bodyStr, &respObj)
	if err != nil {
		return err
	}
	au.Token = respObj.Access_token
	au.ID = respObj.ID
	return nil
}

// Let Gateway onboard to the cloud.
// When there's no error, OnboardGatewayResponse is returned.
func (au *APIAuthor) OnboardGateway(request OnboardGatewayRequest) (*OnboardGatewayResponse, error) {
	var ret OnboardGatewayResponse
	reqJson, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/onboardings", au.App.ThingIFBaseUrl())
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/vnd.kii.onboardingWithVendorThingIDByThing+json")
	req.Header.Set("authorization", "bearer "+au.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyStr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Println("body: " + string(bodyStr))

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		err = json.Unmarshal(bodyStr, &ret)
		if err != nil {
			return nil, err
		}
		return &ret, nil
	} else {
		err = errors.New(string(bodyStr))
		return nil, err
	}
}

// Request access token of end node of gateway.
// When there's no error, GenerateEndNodeTokenResponse is returned.
func (gw *Gateway) GenerateEndNodeToken(endnodeID string) (*GenerateEndNodeTokenResponse, error) {
	var ret GenerateEndNodeTokenResponse
	url := fmt.Sprintf("%s/things/%s/end-nodes/%s/token", gw.App.KiiCloudBaseUrl(), gw.ID, endnodeID)

	reqObj := GenerateEndNodeTokenRequest{
		ExpiresIn: "",
	}
	reqJson, _ := json.Marshal(reqObj)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqJson))
	req.Header.Set("content-type", "application/json")
	req.Header.Set("authorization", "bearer "+gw.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyStr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Println("body: " + string(bodyStr))

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		err = json.Unmarshal(bodyStr, &ret)
		if err != nil {
			return nil, err
		}
		return &ret, nil
	} else {
		err = errors.New(string(bodyStr))
		return nil, err
	}

}

// Register Thing.
// Where there is no error, ThingRegisterResponse is returned
func (app *App) RegisterThing(request ThingRegisterRequest) (*ThingRegisterResponse, error) {
	var ret ThingRegisterResponse

	reqJson, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/things", app.KiiCloudBaseUrl())
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqJson))
	if err != nil {
		return nil, err
	}

	req.Header.Set("content-type", "application/vnd.kii.ThingRegistrationRequest+json")
	req.Header.Set("X-Kii-AppID", app.AppID)
	req.Header.Set("X-Kii-AppKey", app.AppKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyStr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Println("body: " + string(bodyStr))

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		err = json.Unmarshal(bodyStr, &ret)
		if err != nil {
			return nil, err
		}
		return &ret, nil
	} else {
		err = errors.New(string(bodyStr))
		return nil, err
	}
}
