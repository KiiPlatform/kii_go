package kii

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// APIAuthor represents API author.
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

// OnboardGateway lets Gateway onboard to the cloud.
// When there's no error, OnboardResponse is returned.
func (a *APIAuthor) OnboardGateway(r *OnboardGatewayRequest) (*OnboardResponse, error) {
	req, err := a.newRequest("POST", a.App.ThingIFURL("/onboardings"), r)
	if err != nil {
		return nil, err
	}

	bodyStr, err := executeRequest(req)
	if err != nil {
		return nil, err
	}
	var ret OnboardResponse
	err = json.Unmarshal(bodyStr, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

// GenerateEndNodeToken Requests access token of end node of gateway.
// Notes the APIAuthor should be a Gateway.
// When there's no error, EndNodeTokenResponse is returned.
func (a APIAuthor) GenerateEndNodeToken(gatewayID string, endnodeID string, r *EndNodeTokenRequest) (*EndNodeTokenResponse, error) {
	path := fmt.Sprintf("/things/%s/end-nodes/%s/token", gatewayID, endnodeID)
	url := a.App.CloudURL(path)

	req, err := a.newRequest("POST", url, r)
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

// AddEndNode adds an end node thing to gateway
// Notes that the APIAuthor should be a Gateway
func (a APIAuthor) AddEndNode(gatewayID string, endnodeID string) error {
	path := fmt.Sprintf("/things/%s/end-nodes/%s", gatewayID, endnodeID)
	url := a.App.CloudURL(path)

	req, err := a.newRequest("PUT", url, nil)
	if err != nil {
		return err
	}
	// a.newRequest() don't set Content-Type for nil body. So we must set it
	// explicitly.
	req.Header.Set("Content-Type", "application/json")

	if _, err := executeRequest(req); err != nil {
		return err
	}
	return nil
}

// RegisterThing registers a Thing on Kii Cloud.
// The request must consist of the predefined fields(see RegisterThingRequest).
// If you want to add the custom fileds, you can simply make RegisterThingRequest as anonymous field of your defined request struct, like:
//  type MyRegisterThingRequest struct {
//    RegisterThingRequest
//    MyField1             string
//  }
// Where there is no error, RegisterThingResponse is returned
func (a APIAuthor) RegisterThing(request interface{}) (*RegisterThingResponse, error) {
	// TODO: should be checked that request contains RegisterThingResponse.

	url := a.App.CloudURL("/things")
	req, err := a.App.newRequest("POST", url, request)
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

// UpdateState updates Thing state.
// Notes that the APIAuthor should be already initialized as a Gateway or EndNode
func (a APIAuthor) UpdateState(thingID string, request interface{}) error {
	path := fmt.Sprintf("/targets/thing:%s/states", thingID)
	url := a.App.ThingIFURL(path)

	req, err := a.newRequest("PUT", url, request)
	if err != nil {
		return err
	}

	if _, err := executeRequest(req); err != nil {
		return err
	}
	return nil
}

// LoginAsKiiUser logins as a KiiUser.
// If there is no error, UserLoginResponse is returned.
// Notes that after login successfully, api doesn't update token of APIAuthor,
// you should update by yourself with the token in response.
func (a *APIAuthor) LoginAsKiiUser(request UserLoginRequest) (*UserLoginResponse, error) {
	url := fmt.Sprintf("https://%s/api/oauth2/token", a.App.HostName())
	req, err := a.App.newRequest("POST", url, request)
	if err != nil {
		return nil, err
	}

	bodyStr, err := executeRequest(req)
	if err != nil {
		return nil, err
	}

	var ret UserLoginResponse
	if err := json.Unmarshal(bodyStr, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}

// RegisterKiiUser registers a KiiUser.
// If there is no error, UserRegisterResponse is returned.
func (a *APIAuthor) RegisterKiiUser(request UserRegisterRequest) (*UserRegisterResponse, error) {
	url := a.App.CloudURL("/users")
	req, err := a.App.newRequest("POST", url, request)
	if err != nil {
		return nil, err
	}

	bodyStr, err := executeRequest(req)
	if err != nil {
		return nil, err
	}

	var ret UserRegisterResponse
	if err := json.Unmarshal(bodyStr, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}

// PostCommand posts command to Thing.
// Notes that it requires Thing already onboard.
// If there is no error, PostCommandRequest is returned.
func (a APIAuthor) PostCommand(thingID string, request PostCommandRequest) (*PostCommandResponse, error) {
	path := fmt.Sprintf("/targets/THING:%s/commands", thingID)
	url := a.App.ThingIFURL(path)
	req, err := a.newRequest("POST", url, request)
	if err != nil {
		return nil, err
	}
	bodyStr, err := executeRequest(req)
	if err != nil {
		return nil, err
	}

	var ret PostCommandResponse
	if err := json.Unmarshal(bodyStr, &ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

// UpdateCommandResults updates command results.
func (a APIAuthor) UpdateCommandResults(thingID string, commandID string, request UpdateCommandResultsRequest) error {

	path := fmt.Sprintf("/targets/thing:%s/commands/%s/action-results", thingID, commandID)
	url := a.App.ThingIFURL(path)
	req, err := a.newRequest("PUT", url, request)
	if err != nil {
		return err
	}

	_, err = executeRequest(req)
	return err
}

// OnboardThingByOwner onboards a thing by its owner.
func (a *APIAuthor) OnboardThingByOwner(request OnboardByOwnerRequest) (*OnboardResponse, error) {
	url := a.App.ThingIFURL("/onboardings")
	req, err := a.newRequest("POST", url, request)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/vnd.kii.OnboardingWithThingIDByOwner+json")

	bodyStr, err := executeRequest(req)
	if err != nil {
		return nil, err
	}

	var ret OnboardResponse
	if err := json.Unmarshal(bodyStr, &ret); err != nil {
		return nil, err
	}

	return &ret, nil
}
