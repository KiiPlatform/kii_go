package kii

import (
	"fmt"
	"testing"
	"time"
)

var testApp App

func init() {
	testApp = App{
		AppID:    "9ab34d8b",
		AppKey:   "7a950d78956ed39f3b0815f0f001b43b",
		Location: "JP",
	}
}

func TestAnonymousLogin(t *testing.T) {

	author, err := AnonymousLogin(testApp)
	if err != nil {
		t.Errorf("got error on anonymous login %s", err)
	}
	if len(author.Token) < 1 {
		t.Errorf("failed to get author token %+v", author)
	}
}

func TestGatewayOnboard(t *testing.T) {
	author, err := AnonymousLogin(testApp)
	if err != nil {
		t.Errorf("got error on anonymous login %s", err)
	}

	requestObj := OnboardGatewayRequest{
		VendorThingID:  "dummyID",
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

func GatewayOnboard() (gateway *APIAuthor, gatewayID *string, error error) {

	author, err := AnonymousLogin(testApp)
	if err != nil {
		return nil, nil, err
	}
	requestObj := OnboardGatewayRequest{
		VendorThingID:  "dummyEndNodeID",
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
	respObj, err := author.OnboardGateway(requestObj)
	if err != nil {
		return nil, nil, err
	}
	author.Token = respObj.AccessToken
	return author, &respObj.ThingID, nil
}

func TestGenerateEndNodeTokenSuccess(t *testing.T) {
	au, gatewayID, err := GatewayOnboard()
	if err != nil {
		t.Errorf("got error on onboard gateway %s", err)
	}
	endNodeID, err := RegisterAnEndNode(au)
	if err != nil {
		t.Errorf("got error when register an end node %s", err)
	}

	err = au.AddEndNode(*gatewayID, endNodeID)
	if err != nil {
		t.Errorf("got error when add end node %s", err)
	}
	responseObj2, err2 := au.GenerateEndNodeToken(*gatewayID, endNodeID, EndNodeTokenRequest{})
	if err2 != nil {
		t.Errorf("got error when GenerateEndNodeToken %s", err2)
	}
	if responseObj2.AccessToken == "" {
		t.Errorf("got response object failed")
	}
}
func TestGenerateEndNodeTokenFail(t *testing.T) {
	au, gatewayID, err := GatewayOnboard()
	if err != nil {
		t.Errorf("got error on onboard gateway %s", err)
	}
	responseObj2, err2 := au.GenerateEndNodeToken(*gatewayID, "th.notexistThing", EndNodeTokenRequest{})
	if err2 == nil {
		t.Errorf("should fail")
	}

	if responseObj2 != nil {
		t.Errorf("should fail")
	}
}

func TestRegisterEndNodeSuccess(t *testing.T) {
	author, err := AnonymousLogin(testApp)
	if err != nil {
		t.Errorf("anonymouseLogin fail:%s", err)
	}

	VendorThingID := fmt.Sprintf("dummyID%d", time.Now().UnixNano())
	type MyRegisterThingRequest struct {
		RegisterThingRequest
		MyCustomString string                 `json:"myCustomString"`
		MyNumber       int                    `json:"myNumber"`
		MyObject       map[string]interface{} `json:"myObject"`
	}
	requestObj := MyRegisterThingRequest{

		RegisterThingRequest: RegisterThingRequest{
			VendorThingID:  VendorThingID,
			ThingPassword:  "dummyPass",
			ThingType:      "dummyType",
			LayoutPosition: ENDNODE.String(),
		},
		MyCustomString: "str",
		MyNumber:       1,
		MyObject: map[string]interface{}{
			"a": "b",
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
	author, err := AnonymousLogin(testApp)
	if err != nil {
		t.Errorf("anonymouseLogin fail:%s", err)
	}

	requestObj := RegisterThingRequest{
		VendorThingID:  "",
		ThingPassword:  "dummyPass",
		ThingType:      "dummyType",
		LayoutPosition: ENDNODE.String(),
	}
	responseObj, err := author.RegisterThing(requestObj)
	if err == nil {
		t.Errorf("should fail")
	}
	if responseObj != nil {
		t.Errorf("should fail")
	}
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
func TestAddEndNodeSuccess(t *testing.T) {
	author, gatewayID, err := GatewayOnboard()
	if err != nil {
		t.Errorf("got error on onboard gateway %s", err)
	}
	endNodeID, err := RegisterAnEndNode(author)
	if err != nil {
		t.Errorf("got error when register an end node %s", err)
	}

	err = author.AddEndNode(*gatewayID, endNodeID)
	if err != nil {
		t.Errorf("got error when add end node %s", err)
	}
}

func TestAddEndNodeFail(t *testing.T) {

	author, gatewayID, err := GatewayOnboard()
	if err != nil {
		t.Errorf("got error on onboard gateway %s", err)
	}
	err = author.AddEndNode(*gatewayID, "dummyEndNode")
	if err == nil {
		t.Errorf("should fail")
	}
}

func TestEndNodeStateSuccess(t *testing.T) {
	au, gatewayID, err := GatewayOnboard()
	if err != nil {
		t.Errorf("got error on onboard gateway %s", err)
	}
	endNodeID, err := RegisterAnEndNode(au)
	if err != nil {
		t.Errorf("got error when register an end node %s", err)
	}

	err = au.AddEndNode(*gatewayID, endNodeID)
	if err != nil {
		t.Errorf("got error when add end node %s", err)
	}

	responseObj, err := au.GenerateEndNodeToken(*gatewayID, endNodeID, EndNodeTokenRequest{})
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

	endNodeAuthor := APIAuthor{
		Token: responseObj.AccessToken,
		App:   testApp,
	}
	err = endNodeAuthor.UpdateState(endNodeID, request)
	if err != nil {
		t.Errorf("should not fail. %s", err)
	}
}

func TestEndNodeStateFail(t *testing.T) {
	endNodeAuthor := APIAuthor{
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
	err := endNodeAuthor.UpdateState("dummyID", request)
	if err == nil {
		t.Errorf("should fail.")
	}
}
