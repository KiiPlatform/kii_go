package kii

import (
	"fmt"
	"testing"
	"time"
	// dproxy "github.com/koron/go-dproxy"
)

func TestQueryThingsSuccess(t *testing.T) {
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
	author.Token = loginResp.AccessToken

	gwvid := fmt.Sprintf("gwID%d", time.Now().UnixNano())

	_, gwid, err := OnboardAGateway(gwvid, "dummyPass")
	if err != nil {
		t.Errorf("fail to onboard gateway:%s", err)
	}

	oboreq := OnboardByOwnerRequest{
		ThingID:       *gwid,
		Owner:         "user:" + loginResp.ID,
		ThingPassword: "dummyPass",
	}
	_, err = author.OnboardThingByOwner(oboreq)
	if err != nil {
		t.Errorf("fail to onboard gateway by login user:%s", err)
	}

	// create an endnode
	endnodeID, err := RegisterAnEndNode(&author)
	owgreq := OnboardEndnodeWithGatewayThingIDRequest{
		GatewayThingID: *gwid,
		OnboardEndnodeRequestCommon: OnboardEndnodeRequestCommon{
			EndNodeVendorThingID: endnodeID,
			EndNodePassword:      "dummyPass",
			Owner:                "user:" + loginResp.ID,
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

	// query thing by owner
	queryResp1, err := author.QueryThings(
		ThingQueryRequest{
			OwnerID: loginResp.ID,
		},
	)
	if err != nil {
		t.Errorf("should not fail to execute query: %v", err)
	}
	if len(queryResp1.Results) != 2 {
		t.Errorf("results should have 2")
	}

	// query gateway
	queryResp2, err := author.QueryThings(
		ThingQueryRequest{
			OwnerID: loginResp.ID,
			Clause:  EqualsClause("_layoutPosition", "GATEWAY"),
		},
	)
	if err != nil {
		t.Errorf("should not fail to execute query")
	}
	if len(queryResp2.Results) != 1 {
		t.Errorf("results should have 1")
	}
	if queryResp2.NextPaginationKey != "" {
		t.Errorf("nextPaginationKey should be empty")
	}

	// query endnode
	queryResp3, err := author.QueryThings(
		ThingQueryRequest{
			OwnerID: loginResp.ID,
			Clause:  EqualsClause("_layoutPosition", "GATEWAY"),
		},
	)
	if err != nil {
		t.Errorf("should not fail to execute query")
	}
	if len(queryResp3.Results) != 1 {
		t.Errorf("results should have 1")
	}
	if queryResp3.NextPaginationKey != "" {
		t.Errorf("nextPaginationKey should be empty")
	}

	// query with list options
	queryResp4, err := author.QueryThings(
		ThingQueryRequest{
			OwnerID:     loginResp.ID,
			ListRequest: ListRequest{BestEffortLimit: 1},
		},
	)
	if err != nil {
		t.Errorf("should not fail to execute query")
	}
	if len(queryResp4.Results) != 1 {
		t.Errorf("results should have 1")
	}
	if queryResp4.NextPaginationKey == "" {
		t.Errorf("nextPaginationKey should not be empty")
	}

	// query with list options
	queryResp5, err := author.QueryThings(
		ThingQueryRequest{
			OwnerID:     loginResp.ID,
			ListRequest: ListRequest{NextPaginationKey: queryResp4.NextPaginationKey},
		},
	)
	if err != nil {
		t.Errorf("should not fail to execute query")
	}
	if len(queryResp5.Results) != 1 {
		t.Errorf("results should have 1")
	}
	if queryResp5.NextPaginationKey != "" {
		t.Errorf("nextPaginationKey should be empty")
	}
}

func TestQueryThingsFail(t *testing.T) {
	author := APIAuthor{
		Token: "dummyToken",
		App:   testApp,
	}
	lr, err := author.QueryThings(
		ThingQueryRequest{
			OwnerID: "",
		})
	if err == nil {
		t.Errorf("should fail")
	}
	if lr != nil {
		t.Errorf("response should be nil")
	}

	lr, err = author.QueryThings(
		ThingQueryRequest{
			OwnerID: "dummyID",
		})
	if err == nil {
		t.Errorf("should fail")
	}
	if lr != nil {
		t.Errorf("response should be nil")
	}
}
