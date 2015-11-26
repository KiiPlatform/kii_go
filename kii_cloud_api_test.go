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
	tokeBeforeOnboard, idBeforeOnBoard := author.Token, author.ID

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
	if tokeBeforeOnboard == author.Token {
		t.Errorf("token should be updated")
	}
	if idBeforeOnBoard == author.ID {
		t.Errorf("ID should be updated")
	}
}

func GatewayOnboard() (*kii.APIAuthor, error) {

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
	_, err1 := author.OnboardGateway(requestObj)
	if err1 != nil {
		return nil, err1
	}
	return &author, nil
}

func TestGenerateEndNodeTokenSuccess(t *testing.T) {
	au, err := GatewayOnboard()
	if err != nil {
		t.Errorf("got error on onboard gateway %s", err)
	}
	endNodeID, err := RegisterAnEndNode(au)
	if err != nil {
		t.Errorf("got error when register an end node %s", err)
	}

	err = au.AddEndNode(endNodeID)
	if err != nil {
		t.Errorf("got error when add end node %s", err)
	}
	responseObj2, err2 := au.GenerateEndNodeToken(endNodeID, kii.EndNodeTokenRequest{})
	if err2 != nil {
		t.Errorf("got error when GenerateEndNodeToken %s", err2)
	}
	if responseObj2.AccessToken == "" {
		t.Errorf("got response object failed")
	}
}
func TestGenerateEndNodeTokenFail(t *testing.T) {
	au, err := GatewayOnboard()
	if err != nil {
		t.Errorf("got error on onboard gateway %s", err)
	}
	responseObj2, err2 := au.GenerateEndNodeToken("th.notexistThing", kii.EndNodeTokenRequest{})
	if err2 == nil {
		t.Errorf("should fail")
	}

	if responseObj2 != nil {
		t.Errorf("should fail")
	}
}

func TestRegisterEndNodeSuccess(t *testing.T) {
	author, err := AnonymousLogin()
	if err != nil {
		t.Errorf("anonymouseLogin fail:%s", err)
	}

	VendorThingID := fmt.Sprintf("dummyID%d", time.Now().UnixNano())
	requestObj := kii.RegisterThingRequest{
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
	responseObj, err := author.RegisterThing(requestObj)
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
	author, err := AnonymousLogin()
	if err != nil {
		t.Errorf("anonymouseLogin fail:%s", err)
	}

	requestObj := kii.RegisterThingRequest{
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
	responseObj, err := author.RegisterThing(requestObj)
	if err == nil {
		t.Errorf("should fail")
	}
	if responseObj != nil {
		t.Errorf("should fail")
	}
}

func RegisterAnEndNode(author *kii.APIAuthor) (endNodeID string, error error) {

	VendorThingID := fmt.Sprintf("dummyID%d", time.Now().UnixNano())
	requestObj := kii.RegisterThingRequest{
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
	responseObj, err := author.RegisterThing(requestObj)
	if err != nil {
		return "", err
	} else {
		return responseObj.ThingID, nil
	}
}
func TestAddEndNodeSuccess(t *testing.T) {
	author, err := GatewayOnboard()
	if err != nil {
		t.Errorf("got error on onboard gateway %s", err)
	}
	endNodeID, err := RegisterAnEndNode(author)
	if err != nil {
		t.Errorf("got error when register an end node %s", err)
	}

	err = author.AddEndNode(endNodeID)
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

func TestEndNodeStateSuccess(t *testing.T) {
	au, err := GatewayOnboard()
	if err != nil {
		t.Errorf("got error on onboard gateway %s", err)
	}
	endNodeID, err := RegisterAnEndNode(au)
	if err != nil {
		t.Errorf("got error when register an end node %s", err)
	}

	err = au.AddEndNode(endNodeID)
	if err != nil {
		t.Errorf("got error when add end node %s", err)
	}

	responseObj, err := au.GenerateEndNodeToken(endNodeID, kii.EndNodeTokenRequest{})
	if err != nil {
		t.Errorf("got error when GenerateEndNodeToken %s", err)
	}

	type UpdateStateRequest struct {
		Power      bool
		Brightness int
		Color      int
	}

	request := UpdateStateRequest{
		Power:      true,
		Brightness: 81,
		Color:      255,
	}

	endNodeAuthor := kii.APIAuthor{
		ID:    endNodeID,
		Token: responseObj.AccessToken,
		App:   testApp,
	}
	err = endNodeAuthor.UpdateState(request)
	if err != nil {
		t.Errorf("should not fail. %s", err)
	}
}

func TestEndNodeStateFail(t *testing.T) {
	endNodeAuthor := kii.APIAuthor{
		ID:    "dummyID",
		Token: "dummyToken",
		App:   testApp,
	}

	type UpdateStateRequest struct {
		Power      bool
		Brightness int
		Color      int
	}

	request := UpdateStateRequest{
		Power:      true,
		Brightness: 81,
		Color:      255,
	}
	err := endNodeAuthor.UpdateState(request)
	if err == nil {
		t.Errorf("should fail.")
	}
}
