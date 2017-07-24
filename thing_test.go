package kii

import (
	"testing"

	dproxy "github.com/koron/go-dproxy"
)

func TestUpdateThing(t *testing.T) {
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
	endNodeAuthor := APIAuthor{
		Token: responseObj.AccessToken,
		App:   testApp,
	}
	ps := map[string]interface{}{
		"StateUploadDisabled": false,
	}
	err = endNodeAuthor.UpdateThing(endNodeID, ps)
	if err != nil {
		t.Errorf("failed to update thing, %s", err)
	}
	thing, err := endNodeAuthor.GetThing(endNodeID)
	if d, err := dproxy.New(thing).M("StateUploadDisabled").Bool(); err != nil || d != false {
		t.Errorf("properties not update correctly")
	}

	if err = endNodeAuthor.DeleteThing(endNodeID); err != nil {
		t.Errorf("delete thing failed, %s", err)
	}
}

func TestUpdateThingFail(t *testing.T) {
	au := APIAuthor{
		Token: "dummyToken",
		App:   testApp,
	}
	err := au.UpdateThing("dummyID", map[string]interface{}{"Field1": "value"})
	if err == nil {
		t.Error("should fail")
	}
}
