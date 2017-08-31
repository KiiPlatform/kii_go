package kii

import (
	"fmt"
	"os"
	"strings"
	"time"
)

var testApp App

func init() {
	envName := "KIIGO_APP"
	s, ok := os.LookupEnv(envName)
	if !ok || s == "" {
		fmt.Printf("failed to lookup environment value for %s", envName)
		return
	}
	ss := strings.SplitN(s, ":", 3)
	if len(ss) != 3 {
		fmt.Printf("invalid format of %s, it should be {SITE}:{APP_ID}:{APP_KEY}", envName)
		return
	}
	testApp = App{
		Location: ss[0],
		AppID:    ss[1],
		AppKey:   ss[2],
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

func OnboardAGateway(vid, password string) (gateway *APIAuthor, gatewayID *string, error error) {

	author, err := AnonymousLogin(testApp)
	if err != nil {
		return nil, nil, err
	}
	requestObj := OnboardGatewayRequest{
		VendorThingID:  vid,
		ThingPassword:  password,
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
