package kii_cloud_api_test

import (
        "testing"
        kca "kii_cloud_api"
)

func TestAnonymousLogin(t *testing.T) {
        app := kca.KiiApp {
                AppID: "9ab34d8b",
                AppKey: "7a950d78956ed39f3b0815f0f001b43b",
                AppLocation: "JP",
        }
        author := kca.APIAuthor {
                App: app,
        }
        err := author.AnonymousLogin()
        if err != nil {
                t.Errorf("got error on anonymous login %s", err)
        }
        if len(author.Token) < 1 {
                t.Errorf("failed to get author token %s", author.Token)
        }
}
