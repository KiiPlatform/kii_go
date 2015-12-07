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

func (a *APIAuthor) newRequest(method, url string, body interface{}) (*http.Request, error) {
	req, err := newRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+a.Token)
	return req, nil
}

// Let Gateway onboard to the cloud.
// When there's no error, OnboardResponse is returned.
func (au *APIAuthor) OnboardGateway(r *OnboardGatewayRequest) (*OnboardResponse, error) {
	req, err := au.newRequest("POST", au.App.ThingIFURL("/onboardings"), r)
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
	url := au.App.ThingIFURL(path)

	req, err := au.newRequest("PUT", url, request)
	if err != nil {
		return err
	}

	if _, err := executeRequest(req); err != nil {
		return err
	}
	return nil
}

// Login as KiiUser.
// If there is no error, KiiUserLoginResponse is returned.
// Notes that after login successfully, api doesn't update token of APIAuthor,
// you should update by yourself with the token in response.
func (au *APIAuthor) LoginAsKiiUser(request KiiUserLoginRequest) (*KiiUserLoginResponse, error) {
	var ret KiiUserLoginResponse
	reqJson, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("https://%s/api/oauth2/token", au.App.HostName())
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("X-Kii-AppID", au.App.AppID)
	req.Header.Set("X-Kii-AppKey", au.App.AppKey)
	bodyStr, err := executeRequest(req)
	if err != nil {
		return nil, err
	} else {
		err = json.Unmarshal(bodyStr, &ret)
		if err != nil {
			return nil, err
		}
		return &ret, nil
	}

}

// Register KiiUser
// If there is no error, KiiUserRegisterResponse is returned.
func (au *APIAuthor) RegisterKiiUser(request KiiUserRegisterRequest) (*KiiUserRegisterResponse, error) {
	var ret KiiUserRegisterResponse
	reqJson, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	url := au.App.CloudURL("/users")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("X-Kii-AppID", au.App.AppID)
	req.Header.Set("X-Kii-AppKey", au.App.AppKey)
	bodyStr, err := executeRequest(req)
	if err != nil {
		return nil, err
	} else {
		err = json.Unmarshal(bodyStr, &ret)
		if err != nil {
			return nil, err
		}
		return &ret, nil
	}

}

// Post command to Thing.
// Notes that it requires Thing already onboard.
// If there is no error, PostCommandRequest is returned.
func (au APIAuthor) PostCommand(thingID string, request PostCommandRequest) (*PostCommandResponse, error) {
	var ret PostCommandResponse
	reqJson, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/targets/THING:%s/commands", thingID)
	url := au.App.ThingIFURL(path)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("authorization", "Bearer "+au.Token)
	bodyStr, err := executeRequest(req)
	if err != nil {
		return nil, err
	} else {
		err = json.Unmarshal(bodyStr, &ret)
		if err != nil {
			return nil, err
		}
		return &ret, nil
	}
}

// Update command results
func (au APIAuthor) UpdateCommandResults(thingID string, commandID string, request UpdateCommandResultsRequest) error {
	reqJson, err := json.Marshal(request)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/targets/thing:%s/commands/%s/action-results", thingID, commandID)
	url := au.App.ThingIFURL(path)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(reqJson))
	if err != nil {
		return err
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("authorization", "Bearer "+au.Token)

	_, err = executeRequest(req)
	return err
}

func (au *APIAuthor) OnboardThingByOwner(request OnboardByOwnerRequest) (*OnboardResponse, error) {
	var ret OnboardResponse
	reqJson, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	url := au.App.ThingIFURL("/onboardings")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/vnd.kii.OnboardingWithThingIDByOwner+json")
	req.Header.Set("authorization", "Bearer "+au.Token)

	bodyStr, err := executeRequest(req)
	if err != nil {
		return nil, err
	} else {
		err = json.Unmarshal(bodyStr, &ret)
		if err != nil {
			return nil, err
		}
		return &ret, nil
	}
}
