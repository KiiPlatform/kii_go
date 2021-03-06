package kii

import (
	"fmt"
	"testing"
	"time"

	"github.com/koron/go-dproxy"
)

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
	responseObj, err := author.OnboardGateway(&requestObj)
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
	responseObj2, err2 := au.GenerateEndNodeToken(*gatewayID, endNodeID, &EndNodeTokenRequest{})
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
	responseObj2, err2 := au.GenerateEndNodeToken(*gatewayID, "th.notexistThing", &EndNodeTokenRequest{})
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

	responseObj, err := au.GenerateEndNodeToken(*gatewayID, endNodeID, &EndNodeTokenRequest{})
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

	resp, err := endNodeAuthor.GetState(endNodeID)

	fmt.Printf("get state:%#v", resp)

	if p, err := dproxy.New(resp).M("Power").Bool(); err != nil || p != true {
		t.Errorf("should not fail.")
	}

	if b, err := dproxy.New(resp).M("Brightness").Int64(); err != nil || b != 81 {
		t.Errorf("should not fail.")
	}

	if c, err := dproxy.New(resp).M("Color").Int64(); err != nil || c != 255 {
		t.Errorf("should not fail.")
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

func TestRegisterAndLoginKiiUserSuccess(t *testing.T) {
	author := APIAuthor{
		Token: "",
		App:   testApp,
	}

	userName := fmt.Sprintf("user%d", time.Now().UnixNano())
	requestObj := UserRegisterRequest{
		LoginName: userName,
		Password:  "dummyPassword",
	}
	resp, err := author.RegisterKiiUser(requestObj)
	if err != nil {
		t.Errorf("register kiiuser failed. %s", err)
	}

	loginReqObj := UserLoginRequest{
		UserName: resp.LoginName,
		Password: "dummyPassword",
	}
	loginResp, err := author.LoginAsKiiUser(loginReqObj)
	if err != nil {
		t.Errorf("login as kiiuser failed. %s", err)
	}
	if len(loginResp.ID) < 1 {
		t.Errorf("got invalid response object %+v", loginResp)
	}

}

func TestRegisterKiiUserFail(t *testing.T) {
	author := APIAuthor{
		Token: "",
		App:   testApp,
	}

	requestObj := UserRegisterRequest{
		Password: "dummyPassword",
	}
	resp, err := author.RegisterKiiUser(requestObj)
	if err == nil {
		t.Errorf("should fail")
	}
	if resp != nil {
		t.Errorf("should be nil")
	}
}

func TestLoginAsKiiUserFail(t *testing.T) {
	author := APIAuthor{
		Token: "",
		App:   testApp,
	}

	loginReqObj := UserLoginRequest{
		UserName: "dummyUser",
		Password: "dummyPassword",
	}
	loginResp, err := author.LoginAsKiiUser(loginReqObj)
	if err == nil {
		t.Errorf("should fail")
	}
	if loginResp != nil {
		t.Errorf("should be nil")
	}
	if author.Token != "" {
		t.Errorf("access token should not be updated")
	}
}

func TestPostCommandSuccess(t *testing.T) {
	author, userID, err := GetLoginKiiUser()
	if err != nil {
		t.Errorf("fail to get login user")
	}

	endnodeID, err := RegisterAnEndNode(author)

	onboardRequest := OnboardByOwnerRequest{
		ThingID:       endnodeID,
		Owner:         "user:" + userID,
		ThingPassword: "dummyPass",
	}
	_, err = author.OnboardThingByOwner(onboardRequest)

	actions := []map[string]interface{}{
		{
			"turnPower": map[string]interface{}{
				"power": true,
			},
		},
	}
	request := PostCommandRequest{
		Issuer:        "user:" + userID,
		Actions:       actions,
		Schema:        "LED-schema",
		SchemaVersion: 1,
	}
	postResp, err := author.PostCommand(endnodeID, request)
	if err != nil {
		t.Errorf("fail to post command: %s", err)
	}
	if len(postResp.CommandID) < 1 {
		t.Errorf("got invalid response object %+v", postResp)
	}
	if err := author.DeleteThing(endnodeID); err != nil {
		t.Error("should not fail to delete Thing", err)
	}
}

func TestPostCommandFail(t *testing.T) {
	author := APIAuthor{
		Token: "dummyToken",
		App:   testApp,
	}
	actions := []map[string]interface{}{
		{
			"turnPower": map[string]interface{}{
				"power": true,
			},
		},
	}
	request := PostCommandRequest{
		Issuer:        "user:dummyID",
		Actions:       actions,
		Schema:        "LED-schema",
		SchemaVersion: 1,
	}
	postResp, err := author.PostCommand("dummyThing", request)
	if err == nil {
		t.Errorf("should fail")
	}
	if postResp != nil {
		t.Errorf("should fail")
	}

}

func TestUpdateCommandResultsSuccess(t *testing.T) {

	// Post command by endnode owner
	author, userID, err := GetLoginKiiUser()
	if err != nil {
		t.Errorf("fail to get login user")
	}
	endnodeID, err := RegisterAnEndNode(author)

	onboardRequest := OnboardByOwnerRequest{
		ThingID:       endnodeID,
		Owner:         "user:" + userID,
		ThingPassword: "dummyPass",
	}
	_, err = author.OnboardThingByOwner(onboardRequest)

	if err != nil {
		t.Errorf("onboard faild:%s", err)
	}

	actions := []map[string]interface{}{
		{
			"turnPower": map[string]interface{}{
				"power": true,
			},
		},
	}
	request := PostCommandRequest{
		Issuer:        "user:" + userID,
		Actions:       actions,
		Schema:        "LED-schema",
		SchemaVersion: 1,
	}
	postResp, err := author.PostCommand(endnodeID, request)
	if err != nil {
		t.Errorf("fail to post command: %s", err)
	}
	commandID := postResp.CommandID

	// Get endnode token and update command results
	gateway, gatewayID, err := GatewayOnboard()
	if err != nil {
		t.Errorf("onboard gateway fail:%s", err)
	}
	err = gateway.AddEndNode(*gatewayID, endnodeID)
	if err != nil {
		t.Errorf("gateway add endnode fail: %s", err)
	}
	endNodeTokenResp, err := gateway.GenerateEndNodeToken(*gatewayID, endnodeID, &EndNodeTokenRequest{})
	endNodeToken := endNodeTokenResp.AccessToken

	// endnode update Command results
	endnodeAuthor := APIAuthor{
		Token: endNodeToken,
		App:   testApp,
	}
	actionResults := []map[string]interface{}{
		{
			"turnPower": map[string]interface{}{
				"succeeded": false,
			},
		},
	}
	updateActionResultsRequest := UpdateCommandResultsRequest{
		ActionResults: actionResults,
	}
	err = endnodeAuthor.UpdateCommandResults(endnodeID, commandID, updateActionResultsRequest)
	if err != nil {
		t.Errorf("update command results faild: %s", err)
	}

	if err := author.DeleteThing(endnodeID); err != nil {
		t.Error("should not fail to delete Thing", err)
	}
}

func TestUpdateCommandResultsFail(t *testing.T) {
	// endnode update Command results
	endnodeAuthor := APIAuthor{
		Token: "dummyToken",
		App:   testApp,
	}
	actionResults := []map[string]interface{}{
		{
			"turnPower": map[string]interface{}{
				"succeeded": false,
			},
		},
	}
	updateActionResultsRequest := UpdateCommandResultsRequest{
		ActionResults: actionResults,
	}
	err := endnodeAuthor.UpdateCommandResults("dummyThingID", "dummyCommandID", updateActionResultsRequest)
	if err == nil {
		t.Errorf("should fail")
	}
}

func TestOnboardThingByOwnerSuccess(t *testing.T) {
	author, userID, err := GetLoginKiiUser()
	if err != nil {
		t.Errorf("fail to get login user")
	}

	endnodeID, err := RegisterAnEndNode(author)

	onboardRequest := OnboardByOwnerRequest{
		ThingID:       endnodeID,
		Owner:         "user:" + userID,
		ThingPassword: "dummyPass",
	}
	responseObj, err := author.OnboardThingByOwner(onboardRequest)
	if err != nil {
		t.Errorf("onboard by owner fail:%s", err)
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
	if err := author.DeleteThing(endnodeID); err != nil {
		t.Error("should not fail to delete Thing", err)
	}
}
func TestOnboardThingByOwnerFail(t *testing.T) {
	author := APIAuthor{
		Token: "dummyToken",
		App:   testApp,
	}

	onboardRequest := OnboardByOwnerRequest{
		ThingID:       "dummyID",
		Owner:         "user:dummyUser",
		ThingPassword: "dummyPass",
	}
	responseObj, err := author.OnboardThingByOwner(onboardRequest)
	if err == nil {
		t.Errorf("should fail")
	}
	if responseObj != nil {
		t.Errorf("should fail")
	}
}

func TestOnboardEndNodeWithGatewayIDSuccess(t *testing.T) {
	// get a login user
	author, userID, err := GetLoginKiiUser()
	if err != nil {
		t.Errorf("fail to get login user")
	}

	// onboard gateway with login user for ownership
	_, gwid, err := GatewayOnboard()
	if err != nil {
		t.Errorf("fail to onboard gateway:%s", err)
	}
	oboreq := OnboardByOwnerRequest{
		ThingID:       *gwid,
		Owner:         "user:" + userID,
		ThingPassword: "dummyPass",
	}
	_, err = author.OnboardThingByOwner(oboreq)
	if err != nil {
		t.Errorf("fail to onboard gateway by login user:%s", err)
	}

	// create an endnode
	endnodeID, err := RegisterAnEndNode(author)
	owgreq := OnboardEndnodeWithGatewayThingIDRequest{
		GatewayThingID: *gwid,
		OnboardEndnodeRequestCommon: OnboardEndnodeRequestCommon{
			EndNodeVendorThingID: endnodeID,
			EndNodePassword:      "dummyPass",
			Owner:                "user:" + userID,
		},
	}
	owgres, err := author.OnboardEndnodeWithGatewayThingID(owgreq)
	if err != nil {
		t.Errorf("onboard endnode with gateway id fail: %s ", err)
	}
	if owgres.AccessToken == "" {
		t.Errorf("should have accessToken")
	}
	if owgres.EndNodeThingID == "" {
		t.Errorf("should have endnodeThingID")
	}
}

func TestOnboardEndNodeWithGatewayIDFail(t *testing.T) {
	// get a login user
	author, userID, err := GetLoginKiiUser()
	if err != nil {
		t.Errorf("fail to get login user")
	}

	// create an endnode
	owgreq := OnboardEndnodeWithGatewayThingIDRequest{
		GatewayThingID: "dummyGatewayID",
		OnboardEndnodeRequestCommon: OnboardEndnodeRequestCommon{
			EndNodeVendorThingID: "dummyVendorThingID",
			EndNodePassword:      "dummyPass",
			Owner:                "user:" + userID,
		},
	}
	owgres, err := author.OnboardEndnodeWithGatewayThingID(owgreq)
	if err == nil {
		t.Errorf("onboard endnode with gateway id should fail ")
	}
	if owgres != nil {
		t.Errorf("should be nil ")
	}

}

func TestOnboardEndNodeWithGatewayVendorIDSuccess(t *testing.T) {
	// get a login user
	author, userID, err := GetLoginKiiUser()
	if err != nil {
		t.Errorf("fail to get login user")
	}

	// onboard gateway with login user for ownership
	_, gwid, err := GatewayOnboard()
	if err != nil {
		t.Errorf("fail to onboard gateway:%s", err)
	}
	oboreq := OnboardByOwnerRequest{
		ThingID:       *gwid,
		Owner:         "user:" + userID,
		ThingPassword: "dummyPass",
	}
	_, err = author.OnboardThingByOwner(oboreq)
	if err != nil {
		t.Errorf("fail to onboard gateway by login user:%s", err)
	}

	// create an endnode
	endnodeID, err := RegisterAnEndNode(author)
	owgreq := OnboardEndnodeWithGatewayVendorThingIDRequest{
		GatewayVendorThingID: "gatewayID",
		OnboardEndnodeRequestCommon: OnboardEndnodeRequestCommon{
			EndNodeVendorThingID: endnodeID,
			EndNodePassword:      "dummyPass",
			Owner:                "user:" + userID,
		},
	}
	owgres, err := author.OnboardEndnodeWithGatewayVendorThingID(owgreq)
	if err != nil {
		t.Errorf("onboard endnode with gateway id fail: %s ", err)
	}
	if owgres.AccessToken == "" {
		t.Errorf("should have accessToken")
	}
	if owgres.EndNodeThingID == "" {
		t.Errorf("should have endnodeThingID")
	}
}

func TestListEndnodeSuccess(t *testing.T) {
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

	lr, err := au.ListEndNodes(*gatewayID, ListRequest{})
	if err != nil {
		t.Errorf("got error when list endnode %s", err)
	}
	if len(lr.Results) < 1 {
		t.Errorf("results should be more than 1")
	}

	// register another endnode
	endNodeID, err = RegisterAnEndNode(au)
	if err != nil {
		t.Errorf("got error when register an end node %s", err)
	}

	err = au.AddEndNode(*gatewayID, endNodeID)
	if err != nil {
		t.Errorf("got error when add end node %s", err)
	}
	lr, err = au.ListEndNodes(*gatewayID, ListRequest{BestEffortLimit: 1})
	if err != nil {
		t.Errorf("got error when list endnode %s", err)
	}
	if len(lr.Results) != 1 {
		t.Errorf("results should be 1")
	}
	if lr.NextPaginationKey == "" {
		t.Errorf("nextPaginationKey should not be empty")
	}

	// request with nextPaginationKey
	lr, err = au.ListEndNodes(*gatewayID, ListRequest{NextPaginationKey: lr.NextPaginationKey})
	if err != nil {
		t.Errorf("got error when list endnode %s", err)
	}
	if len(lr.Results) < 1 {
		t.Errorf("results should be greater than 1")
	}

}

func TestListEndnodeFail(t *testing.T) {
	// dummy gateway
	gwAuthor := APIAuthor{
		Token: "dummyToken",
		App:   testApp,
	}
	lr, err := gwAuthor.ListEndNodes("dummyId", ListRequest{})
	if err == nil {
		t.Errorf("should fail")
	}
	if lr != nil {
		t.Errorf("response should be nil")
	}

}

func TestCreateThingScopeObjectSuccess(t *testing.T) {
	thingBucket := fmt.Sprintf("myBucket%d", time.Now().UnixNano())

	au, gwID, err := GatewayOnboard()
	if err != nil {
		t.Errorf("fail to onboard gateway %s", err)
	}

	object := map[string]interface{}{
		"age":     23,
		"country": "cn",
	}

	resp, err := au.CreateThingScopeObject(*gwID, thingBucket, object)
	if err != nil {
		t.Errorf("fail to create object :%s", err)
	}
	if resp == nil {
		t.Error("response is nil")
	}

	object2 := map[string]interface{}{
		"age":     25,
		"country": "us",
	}

	resp, err = au.CreateThingScopeObject(*gwID, thingBucket, object2)
	if err != nil {
		t.Errorf("fail to create object :%s", err)
	}
	if resp == nil {
		t.Error("response is nil")
	}

	listResp, err := au.ListAllThingScopeObjects(*gwID, thingBucket, ListRequest{BestEffortLimit: 1})
	if err != nil {
		t.Errorf("fail to list all object of thing scope:%s", err)
	}
	if listResp == nil {
		t.Error("listResp is nil")
	} else {
		if len(listResp.Results) != 1 {
			t.Errorf("should have 1 object :%d\n", len(listResp.Results))
		}

		listResp, err = au.ListAllThingScopeObjects(*gwID, thingBucket, ListRequest{NextPaginationKey: listResp.NextPaginationKey})
		if err != nil {
			t.Errorf("fail to list all object of thing scope:%s", err)
		}
		if listResp == nil {
			t.Error("listResp is nil")
		} else {
			if len(listResp.Results) != 1 {
				t.Errorf("should have 1 object :%d\n", len(listResp.Results))
			}
		}
	}

	// Delete the bucket
	if err := au.DeleteThingScopeBucket(*gwID, thingBucket); err != nil {
		t.Errorf("delete bucket fail:%s\n", err)
	}

}

func TestQueryObjectSuccess(t *testing.T) {
	thingBucket := fmt.Sprintf("myBucket%d", time.Now().UnixNano())

	au, gwID, err := GatewayOnboard()
	if err != nil {
		t.Errorf("fail to onboard gateway %s", err)
	}

	object := map[string]interface{}{
		"age":     23,
		"country": "cn",
	}

	_, err = au.CreateThingScopeObject(*gwID, thingBucket, object)
	if err != nil {
		t.Errorf("fail to create object :%s", err)
	}

	object2 := map[string]interface{}{
		"age":     25,
		"country": "us",
	}

	_, err = au.CreateThingScopeObject(*gwID, thingBucket, object2)
	if err != nil {
		t.Errorf("fail to create object :%s", err)
	}

	//Test QueryObjects
	cCluase := EqualsClause("country", "us")
	aClause := EqualsClause("age", 25)
	qClause := AndClause(cCluase, aClause)
	qreq := QueryObjectsRequest{
		BucketQuery: BucketQuery{
			Clause:     qClause,
			OrderBy:    "age",
			Descending: false,
		},
	}

	queryResp, err := au.QueryObjects(*gwID, thingBucket, qreq)
	if err != nil {
		t.Errorf("fail to list all object of thing scope:%s", err)
	}
	if queryResp == nil {
		t.Error("listResp is nil")
	} else {
		if len(queryResp.Results) != 1 {
			t.Errorf("should have 1 object :%d\n", len(queryResp.Results))
		}
	}

	// Delete the bucket
	if err := au.DeleteThingScopeBucket(*gwID, thingBucket); err != nil {
		t.Errorf("delete bucket fail:%s\n", err)
	}
}
func TestCreateThingScopeObjectFail(t *testing.T) {

	// dummy gateway
	au := APIAuthor{
		Token: "dummyToken",
		App:   testApp,
	}

	object := map[string]interface{}{
		"age":     23,
		"country": "cn",
	}

	resp, err := au.CreateThingScopeObject("dummyThingID", "dummyBucket", object)
	if err == nil {
		t.Error("should fail")
	}
	if resp != nil {
		t.Error("response should be nil")
	}
}

func TestListAllThingScopeObjectsFail(t *testing.T) {

	// dummy gateway
	au := APIAuthor{
		Token: "dummyToken",
		App:   testApp,
	}

	resp, err := au.ListAllThingScopeObjects("dummyID", "dummyBucket", ListRequest{})
	if err == nil {
		t.Error("should fail")
	}
	if resp != nil {
		t.Error("response should be nil")
	}

}

func TestDeleteThingScopeBucketFail(t *testing.T) {
	// dummy gateway
	au := APIAuthor{
		Token: "dummyToken",
		App:   testApp,
	}

	err := au.DeleteThingScopeBucket("dummyID", "dummyBucket")
	if err == nil {
		t.Error("should fail")
	}
	code := err.(*CloudError).ErrorCode
	httpCode := err.(*CloudError).HTTPStatus
	if code != "WRONG_TOKEN" {
		t.Errorf("unexpected error object: %v+", err)
	}
	if httpCode != 403 {
		t.Errorf("unexpected error object: %v+", err)
	}

}

func TestQueryObjectsFail(t *testing.T) {

	// dummy gateway
	au := APIAuthor{
		Token: "dummyToken",
		App:   testApp,
	}

	//Test QueryObjects
	cCluase := EqualsClause("country", "us")
	aClause := EqualsClause("age", 25)
	qClause := AndClause(cCluase, aClause)
	qreq := QueryObjectsRequest{
		BucketQuery: BucketQuery{
			Clause:     qClause,
			OrderBy:    "age",
			Descending: false,
		},
	}

	queryResp, err := au.QueryObjects("dummyID", "dummyBucket", qreq)

	if err == nil {
		t.Error("should fail")
	}
	if queryResp != nil {
		t.Error("response should be nil")
	}

}

func TestUpdateVendorThingIDSuccess(t *testing.T) {
	author, userID, err := GetLoginKiiUser()
	if err != nil {
		t.Error("fail to get login user", err)
	}

	newVid := fmt.Sprintf("newVID%d", time.Now().UnixNano())
	thingID, err := RegisterAnEndNode(author)
	if err != nil {
		t.Error("should not fail to register thing", err)
	}

	or := OnboardByOwnerRequest{
		ThingID:       thingID,
		ThingPassword: "dummyPass",
		Owner:         "user:" + userID,
	}
	// onboard to get ownership
	_, err = author.OnboardThingByOwner(or)
	if err != nil {
		t.Error("fail to onboard", err)
	}

	request := UpdateVendorThingIDRequest{
		VendorThingID: newVid,
		Password:      "newPass",
	}
	err = author.UpdateVendorThingID(thingID, request)
	if err != nil {
		t.Error("should not fail to update vendorThingID", err)
	}

	resp, err := author.GetThing(thingID)

	if err != nil {
		t.Error("should not fail to get Thing", err)
	} else {

		if vid, err := dproxy.New(resp).M("_vendorThingID").String(); err != nil || vid != newVid {
			t.Error("vendorThingID not updated correctly")
		}
	}

	if err := author.DeleteThing(thingID); err != nil {
		t.Error("should not fail to delete Thing", err)
	}
}

func TestUpdateVendorThingIDFail(t *testing.T) {

	au := APIAuthor{
		Token: "dummyToken",
		App:   testApp,
	}
	err := au.UpdateVendorThingID("dummyID", UpdateVendorThingIDRequest{
		VendorThingID: "newVendorThingiD",
		Password:      "newPass",
	})
	if err == nil {
		t.Error("should fail")
	}
}

func TestGetThingFail(t *testing.T) {
	au := APIAuthor{
		Token: "dummyToken",
		App:   testApp,
	}
	resp, err := au.GetThing("dummyID")
	if err == nil {
		t.Error("should fail")
	}
	if resp != nil {
		t.Error("response should be nil")
	}
}

func TestDeleteThingFail(t *testing.T) {
	au := APIAuthor{
		Token: "dummyToken",
		App:   testApp,
	}
	err := au.DeleteThing("dummyID")
	if err == nil {
		t.Error("should fail")
	}
}
func TestResetThingPasswordSuccess(t *testing.T) {
	author, err := AdminLogin(testApp, clientID, clientSecret)
	if err != nil {
		t.Errorf("anonymouseLogin fail:%s", err)
	}
	vid := fmt.Sprintf("dummyID%d", time.Now().UnixNano())
	req := RegisterThingRequest{
		VendorThingID:  vid,
		ThingPassword:  "dummyPass",
		ThingType:      "dummyType",
		LayoutPosition: ENDNODE.String(),
	}
	thing, err := author.RegisterThing(req)
	if err != nil {
		t.Errorf("failed to register thing: %v", err)
	}

	err = author.ResetThingPassword(thing.ThingID, "newPass")
	if err != nil {
		t.Errorf("failed to reset password: %v", err)
	}

	err = author.DeleteThing(thing.ThingID)
	if err != nil {
		t.Errorf("failed to delete thing: %v", err)
	}
}

func TestResetThingPasswordFail(t *testing.T) {
	au := APIAuthor{
		Token: "dummyToken",
		App:   testApp,
	}
	err := au.ResetThingPassword("dummyID", "newPass")
	if err == nil {
		t.Error("should fail")
	}
}
