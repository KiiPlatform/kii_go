package kii

import (
	"testing"
)

func TestGetMqttEndpoint(t *testing.T) {
	// get a login user
	author, _, err := GetLoginKiiUser()
	if err != nil {
		t.Errorf("fail to get login user")
	}

	installID, err := author.InstallMqtt(false)
	if err != nil {
		t.Errorf("fail to install mqtt")
	}
	if installID == "" {
		t.Errorf("installationID is empty")
	}

	endpoint, err := author.GetMqttEndpoint(installID)
	if err != nil {
		t.Errorf("fail to get mqtt endpoint")
	}
	if endpoint == nil {
		t.Errorf("endpint is null")
	} else {
		if endpoint.Password == "" || endpoint.InstallationID == "" ||
			endpoint.Host == "" || endpoint.MqttTopic == "" ||
			endpoint.Username == "" || endpoint.PortSSL == 0 ||
			endpoint.PortTCP == 0 || endpoint.PortWS == 0 ||
			endpoint.PortWSS == 0 || endpoint.XMqttTTL == 0 {
			t.Errorf("endpoint is invalid")
		}
	}
}
