package kii_test

import (
        "testing"
        kii "github.com/KiiPlatform/kii_go"
)

func TestAnonymousLogin(t *testing.T) {
        app := kii.App {
                AppID: "9ab34d8b",
                AppKey: "7a950d78956ed39f3b0815f0f001b43b",
                AppLocation: "JP",
        }
        author := kii.APIAuthor {
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
