package kii

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Struct represents API author.
// Can be Gateway, EndNode or KiiUser, depending on the token.
type APIAuthor struct {
	Token string
	App   App
}

// Let Gateway onboard to the cloud.
// When there's no error, OnboardGatewayResponse is returned.
func (au *APIAuthor) OnboardGateway(request OnboardGatewayRequest) (*OnboardGatewayResponse, error) {
	reqJSON, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	url := au.App.ThingURL("/onboardings")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqJSON))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/vnd.kii.onboardingWithVendorThingIDByThing+json")
	req.Header.Set("Authorization", "Bearer "+au.Token)

	bodyStr, err := executeRequest(req)
	if err != nil {
		return nil, err
	}
	var ret OnboardGatewayResponse
	err = json.Unmarshal(bodyStr, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

// Request access token of end node of gateway.
// Notes the APIAuthor should be a Gateway.
// When there's no error, EndNodeTokenResponse is returned.
func (au APIAuthor) GenerateEndNodeToken(gatewayID string, endnodeID string, request EndNodeTokenRequest) (*EndNodeTokenResponse, error) {
	path := fmt.Sprintf("/things/%s/end-nodes/%s/token", gatewayID, endnodeID)
	url := au.App.CloudURL(path)

	reqJSON, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+au.Token)

	bodyStr, err := executeRequest(req)
	if err != nil {
		return nil, err
	}
	var ret EndNodeTokenResponse
	err = json.Unmarshal(bodyStr, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

// Add an end node thing to gateway
// Notes that the APIAuthor should be a Gateway
func (au APIAuthor) AddEndNode(gatewayID string, endnodeID string) error {
	path := fmt.Sprintf("/things/%s/end-nodes/%s", gatewayID, endnodeID)
	url := au.App.CloudURL(path)

	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+au.Token)

	if _, err := executeRequest(req); err != nil {
		return err
	}
	return nil
}

// Register Thing.
// The request must consist of the predefined fields(see RegisterThingRequest).
// If you want to add the custom fileds, you can simply make RegisterThingRequest as anonymous field of your defined request struct, like:
//  type MyRegisterThingRequest struct {
//    RegisterThingRequest
//    MyField1             string
//  }
// Where there is no error, RegisterThingResponse is returned
func (au APIAuthor) RegisterThing(request interface{}) (*RegisterThingResponse, error) {
	reqJSON, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	url := au.App.CloudURL("/things")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqJSON))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/vnd.kii.ThingRegistrationRequest+json")
	req.Header.Set("X-Kii-AppID", au.App.AppID)
	req.Header.Set("X-Kii-AppKey", au.App.AppKey)

	bodyStr, err := executeRequest(req)
	if err != nil {
		return nil, err
	}
	var ret RegisterThingResponse
	err = json.Unmarshal(bodyStr, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

// Update Thing state.
// Notes that the APIAuthor should be already initialized as a Gateway or EndNode
func (au APIAuthor) UpdateState(thingID string, request interface{}) error {

	reqJSON, err := json.Marshal(request)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/targets/thing:%s/states", thingID)
	url := au.App.ThingURL(path)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(reqJSON))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+au.Token)

	if _, err := executeRequest(req); err != nil {
		return err
	}
	return nil
}
