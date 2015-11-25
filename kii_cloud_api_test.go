package kii_test

import (
	"fmt"
	kii "github.com/KiiPlatform/kii_go"
	"testing"
	"time"
)

var testApp kii.App

func init() {
	testApp = kii.App{
		AppID:       "9ab34d8b",
		AppKey:      "7a950d78956ed39f3b0815f0f001b43b",
		AppLocation: "JP",
	}
}

func TestAnonymousLogin(t *testing.T) {
	author := kii.APIAuthor{
		App: testApp,
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
	author := kii.APIAuthor{
		App: testApp,
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
func TestGenerateEndNodeTokenSuccess(t *testing.T) {
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
func TestGenerateEndNodeTokenFail(t *testing.T) {
	gateway, err := GatewayOnboard()
	if err != nil {
		t.Errorf("got error on onboard gateway %s", err)
	}
	responseObj2, err2 := gateway.GenerateEndNodeToken("th.notexistThing")
	if err2 == nil {
		t.Errorf("should fail")
	}

	if responseObj2 != nil {
		t.Errorf("should fail")
	}
}

func TestRegisterEndNodeSuccess(t *testing.T) {
	VendorThingID := fmt.Sprintf("dummyID%d", time.Now().UnixNano())
	requestObj := kii.ThingRegisterRequest{
		VendorThingID:  VendorThingID,
		ThingPassword:  "dummyPass",
		ThingType:      "dummyType",
		LayoutPosition: kii.ENDNODE.String(),
		ThingProperties: map[string]interface{}{
			"myCustomString": "str",
			"myNumber":       1,
			"myObject": map[string]interface{}{
				"a": "b",
			},
		},
	}
	responseObj, err := testApp.RegisterThing(requestObj)
	if err != nil {
		t.Errorf("fail to register thing")
	}
	if len(responseObj.ThingID) < 1 {
		t.Errorf("got invalid response object %+v", responseObj)
	}

	if len(responseObj.VendorThingID) < 1 {
		t.Errorf("got invalid response object %+v", responseObj)
	}

	if len(responseObj.ThingType) < 1 {
		t.Errorf("got invalid response object %+v", responseObj)
	}

	if len(responseObj.LayoutPosition) < 1 {
		t.Errorf("got invalid response object %+v", responseObj)
	}
}

func TestRegisterEndNodeFail(t *testing.T) {
	requestObj := kii.ThingRegisterRequest{
		VendorThingID:  "",
		ThingPassword:  "dummyPass",
		ThingType:      "dummyType",
		LayoutPosition: kii.ENDNODE.String(),
		ThingProperties: map[string]interface{}{
			"myCustomString": "str",
			"myNumber":       1,
			"myObject": map[string]interface{}{
				"a": "b",
			},
		},
	}
	responseObj, err := testApp.RegisterThing(requestObj)
	if err == nil {
		t.Errorf("should fail")
	}
	if responseObj != nil {
		t.Errorf("should fail")
	}
}

func RegisterAnEndNode() (endNodeID *string, error error) {
	VendorThingID := fmt.Sprintf("dummyID%d", time.Now().UnixNano())
	requestObj := kii.ThingRegisterRequest{
		VendorThingID:  VendorThingID,
		ThingPassword:  "dummyPass",
		ThingType:      "dummyType",
		LayoutPosition: kii.ENDNODE.String(),
		ThingProperties: map[string]interface{}{
			"myCustomString": "str",
			"myNumber":       1,
			"myObject": map[string]interface{}{
				"a": "b",
			},
		},
	}
	responseObj, err := testApp.RegisterThing(requestObj)
	if err != nil {
		return nil, err
	} else {
		return &responseObj.ThingID, nil
	}
}
func TestAddEndNodeSuccess(t *testing.T) {

	gateway, err := GatewayOnboard()
	if err != nil {
		t.Errorf("got error on onboard gateway %s", err)
	}
	endNodeID, err := RegisterAnEndNode()
	if err != nil {
		t.Errorf("got error when register an end node %s", err)
	}

	err = gateway.AddEndNode(*endNodeID)
	if err != nil {
		t.Errorf("got error when add end node %s", err)
	}
}

func TestAddEndNodeFail(t *testing.T) {

	gateway, err := GatewayOnboard()
	if err != nil {
		t.Errorf("got error on onboard gateway %s", err)
	}
	err = gateway.AddEndNode("dummyEndNode")
	if err == nil {
		t.Errorf("should fail")
	}
}
