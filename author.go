package kii

import (
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

func (a *APIAuthor) newRequest(method, url string, body interface{}) (*http.Request, error) {
	req, err := newRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+a.Token)
	return req, nil
}

// Let Gateway onboard to the cloud.
// When there's no error, OnboardGatewayResponse is returned.
func (au *APIAuthor) OnboardGateway(r *OnboardGatewayRequest) (*OnboardGatewayResponse, error) {
	req, err := au.newRequest("POST", au.App.ThingURL("/onboardings"), r)
	if err != nil {
		return nil, err
	}

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
func (au APIAuthor) GenerateEndNodeToken(gatewayID string, endnodeID string, r *EndNodeTokenRequest) (*EndNodeTokenResponse, error) {
	path := fmt.Sprintf("/things/%s/end-nodes/%s/token", gatewayID, endnodeID)
	url := au.App.CloudURL(path)

	req, err := au.newRequest("POST", url, r)
	if err != nil {
		return nil, err
	}

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

	req, err := au.newRequest("PUT", url, nil)
	if err != nil {
		return err
	}
	// au.newRequest() don't set Content-Type for nil body. So we must set it
	// explicitly.
	req.Header.Set("Content-Type", "application/json")

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
	// TODO: should be checked that request contains RegisterThingResponse.

	url := au.App.CloudURL("/things")
	req, err := au.App.newRequest("POST", url, request)
	if err != nil {
		return nil, err
	}
	// replace default Content-Type.
	req.Header.Set("Content-Type", "application/vnd.kii.ThingRegistrationRequest+json")

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
	path := fmt.Sprintf("/targets/thing:%s/states", thingID)
	url := au.App.ThingURL(path)

	req, err := au.newRequest("PUT", url, request)
	if err != nil {
		return err
	}

	if _, err := executeRequest(req); err != nil {
		return err
	}
	return nil
}
