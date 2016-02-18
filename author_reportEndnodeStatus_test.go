package kii

import "testing"

func TestReportEndnodeStatusSuccess(t *testing.T) {
	// get a login user
	author, userID, err := GetLoginKiiUser()
	if err != nil {
		t.Errorf("fail to get login user")
	}

	// onboard gateway with login user for ownership
	gwAu, gwid, err := GatewayOnboard()
	if err != nil {
		t.Errorf("fail to onboard gateway:%s", err)
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
	_, err = author.OnboardEndnodeWithGatewayVendorThingID(owgreq)
	if err != nil {
		t.Errorf("onboard endnode with gateway id fail: %s ", err)
	}

	// report endnode connection to true
	recReq := ReportEndnodeStatusRequest{
		Online: true,
	}
	err = gwAu.ReportEndnodeStatus(*gwid, endnodeID, recReq)
	if err != nil {
		t.Error("report endnode ")
	}

	// confirm online status
	oboreq := OnboardByOwnerRequest{
		ThingID:       endnodeID,
		Owner:         "user:" + userID,
		ThingPassword: "dummyPass",
	}
	_, err = author.OnboardThingByOwner(oboreq)

	gtResp, err := author.GetThing(endnodeID)
	if err != nil {
		t.Error("get thing fail", err)
	}
	if gtResp.Online != true {
		t.Error("online status of thing should be true")
	}

	// report endnode connection to false
	recReq.Online = false
	err = gwAu.ReportEndnodeStatus(*gwid, endnodeID, recReq)
	if err != nil {
		t.Error("report endnode ")
	}
	// confirm online status
	gtResp, err = author.GetThing(endnodeID)
	if err != nil {
		t.Error("get thing fail", err)
	}
	if gtResp.Online != false {
		t.Error("online status of thing should be true")
	}

	// delete the endnode
	err = author.DeleteThing(endnodeID)
	if err != nil {
		t.Error("fail to delete endnode after test", err)
	}
}

func TestReportEndnodeStatusFail(t *testing.T) {
	au := APIAuthor{
		Token: "dummyToken",
		App:   testApp,
	}
	err := au.ReportEndnodeStatus("dummyGWID", "dummyEnID", ReportEndnodeStatusRequest{Online: true})
	if err == nil {
		t.Error("should fail")
	}
}
