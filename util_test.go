package kii

import (
	"fmt"
	"time"
)

var testApp App

func init() {
	testApp = App{
		AppID:    "crju493ckopg",
		AppKey:   "408d090d161c40d2b24ae289030351df",
		Location: "US",
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
