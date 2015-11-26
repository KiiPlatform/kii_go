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

func executeRequest(request http.Request) (respBody []byte, error error) {

	client := &http.Client{}
	resp, err := client.Do(&request)
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
		return bodyStr, nil
	} else {
		err = errors.New(string(bodyStr))
		return nil, err
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
// Can be Gateway, EndNode or KiiUser, depending on the token.
type APIAuthor struct {
	Token string
	App   App
}

// Struct for requesting end node token
type EndNodeTokenRequest struct {
	ExpiresIn string `json:"expires_in,omitempty"`
}

// Struct for receiving response of end node token
type EndNodeTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	ThingID      string `json:"id"`
	RefreshToken string `json:"refresh_token"`
}

// Struct for requesting Thing Registration.
type RegisterThingRequest struct {
	VendorThingID   string                 `json:"_vendorThingID"`
	ThingPassword   string                 `json:"_password"`
	ThingType       string                 `json:"_thingType,omitempty"`
	LayoutPosition  string                 `json:"_layoutPosition,omitempty"`
	Vendor          string                 `json:"_vendor,omitempty"`
	FirmwareVersion string                 `json:"_firmwareVersion,omitempty"`
	Iot             string                 `json:"_iot,omitempty"`
	StringField1    string                 `json:"_stringField1,omitempty"`
	StringField2    string                 `json:"_stringField2,omitempty"`
	StringField3    string                 `json:"_stringField3,omitempty"`
	StringField4    string                 `json:"_stringField4,omitempty"`
	StringField5    string                 `json:"_stringField5,omitempty"`
	NumberField1    int64                  `json:"_numberField1,omitempty"`
	NumberField2    int64                  `json:"_numberField2,omitempty"`
	NumberField3    int64                  `json:"_numberField3,omitempty"`
	NumberField4    int64                  `json:"_numberField4,omitempty"`
	NumberField5    int64                  `json:"_numberField5,omitempty"`
	ThingProperties map[string]interface{} `json:"thingProperties"`
}

// Struct for receiving response of end node token
type RegisterThingResponse struct {
	ThingID        string `json:"_thingID"`
	VendorThingID  string `json:"_vendorThingID"`
	ThingType      string `json:"_thingType"`
	LayoutPosition string `json:"_layoutPosition"`
	Created        int    `json:"_created"`
	Disabled       bool   `json:"_disabled"`
}

// Login as Anonymous user.
// When there's no error, Token is updated with Anonymous token.
func AnonymousLogin(app App) (*APIAuthor, error) {
	type AnonymousLoginRequest struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		GrantType    string `json:"grant_type"`
	}
	type AnonymousLoginResponse struct {
		ID          string `json:"id"`
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		TokenType   string `json:"token_type"`
	}
	reqObj := AnonymousLoginRequest{
		ClientID:     app.AppID,
		ClientSecret: app.AppKey,
		GrantType:    "client_credentials",
	}
	reqJson, err := json.Marshal(reqObj)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/oauth2/token", app.KiiCloudBaseUrl())
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
	}

	bodyStr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Println("body: " + string(bodyStr))

	var respObj AnonymousLoginResponse
	err = json.Unmarshal(bodyStr, &respObj)
	if err != nil {
		return nil, err
	}
	au := APIAuthor{
		Token: respObj.AccessToken,
		App:   app,
	}
	return &au, nil
}

// Let Gateway onboard to the cloud.
// Notes that the APIAuthor must be Anonymous user.
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
	req.Header.Set("authorization", "Bearer "+au.Token)

	bodyStr, err := executeRequest(*req)
	if err != nil {
		return nil, err
	} else {
		err = json.Unmarshal(bodyStr, &ret)
		if err != nil {
			return nil, err
		}
		return &ret, nil
	}
}

// Request access token of end node of gateway.
// Notes the APIAuthor should be a Gateway.
// When there's no error, EndNodeTokenResponse is returned.
func (au APIAuthor) GenerateEndNodeToken(gatewayID string, endnodeID string, request EndNodeTokenRequest) (*EndNodeTokenResponse, error) {
	var ret EndNodeTokenResponse
	url := fmt.Sprintf("%s/things/%s/end-nodes/%s/token", au.App.KiiCloudBaseUrl(), gatewayID, endnodeID)

	reqJson, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqJson))
	req.Header.Set("content-type", "application/json")
	req.Header.Set("authorization", "Bearer "+au.Token)

	bodyStr, err := executeRequest(*req)
	if err != nil {
		return nil, err
	} else {
		err = json.Unmarshal(bodyStr, &ret)
		if err != nil {
			return nil, err
		}
		return &ret, nil
	}
}

// Add an end node thing to gateway
// Notes that the APIAuthor should be a Gateway
// when it succeeds, error is nil
func (au APIAuthor) AddEndNode(gatewayID string, endnodeID string) error {
	url := fmt.Sprintf("%s/things/%s/end-nodes/%s", au.App.KiiCloudBaseUrl(), gatewayID, endnodeID)

	req, err := http.NewRequest("PUT", url, nil)
	req.Header.Set("content-type", "application/json")
	req.Header.Set("authorization", "Bearer "+au.Token)
	if err != nil {
		return err
	}

	_, err1 := executeRequest(*req)
	return err1
}

// Register Thing.
// Where there is no error, RegisterThingResponse is returned
func RegisterThing(app App, request RegisterThingRequest) (*RegisterThingResponse, error) {
	var ret RegisterThingResponse

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

	bodyStr, err := executeRequest(*req)
	if err != nil {
		return nil, err
	} else {
		err = json.Unmarshal(bodyStr, &ret)
		if err != nil {
			return nil, err
		}
		return &ret, nil
	}
}

// Update Thing state.
// Notes that the APIAuthor should be already initialized as a Gateway or EndNode
func (au APIAuthor) UpdateState(thingID string, request interface{}) error {

	reqJson, err := json.Marshal(request)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/targets/thing:%s/states", au.App.ThingIFBaseUrl(), thingID)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(reqJson))
	if err != nil {
		return err
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("authorization", "Bearer "+au.Token)

	_, err1 := executeRequest(*req)
	return err1
}
