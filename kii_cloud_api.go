package kii

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type App struct {
	AppID       string
	AppKey      string
	AppLocation string
}

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

func (ka *App) ThingIFBaseUrl() string {
	return fmt.Sprintf("https://%s/thing-if/apps/%s", ka.HostName(), ka.AppID)
}

func (ka *App) KiiCloudBaseUrl() string {
	return fmt.Sprintf("https://%s/api/apps/%s", ka.HostName(), ka.AppID)
}

type LayoutPosition int
const (
	ENDNODE LayoutPosition = iota
	STANDALONE
	GATEWAY
)
func (lp LayoutPosition) String() string {
	switch lp {
	case ENDNODE:
		return "ENDNOE"
	case STANDALONE:
		return "STANDALONE"
	case GATEWAY:
		return "GATEWAY"
	default:
		log.Fatal("never reache here")
		return "invalid layout"
	}
}

type OnboardGatewayRequest struct {
	VendorThingID   string                 `json:"vendorThingID"`
	ThingPassword   string                 `json:"thingPassword"`
	ThingType       string                 `json:"thingType"`
	LayoutPosition  string                 `json:"layoutPosition"`
	ThingProperties map[string]interface{} `json:"thingProperties"`
}

type OnboardGatewayResponse struct {
	ThingID      string       `json:"thingID"`
	AccessToken  string       `json:"accessToken"`
	MqttEndpoint MqttEndpoint `json:"mqttEndpoint"`
}

type MqttEndpoint struct {
	InstallationID string `json:"installationID"`
	Host           string `json:"host"`
	MqttTopic      string `json:"mqttTopic"`
	Username       string `json:"userName"`
	Password       string `json:"password"`
	PortSSL        int    `json:"portSSL"`
	PortTCP        int    `json:"portTCP"`
}

type APIAuthor struct {
	Token string
	App   App
}

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
	req, err2 := http.NewRequest("POST", url, bytes.NewBuffer(reqJson))
	if err2 != nil {
		return err2
	}
	req.Header.Set("content-type", "application/json")

	client := &http.Client{}
	resp, err3 := client.Do(req)
	if err3 != nil {
		return err3
	}
	defer resp.Body.Close()

	bodyStr, err4 := ioutil.ReadAll(resp.Body)
	if err4 != nil {
		return err4
	}
	log.Println("body: " + string(bodyStr))

	var respObj AnonymousLoginResponse
	err5 := json.Unmarshal(bodyStr, &respObj)
	if err5 != nil {
		return err5
	}
	au.Token = respObj.Access_token
	return nil
}

func (au *APIAuthor) OnboardGateway(request *OnboardGatewayRequest) (OnboardGatewayResponse, error) {
	var ret OnboardGatewayResponse
	reqJson, err := json.Marshal(request)
	if err != nil {
		return ret, err
	}
	url := fmt.Sprintf("%s/onboardings", au.App.ThingIFBaseUrl())
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqJson))
	if err != nil {
		return ret, err
	}
	req.Header.Set("content-type", "application/vnd.kii.onboardingWithVendorThingIDByThing+json")
	req.Header.Set("authorization", "bearer " + au.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ret, err
	}
	defer resp.Body.Close()

	bodyStr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ret, err
	}
	log.Println("body: " + string(bodyStr))

	err = json.Unmarshal(bodyStr, &ret)
	if err != nil {
		return ret, err
	}
	return ret, nil
}
