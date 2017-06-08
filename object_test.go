package kii

import (
	"fmt"
	"testing"
	"time"

	dproxy "github.com/koron/go-dproxy"
)

func TestAppScopeObjectSuccess(t *testing.T) {

	author, _, err := GetLoginKiiUser()
	if err != nil {
		t.Error("fail to get login user", err)
	}

	bn := fmt.Sprintf("myBucket%d", time.Now().UnixNano())

	data := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}
	b := AppBucket{
		BucketName: bn,
	}
	resp, err := author.PostObject(b, data)

	if err != nil {
		t.Errorf("create app scope object failed: %s\n", err)
	}
	id := resp.ObjectID
	obj, err := author.GetObject(b, id)
	if err != nil {
		t.Errorf("failed to GetObject, err: %s", err)
	}

	if k1, _ := dproxy.New(obj).M("key1").String(); k1 != "value1" {
		t.Errorf("value of key1 is invalid")
	}
	if k2, _ := dproxy.New(obj).M("key2").String(); k2 != "value2" {
		t.Errorf("value of key2 is invalid")
	}

	if err := author.DeleteObject(b, id); err != nil {
		t.Errorf("delete object fail:%s\n", err)
	}

	if err := author.DeleteBucket(b); err != nil {
		t.Errorf("delete bucket fail:%s\n", err)
	}
}

func TestUserScopeObjectSuccess(t *testing.T) {

	author, userID, err := GetLoginKiiUser()
	if err != nil {
		t.Error("fail to get login user", err)
	}

	bn := fmt.Sprintf("myBucket%d", time.Now().UnixNano())

	data := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}

	b := UserBucket{
		BucketName: bn,
		UserID:     userID,
	}
	resp, err := author.PostObject(b, data)

	if err != nil {
		t.Errorf("create user scope object failed: %s\n", err)
	}
	id := resp.ObjectID
	obj, err := author.GetObject(b, id)
	if err != nil {
		t.Errorf("failed to GetObject, err: %s", err)
	}

	if k1, _ := dproxy.New(obj).M("key1").String(); k1 != "value1" {
		t.Errorf("value of key1 is invalid")
	}
	if k2, _ := dproxy.New(obj).M("key2").String(); k2 != "value2" {
		t.Errorf("value of key2 is invalid")
	}

	if err := author.DeleteObject(b, id); err != nil {
		t.Errorf("delete object fail:%s\n", err)
	}

	if err := author.DeleteBucket(b); err != nil {
		t.Errorf("delete bucket fail:%s\n", err)
	}
}

func TestThingScopeObjectSuccess(t *testing.T) {
	thingBucket := fmt.Sprintf("myBucket%d", time.Now().UnixNano())

	au, gwID, err := GatewayOnboard()
	if err != nil {
		t.Errorf("fail to onboard gateway %s", err)
	}

	data := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}

	b := ThingBucket{
		BucketName: thingBucket,
		ThingID:    *gwID,
	}
	resp, err := au.PostObject(b, data)
	if err != nil {
		t.Errorf("create user scope object failed: %s\n", err)
	}
	id := resp.ObjectID
	obj, err := au.GetObject(b, id)
	if err != nil {
		t.Errorf("failed to GetObject, err: %s", err)
	}

	if k1, _ := dproxy.New(obj).M("key1").String(); k1 != "value1" {
		t.Errorf("value of key1 is invalid")
	}
	if k2, _ := dproxy.New(obj).M("key2").String(); k2 != "value2" {
		t.Errorf("value of key2 is invalid")
	}

	if err := au.DeleteObject(b, id); err != nil {
		t.Errorf("delete object fail:%s\n", err)
	}

	if err := au.DeleteBucket(b); err != nil {
		t.Errorf("delete bucket fail:%s\n", err)
	}

}
