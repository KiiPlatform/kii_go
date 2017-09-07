package kii

import (
	"testing"
	// dproxy "github.com/koron/go-dproxy"
)

func TestQueryUserSuccess(t *testing.T) {
	// get admin token
	admin, err := AdminLogin(testApp, clientID, clientSecret)
	if err != nil {
		t.Errorf("admin login fail:%s", err)
	}
	// author for register kii user
	author := APIAuthor{
		Token: "",
		App:   testApp,
	}

	var userIDs []string

	loginName1 := randString()
	displayName := randString()

	user, err := author.RegisterKiiUser(
		UserRegisterRequest{
			LoginName:   loginName1,
			DisplayName: displayName,
			Country:     "JP",
			Password:    "password",
		},
	)
	if err != nil {
		t.Error("failed to register user", err)
	} else {
		userIDs = append(userIDs, user.UserID)
	}

	user, err = author.RegisterKiiUser(
		UserRegisterRequest{
			LoginName:   randString(),
			DisplayName: displayName,
			Country:     "JP",
			Password:    "password",
		},
	)
	if err != nil {
		t.Error("failed to register user", err)
	} else {
		userIDs = append(userIDs, user.UserID)
	}

	user, err = author.RegisterKiiUser(
		UserRegisterRequest{
			LoginName:   randString(),
			DisplayName: displayName,
			Country:     "JP",
			Password:    "password",
		},
	)
	if err != nil {
		t.Error("failed to register user", err)
	} else {
		userIDs = append(userIDs, user.UserID)
	}

	// query with loginName
	qreq := QueryUsersRequest{
		UserQuery: Query{
			Clause: EqualsClause("loginName", loginName1),
		},
	}
	qres, err := admin.QueryUsers(qreq)
	if err != nil {
		t.Error("failed to query user: ", err)
	}
	if len(qres.Results) != 1 {
		t.Error("results should have one user")
	}

	// query with displayName
	qreq = QueryUsersRequest{
		UserQuery: Query{
			Clause:     EqualsClause("displayName", displayName),
			OrderBy:    "loginName",
			Descending: true,
		},
	}
	qres, err = admin.QueryUsers(qreq)
	if err != nil {
		t.Error("failed to query user: ", err)
	}
	if len(qres.Results) != 3 {
		t.Error("results should have 3 user, but have ", len(qres.Results))
	}

	// query with bestEffortLimit
	qreq = QueryUsersRequest{
		UserQuery: Query{
			Clause:     EqualsClause("displayName", displayName),
			OrderBy:    "loginName",
			Descending: true,
		},
		BestEffortLimit: "2",
	}
	qres, err = admin.QueryUsers(qreq)
	if err != nil {
		t.Error("failed to query user: ", err)
	}
	if len(qres.Results) != 2 {
		t.Error("results should have 2 user, but have ", len(qres.Results))
	}
	if qres.NextPaginationKey == "" {
		t.Error("NextPaginationKey should not empty")
	}

	// query with NextPaginationKey
	qreq = QueryUsersRequest{
		UserQuery: Query{
			Clause:     EqualsClause("displayName", displayName),
			OrderBy:    "loginName",
			Descending: true,
		},
		PaginationKey: qres.NextPaginationKey,
	}
	qres, err = admin.QueryUsers(qreq)
	if err != nil {
		t.Error("failed to query user: ", err)
	}
	if len(qres.Results) != 1 {
		t.Error("results should have 1 user, but have ", len(qres.Results))
	}

	// delete users after tests
	for _, user := range userIDs {

		err = admin.DeleteKiiUser(user)
		if err != nil {
			t.Error("failed to delete user: ", err)
		}
	}
}

func TestQueryUserFail(t *testing.T) {
	admin := APIAuthor{
		Token: "dummyToken",
		App:   testApp,
	}

	qreq := QueryUsersRequest{
		UserQuery: Query{
			Clause: EqualsClause("loginName", "test user"),
		},
	}
	_, err := admin.QueryUsers(qreq)
	if err == nil {
		t.Error("shoudl fail")
	}
}

func TestDeleteKiiUserSuccess(t *testing.T) {
	// get admin token
	admin, err := AdminLogin(testApp, clientID, clientSecret)
	if err != nil {
		t.Errorf("admin login fail:%s", err)
	}
	// author for register kii user
	author := APIAuthor{
		Token: "",
		App:   testApp,
	}

	loginName1 := randString()

	user, err := author.RegisterKiiUser(
		UserRegisterRequest{
			LoginName:   loginName1,
			DisplayName: "test user 1",
			Country:     "JP",
			Password:    "password",
		},
	)
	if err != nil {
		t.Error("failed to register user", err)
	}

	err = admin.DeleteKiiUser(user.UserID)
	if err != nil {
		t.Error("failed to delete kii user: ", err)
	}

	// query with bestEffortLimit
	qreq := QueryUsersRequest{
		UserQuery: Query{
			Clause: EqualsClause("loginName", loginName1),
		},
	}
	qres, err := admin.QueryUsers(qreq)
	if err != nil {
		t.Error("failed to query user: ", err)
	}
	if len(qres.Results) != 0 {
		t.Error("results should be empty, but has ", len(qres.Results))
	}
}

func TestDeleteKiiUserFail(t *testing.T) {
	admin := APIAuthor{
		Token: "dummyToken",
		App:   testApp,
	}
	err := admin.DeleteKiiUser("dummyId")
	if err == nil {
		t.Error("should fail to delete user")
	}
}
