package kii_test

import (
	kii "github.com/KiiPlatform/kii_go"
	"testing"
)

func TestAnonymousLogin(t *testing.T) {
	app := kii.App{
		AppID:       "9ab34d8b",
		AppKey:      "7a950d78956ed39f3b0815f0f001b43b",
		AppLocation: "JP",
	}
	author := kii.APIAuthor{
		App: app,
	}
	err := author.AnonymousLogin()
	if err != nil {
		t.Errorf("got error on anonymous login %s", err)
	}
	if len(author.Token) < 1 {
		t.Errorf("failed to get author token %+v", author)
	}
	if len(author.ID) < 1 {
		t.Errorf("failed to get author ID %+v", author)
	}
}

func AnonymousLogin() (kii.APIAuthor, error) {
	app := kii.App{
		AppID:       "9ab34d8b",
		AppKey:      "7a950d78956ed39f3b0815f0f001b43b",
		AppLocation: "JP",
	}
	author := kii.APIAuthor{
		App: app,
	}
	err := author.AnonymousLogin()
	if err != nil {
		return author, err
	}
	return author, nil
}

func TestGatewayOnboard(t *testing.T) {
	author, err := AnonymousLogin()
	if err != nil {
		t.Errorf("got error on anonymous login %s", err)
	}
	requestObj := kii.OnboardGatewayRequest{
		VendorThingID:  "dummyID",
		ThingPassword:  "dummyPass",
		ThingType:      "dummyType",
		LayoutPosition: kii.GATEWAY.String(),
		ThingProperties: map[string]interface{}{
			"myCustomString": "str",
			"myNumber":       1,
			"myObject": map[string]interface{}{
				"a": "b",
			},
		},
	}
	responseObj, err := author.OnboardGateway(requestObj)
	if err != nil {
		t.Errorf("got error on Onboarding %s", err)
	}
	if len(responseObj.ThingID) < 1 {
		t.Errorf("got invalid response object %+v", responseObj)
	}
	if len(responseObj.AccessToken) < 1 {
		t.Errorf("got invalid response object %+v", responseObj)
	}
	if len(responseObj.MqttEndpoint.InstallationID) < 1 {
		t.Errorf("got invalid endpoint object %+v", responseObj.MqttEndpoint)
	}
	if len(responseObj.MqttEndpoint.Host) < 1 {
		t.Errorf("got invalid endpoint object %+v", responseObj.MqttEndpoint)
	}
	if len(responseObj.MqttEndpoint.MqttTopic) < 1 {
		t.Errorf("got invalid endpoint object %+v", responseObj.MqttEndpoint)
	}
	if len(responseObj.MqttEndpoint.Username) < 1 {
		t.Errorf("got invalid endpoint object %+v", responseObj.MqttEndpoint)
	}
	if len(responseObj.MqttEndpoint.Password) < 1 {
		t.Errorf("got invalid endpoint object %+v", responseObj.MqttEndpoint)
	}
	if responseObj.MqttEndpoint.PortSSL < 1 {
		t.Errorf("got invalid endpoint object %+v", responseObj.MqttEndpoint)
	}
	if responseObj.MqttEndpoint.PortTCP < 1 {
		t.Errorf("got invalid endpoint object %+v", responseObj.MqttEndpoint)
	}
}

func GatewayOnboard() (*kii.Gateway, error) {

	author, err := AnonymousLogin()
	if err != nil {
		return nil, err
	}
	requestObj := kii.OnboardGatewayRequest{
		VendorThingID:  "dummyEndNodeID",
		ThingPassword:  "dummyPass",
		ThingType:      "dummyType",
		LayoutPosition: kii.GATEWAY.String(),
		ThingProperties: map[string]interface{}{
			"myCustomString": "str",
			"myNumber":       1,
			"myObject": map[string]interface{}{
				"a": "b",
			},
		},
	}
	responseObj, err := author.OnboardGateway(requestObj)
	if err != nil {
		return nil, err
	}
	gateway := kii.Gateway{
		Token: responseObj.AccessToken,
		ID:    responseObj.ThingID,
		App:   author.App,
	}
	return &gateway, err
}
func TestGenerateEndNodeToken(t *testing.T) {
	gateway, err := GatewayOnboard()
	if err != nil {
		t.Errorf("got error on onboard gateway %s", err)
	}
	responseObj2, err2 := gateway.GenerateEndNodeToken("th.350948a00022-10ca-5e11-6829-0ffc0c06")
	if err2 != nil {
		t.Errorf("got error when GenerateEndNodeToken %s", err2)
	}
	if responseObj2.AccessToken == "" {
		t.Errorf("got response object failed")
	}
}
