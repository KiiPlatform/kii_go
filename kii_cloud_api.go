packagea thing_if_gateway

import (
        "bytes"
        "encoding/json"
        "fmt"
        "io/ioutil"
        "net/http"
        "strings"
        "log"
)

type KiiApp struct {
        AppID string
        AppKey string
        AppLocation string
}

func (ka *KiiApp) HostName() string {
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

func (ka *KiiApp) ThingIFBaseUrl() string {
        return fmt.Sprintf("https://%s/thing-if/apps/%s", ka.HostName(), ka.AppID)
}

func (ka *KiiApp) KiiCloudBaseUrl() string {
        return fmt.Sprintf("https://%s/api/apps/%s", ka.HostName(), ka.AppID)
}

type OnboardGatewayRequest struct {
        // TODO: implement marshaler for thingProperties
        VendorThingID string `json:"vendorThingID"`
        ThingPassword string `json:"thingPasssword"`
        thingType string `json:"thingType"`
        layoutPosition string `json:"layoutPosition"`
}

type OnboardGatewayResponse struct {
        ThingID string `json:"thingID"`
        AccessToken string `json:"accessToken"`
        MqttEndPoint MqttEndPoint `json:"mqttEndpooint"`
}

type MqttEndPoint struct {
        InstallationID string `json:"installationID"`
        Host string `json:"host"`
        MqttTopic string `json:"mqttTopic"`
        Username string `json:"userName"`
        Password string `json:"password"`
        PortSSL int `json:"portSSL"`
        PortTCP int `json:"portTCP"`
}

type APIAuthor struct {
        Token string
        App KiiApp
}

func (au *APIAuthor) AnonymousLogin() error {
        type AnonymousLoginRequest struct {
                ClientID string `json:"client_id"`
                ClientSecret string `json:"client_secret"`
                GrantType string `json:"grant_type"`
        }
        type AnonymousLoginResponse struct {
                ID string `json:"id"`
                Access_token string `json:"access_token"`
                Expires_in int `json:"expires_in"`
                Token_type string `json:"token_type"`
        }
        reqObj := AnonymousLoginRequest {
                ClientID: au.App.AppID,
                ClientSecret: au.App.AppKey,
                GrantType: "client_credentials",
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

func (au *APIAuthor) Onboard(request *OnboardGatewayRequest) (OnboardGatewayResponse, error) {
        // TODO: implement it.
        var ret OnboardGatewayResponse
        return ret, nil
}
