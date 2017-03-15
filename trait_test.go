package kii

import (
	"fmt"
	"testing"
	"time"

	dproxy "github.com/koron/go-dproxy"
)

var thingType string
var firmwareVersion string
var alias string

func init() {
	thingType = "MyAirConditioner"
	firmwareVersion = "v1"
	alias = "AirConditionerAlias"
}

type AirConditonerState struct {
	Power       bool  `json:"power"`
	Temperature int64 `json:"currentTemperature"`
}

func RegisterATraitEnabledEndNode(author *APIAuthor) (endNodeID string, error error) {

	VendorThingID := fmt.Sprintf("dummyID%d", time.Now().UnixNano())
	requestObj := RegisterThingRequest{
		VendorThingID:   VendorThingID,
		ThingPassword:   "dummyPass",
		ThingType:       thingType,
		FirmwareVersion: firmwareVersion,
		LayoutPosition:  ENDNODE.String(),
	}
	responseObj, err := author.RegisterThing(requestObj)
	if err != nil {
		return "", err
	}
	return responseObj.ThingID, nil
}

func TestOnboadWithFirmwareVersion(t *testing.T) {
	// get a login user
	author, userID, err := GetLoginKiiUser()
	if err != nil {
		t.Errorf("fail to get login user")
	}

	// onboard gateway with login user for ownership
	_, gwid, err := GatewayOnboard()
	if err != nil {
		t.Errorf("got error on onboard gateway %s", err)
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

	endnodeThingID := fmt.Sprintf("dummyID%d", time.Now().UnixNano())
	owgrep := OnboardEndnodeWithGatewayThingIDRequest{
		GatewayThingID: *gwid,
		OnboardEndnodeRequestCommon: OnboardEndnodeRequestCommon{
			EndNodeVendorThingID:   endnodeThingID,
			EndNodePassword:        "dummyPass",
			Owner:                  "user:" + userID,
			EndNodeThingType:       thingType,
			EndNodeFirmwareVersion: firmwareVersion,
		},
	}
	owgres, err := author.OnboardEndnodeWithGatewayThingID(owgrep)
	if err != nil {
		t.Errorf("onboard endnode with gateway id fail: %s", err)
	}
	if owgres.AccessToken == "" {
		t.Errorf("should have accessToken")
	}
	if owgres.EndNodeThingID == "" {
		t.Errorf("should have endnodeThingID")
	}

	thingRes, err := author.GetThing(owgres.EndNodeThingID)
	if err != nil {
		t.Errorf("get thing failed:%s", err)
	}
	if thingRes.ThingType != thingType {
		t.Errorf("thingType is wrong %s", thingRes.ThingType)
	}
	if thingRes.FirmwareVersion != firmwareVersion {
		t.Errorf("firmwareVersion is wrong %s", thingRes.FirmwareVersion)
	}
}

func TestUpdateMultipleTraitStateSuccess(t *testing.T) {
	au, gatewayID, err := GatewayOnboard()
	if err != nil {
		t.Errorf("got error on onboard gateway %s", err)
	}
	endNodeID, err := RegisterATraitEnabledEndNode(au)
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

	request := map[string]interface{}{
		alias: AirConditonerState{
			Power:       true,
			Temperature: 23,
		},
	}

	endNodeAuthor := APIAuthor{
		Token: responseObj.AccessToken,
		App:   testApp,
	}
	err = endNodeAuthor.UpdateMultipleTraitState(endNodeID, request)
	if err != nil {
		t.Errorf("should not fail. %s", err)
	}

	resp, err := endNodeAuthor.GetState(endNodeID)

	fmt.Printf("get state:%#v", resp)

	if p, err := dproxy.New(resp).M(alias).M("power").Bool(); err != nil || p != true {
		t.Errorf("should not fail.")
	}

	if c, err := dproxy.New(resp).M(alias).M("currentTemperature").Int64(); err != nil || c != 23 {
		t.Errorf("should not fail.")
	}
}

func TestUpdateMultipleTraitsStateFail(t *testing.T) {
	au, gatewayID, err := GatewayOnboard()
	if err != nil {
		t.Errorf("got error on onboard gateway %s", err)
	}
	endNodeID, err := RegisterATraitEnabledEndNode(au)
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

	request := map[string]interface{}{
		"not-existing-alias": AirConditonerState{
			Power:       true,
			Temperature: 23,
		},
	}

	endNodeAuthor := APIAuthor{
		Token: responseObj.AccessToken,
		App:   testApp,
	}
	err = endNodeAuthor.UpdateMultipleTraitState(endNodeID, request)
	if err == nil {
		t.Errorf("should fail.")
	}
}

func TestUpdateSingleTraitStateSuccess(t *testing.T) {
	au, gatewayID, err := GatewayOnboard()
	if err != nil {
		t.Errorf("got error on onboard gateway %s", err)
	}
	endNodeID, err := RegisterATraitEnabledEndNode(au)
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

	endNodeAuthor := APIAuthor{
		Token: responseObj.AccessToken,
		App:   testApp,
	}
	err = endNodeAuthor.UpdateTraitState(
		endNodeID,
		alias,
		AirConditonerState{
			Power:       true,
			Temperature: 23,
		})
	if err != nil {
		t.Errorf("should not fail. %s", err)
	}

	resp, err := endNodeAuthor.GetState(endNodeID)

	fmt.Printf("get state:%#v", resp)

	if p, err := dproxy.New(resp).M(alias).M("power").Bool(); err != nil || p != true {
		t.Errorf("should not fail.")
	}

	if c, err := dproxy.New(resp).M(alias).M("currentTemperature").Int64(); err != nil || c != 23 {
		t.Errorf("should not fail.")
	}
}

func TestUpdateSingleTraitsStateFail(t *testing.T) {
	au, gatewayID, err := GatewayOnboard()
	if err != nil {
		t.Errorf("got error on onboard gateway %s", err)
	}
	endNodeID, err := RegisterATraitEnabledEndNode(au)
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

	endNodeAuthor := APIAuthor{
		Token: responseObj.AccessToken,
		App:   testApp,
	}
	err = endNodeAuthor.UpdateTraitState(
		endNodeID,
		"not-existing-alias",
		AirConditonerState{
			Power:       true,
			Temperature: 23,
		})
	if err == nil {
		t.Errorf("should fail.")
	}
}
