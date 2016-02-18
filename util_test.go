package kii

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/koron/go-dproxy"
)

var testApp App

func init() {
	testApp = App{
		AppID:    "9ab34d8b",
		AppKey:   "7a950d78956ed39f3b0815f0f001b43b",
		Location: "JP",
	}
	// If you want to make log enabled, uncomment below line.
	//Logger = log.New(os.Stderr, "", log.LstdFlags)
}

func GatewayOnboard() (gateway *APIAuthor, gatewayID *string, error error) {

	author, err := AnonymousLogin(testApp)
	if err != nil {
		return nil, nil, err
	}
	requestObj := OnboardGatewayRequest{
		VendorThingID:  "gatewayID",
		ThingPassword:  "dummyPass",
		ThingType:      "dummyType",
		LayoutPosition: GATEWAY.String(),
		ThingProperties: map[string]interface{}{
			"myCustomString": "str",
			"myNumber":       1,
			"myObject": map[string]interface{}{
				"a": "b",
			},
		},
	}
	respObj, err := author.OnboardGateway(&requestObj)
	if err != nil {
		return nil, nil, err
	}
	author.Token = respObj.AccessToken
	return author, &respObj.ThingID, nil
}

func RegisterAnEndNode(author *APIAuthor) (endNodeID string, error error) {

	VendorThingID := fmt.Sprintf("dummyID%d", time.Now().UnixNano())
	requestObj := RegisterThingRequest{
		VendorThingID:  VendorThingID,
		ThingPassword:  "dummyPass",
		ThingType:      "dummyType",
		LayoutPosition: ENDNODE.String(),
	}
	responseObj, err := author.RegisterThing(requestObj)
	if err != nil {
		return "", err
	}
	return responseObj.ThingID, nil
}

func GetLoginKiiUser() (loginAuthor *APIAuthor, userID string, error error) {
	author := APIAuthor{
		Token: "",
		App:   testApp,
	}

	// userName := fmt.Sprintf("user%d", time.Now().UnixNano())
	userName := "user4KiiGoTest"
	password := "dummyPassword"

	// login or register a user
	loginReqObj := UserLoginRequest{
		UserName: userName,
		Password: password,
	}
	respObj, err := author.LoginAsKiiUser(loginReqObj)
	if err != nil {
		var v interface{}

		if err := json.Unmarshal([]byte(err.Error()), &v); err != nil {
			return nil, "", err
		}
		fmt.Println("ok here")
		errCode, err := dproxy.New(v).M("errorCode").String()
		if err != nil {
			return nil, "", err
		}

		if errCode != "invalid_grant" {
			return nil, "", err
		}

		requestObj := UserRegisterRequest{
			LoginName: userName,
			Password:  password,
		}
		_, err = author.RegisterKiiUser(requestObj)
		if err != nil {
			return nil, "", err
		}
		// login again to get token
		respObj, err = author.LoginAsKiiUser(loginReqObj)
		if err != nil {
			return nil, "", err
		}
	}

	author.Token = respObj.AccessToken
	return &author, respObj.ID, nil
}
