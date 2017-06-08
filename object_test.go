package kii

import (
	"fmt"
	"testing"
	"time"

	dproxy "github.com/koron/go-dproxy"
)

func TestCreateAppScopeObjectSuccess(t *testing.T) {

	author, _, err := GetLoginKiiUser()
	if err != nil {
		t.Error("fail to get login user", err)
	}

	scope := Scope{
		Type: APP,
	}

	bn := fmt.Sprintf("myBucket%d", time.Now().UnixNano())

	data := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}
	resp, err := author.CreateObject(scope, bn, data)

	if err != nil {
		t.Errorf("create app scope object failed: %s\n", err)
	}
	id := resp.ObjectID
	obj, err := author.GetObject(scope, bn, id)
	if err != nil {
		t.Errorf("failed to GetObject, err: %s", err)
	}

	if k1, _ := dproxy.New(obj).M("key1").String(); k1 != "value1" {
		t.Errorf("value of key1 is invalid")
	}
	if k2, _ := dproxy.New(obj).M("key2").String(); k2 != "value2" {
		t.Errorf("value of key2 is invalid")
	}

	if err := author.DeleteObject(scope, bn, id); err != nil {
		t.Errorf("delete object fail:%s\n", err)
	}

	if err := author.DeleteBucket(scope, bn); err != nil {
		t.Errorf("delete bucket fail:%s\n", err)
	}
}

func TestCreateUserScopeObjectSuccess(t *testing.T) {

	author, userID, err := GetLoginKiiUser()
	if err != nil {
		t.Error("fail to get login user", err)
	}

	scope := Scope{
		Type: USER,
		ID:   userID,
	}

	bn := fmt.Sprintf("myBucket%d", time.Now().UnixNano())

	data := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}
	resp, err := author.CreateObject(scope, bn, data)

	if err != nil {
		t.Errorf("create app scope object failed: %s\n", err)
	}
	id := resp.ObjectID
	obj, err := author.GetObject(scope, bn, id)
	if err != nil {
		t.Errorf("failed to GetObject, err: %s", err)
	}

	if k1, _ := dproxy.New(obj).M("key1").String(); k1 != "value1" {
		t.Errorf("value of key1 is invalid")
	}
	if k2, _ := dproxy.New(obj).M("key2").String(); k2 != "value2" {
		t.Errorf("value of key2 is invalid")
	}

	if err := author.DeleteObject(scope, bn, id); err != nil {
		t.Errorf("delete object fail:%s\n", err)
	}

	if err := author.DeleteBucket(scope, bn); err != nil {
		t.Errorf("delete bucket fail:%s\n", err)
	}

}
