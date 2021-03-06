package kii

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"

	dproxy "github.com/koron/go-dproxy"
)

// APIAuthor represents API author.
// Can be Gateway, EndNode or KiiUser, depending on the token.
type APIAuthor struct {
	Token string
	App   App
}

func (a *APIAuthor) newRequest(method, url string, body interface{}) (*request, error) {
	req, err := newRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+a.Token)
	return req, nil
}

// OnboardGateway lets Gateway onboard to the cloud.
// When there's no error, OnboardGatewayResponse is returned.
func (a *APIAuthor) OnboardGateway(r *OnboardGatewayRequest) (*OnboardGatewayResponse, error) {
	req, err := a.newRequest("POST", a.App.ThingIFURL("/onboardings"), r)
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

// GetState get Thing state.
func (a APIAuthor) GetState(thingID string) (interface{}, error) {
	path := fmt.Sprintf("/targets/thing:%s/states", thingID)
	url := a.App.ThingIFURL(path)

	req, err := a.newRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := executeRequest(req)
	if err != nil {
		return nil, err
	}

	var state interface{}
	err = json.Unmarshal(resp, &state)
	if err != nil {
		return nil, err
	}

	return state, nil
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

// DeleteKiiUser deletes kii user by id.
func (a *APIAuthor) DeleteKiiUser(userID string) error {
	url := a.App.CloudURL("/users/" + userID)
	req, err := a.newRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	_, err = executeRequest(req)
	return err
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

// PostTraitCommand posts trait command to Thing.
// Notes that it requires Thing already onboard.
// If there is no error, PostCommandResponse is returned.
func (a APIAuthor) PostTraitCommand(thingID string, request PostCommandRequest) (*PostCommandResponse, error) {
	path := fmt.Sprintf("/targets/THING:%s/commands", thingID)
	url := a.App.ThingIFURL(path)
	req, err := a.newRequest("POST", url, request)
	req.Header.Set("Content-Type", "application/vnd.kii.CommandCreationRequest+json")

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

// UpdateTraitCommandResults updates trait format command results.
func (a APIAuthor) UpdateTraitCommandResults(thingID string, commandID string, request UpdateCommandResultsRequest) error {

	path := fmt.Sprintf("/targets/thing:%s/commands/%s/action-results", thingID, commandID)
	url := a.App.ThingIFURL(path)
	req, err := a.newRequest("PUT", url, request)
	req.Header.Set("Content-Type", "application/vnd.kii.CommandResultsUpdateRequest+json")
	if err != nil {
		return err
	}
	_, err = executeRequest(req)
	return err
}

// GetCommand gets command info
func (a *APIAuthor) GetCommand(thingID, commandID string) (*GetCommandResponse, error) {
	path := fmt.Sprintf("/targets/thing:%s/commands/%s", thingID, commandID)
	url := a.App.ThingIFURL(path)
	req, err := a.newRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	bodyStr, err := executeRequest(req)
	if err != nil {
		return nil, err
	}

	var ret GetCommandResponse
	if err := json.Unmarshal(bodyStr, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}

// OnboardThingByOwner onboards a thing by its owner.
func (a *APIAuthor) OnboardThingByOwner(request OnboardByOwnerRequest) (*OnboardGatewayResponse, error) {
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

	var ret OnboardGatewayResponse
	if err := json.Unmarshal(bodyStr, &ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

// onboardEndnodeWithGateway onboards an endnode
// request must be either OnboardEndnodeWithGatewayVendorThingIDRequest or OnboardEndnodeWithGatewayThingIDRequest
func (a *APIAuthor) onboardEndnodeWithGateway(request interface{}) (*OnboardEndnodeResponse, error) {
	var contentType string
	if reflect.TypeOf(request) == reflect.TypeOf(OnboardEndnodeWithGatewayThingIDRequest{}) {
		contentType = "application/vnd.kii.OnboardingEndNodeWithGatewayThingID+json"
	} else if reflect.TypeOf(request) == reflect.TypeOf(OnboardEndnodeWithGatewayVendorThingIDRequest{}) {
		contentType = "application/vnd.kii.OnboardingEndNodeWithGatewayVendorThingID+json"
	} else {
		return nil, errors.New("request must be either OnboardEndnodeWithGatewayThingIDRequest or OnboardEndnodeWithGatewayVendorThingIDRequest")
	}

	url := a.App.ThingIFURL("/onboardings")

	req, err := a.newRequest("POST", url, request)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)

	bodyStr, err := executeRequest(req)
	if err != nil {
		return nil, err
	}

	var ret OnboardEndnodeResponse
	if err := json.Unmarshal(bodyStr, &ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

// OnboardEndnodeWithGatewayThingID onboards an endnode with thingID of gateway
func (a *APIAuthor) OnboardEndnodeWithGatewayThingID(request OnboardEndnodeWithGatewayThingIDRequest) (*OnboardEndnodeResponse, error) {
	return a.onboardEndnodeWithGateway(request)
}

// OnboardEndnodeWithGatewayVendorThingID onboards an endnode with vendorThingID of gateway
func (a *APIAuthor) OnboardEndnodeWithGatewayVendorThingID(request OnboardEndnodeWithGatewayVendorThingIDRequest) (*OnboardEndnodeResponse, error) {
	return a.onboardEndnodeWithGateway(request)
}

// ListEndNodes request list of endnodes belong to geateway
func (a *APIAuthor) ListEndNodes(gatewayID string, listPara ListRequest) (*ListEndNodesResponse, error) {
	path := fmt.Sprintf("/things/%s/end-nodes", gatewayID)
	v := url.Values{}
	if listPara.BestEffortLimit != 0 {
		v.Set("bestEffortLimit", strconv.Itoa(listPara.BestEffortLimit))
	}
	if listPara.NextPaginationKey != "" {
		v.Set("paginationKey", listPara.NextPaginationKey)
	}
	if len(v) > 0 {
		path += "?" + v.Encode()
	}

	url := a.App.ThingIFURL(path)
	req, err := a.newRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	bodyStr, err := executeRequest(req)
	if err != nil {
		return nil, err
	}
	var ret ListEndNodesResponse
	if err := json.Unmarshal(bodyStr, &ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

// CreateThingScopeObject create Thing scope object
//
// Deprecated: Use PostObject instead.
func (a APIAuthor) CreateThingScopeObject(thingID, bucketName string, object map[string]interface{}) (*CreateObjectResponse, error) {
	ts := ThingBucket{
		BucketName: bucketName,
		ThingID:    thingID,
	}
	return a.PostObject(ts, object)
}

// ListAllThingScopeObjects list all objects of the specified thing scope bucket
func (a APIAuthor) ListAllThingScopeObjects(thingID, bucketName string, listPara ListRequest) (*ListObjectsResponse, error) {
	clause := AllQueryClause()
	request := QueryObjectsRequest{
		BucketQuery: BucketQuery{
			Clause: clause,
		},
	}
	if listPara.BestEffortLimit != 0 {
		request.BestEffortLimit = strconv.Itoa(listPara.BestEffortLimit)
	}
	if listPara.NextPaginationKey != "" {
		request.PaginationKey = listPara.NextPaginationKey
	}
	resp, err := a.QueryObjects(thingID, bucketName, request)
	if err != nil {
		return nil, err
	}
	return &ListObjectsResponse{
		Results:           resp.Results,
		NextPaginationKey: resp.NextPaginationKey,
	}, nil
}

// DeleteThingScopeBucket delete ThingScope bucket
//
// Deprecated: Use DeleBucket instead.
func (a APIAuthor) DeleteThingScopeBucket(thingID, bucketName string) error {
	ts := ThingBucket{
		BucketName: bucketName,
		ThingID:    thingID,
	}
	return a.DeleteBucket(ts)
}

//QueryObjects query objects of bucket under Thing Scope
func (a APIAuthor) QueryObjects(thingID, bucketName string, request QueryObjectsRequest) (*QueryObjectResponse, error) {
	path := fmt.Sprintf("/things/%s/buckets/%s/query", thingID, bucketName)
	url := a.App.CloudURL(path)

	req, err := a.newRequest("POST", url, request)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-type", "application/vnd.kii.QueryRequest+json")
	if err != nil {
		return nil, err
	}

	bodyStr, err := executeRequest(req)
	if err != nil {
		return nil, err
	}

	var ret QueryObjectResponse
	if err := json.Unmarshal(bodyStr, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}

//QueryUsers query users
func (a APIAuthor) QueryUsers(request QueryUsersRequest) (*QueryUsersResponse, error) {
	url := a.App.CloudURL("/users/query")

	req, err := a.newRequest("POST", url, request)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-type", "application/vnd.kii.UserQueryRequest+json")
	if err != nil {
		return nil, err
	}

	bodyStr, err := executeRequest(req)
	if err != nil {
		return nil, err
	}

	var ret QueryUsersResponse
	if err := json.Unmarshal(bodyStr, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}

//UpdateVendorThingID update Vendor ThingID of exsiting Thing
func (a APIAuthor) UpdateVendorThingID(thingID string, request UpdateVendorThingIDRequest) error {
	path := fmt.Sprintf("/things/%s/vendor-thing-id", thingID)
	url := a.App.CloudURL(path)

	req, err := a.newRequest("PUT", url, request)
	if err != nil {
		return err
	}
	req.Header.Set("Content-type", "application/vnd.kii.VendorThingIDUpdateRequest+json")
	if err != nil {
		return err
	}

	_, err = executeRequest(req)
	if err != nil {
		return err
	}
	return nil
}

// GetThing get thing info
func (a APIAuthor) GetThing(thingID string) (interface{}, error) {
	path := fmt.Sprintf("/things/%s", thingID)
	url := a.App.CloudURL(path)
	req, err := a.newRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	bodyStr, err := executeRequest(req)
	if err != nil {
		return nil, err
	}

	var obj interface{}
	err = json.Unmarshal(bodyStr, &obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

// UpdateThing update thing properites.
func (a APIAuthor) UpdateThing(thingID string, data map[string]interface{}) error {
	path := fmt.Sprintf("/things/%s", thingID)
	url := a.App.CloudURL(path)
	req, err := a.newRequest("PATCH", url, data)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/vnd.kii.ThingUpdateRequest+json")

	if _, err := executeRequest(req); err != nil {
		return err
	}
	return nil
}

// DeleteThing delete an exsiting Thing
func (a APIAuthor) DeleteThing(thingID string) error {
	path := fmt.Sprintf("/things/%s", thingID)
	url := a.App.CloudURL(path)
	req, err := a.newRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	_, err = executeRequest(req)
	if err != nil {
		return err
	}
	return nil
}

// ReportEndnodeStatus reports online status of endnode by gateway
func (a APIAuthor) ReportEndnodeStatus(gatewayID, endnodeID string, request ReportEndnodeStatusRequest) error {
	path := fmt.Sprintf("/things/%s/end-nodes/%s/connection", gatewayID, endnodeID)
	url := a.App.ThingIFURL(path)

	req, err := a.newRequest("PUT", url, request)
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

// UpdateMultipleTraitState updates multiple traits formatted Thing state. You should first define trait, thingType,
// firmwareVersion and alias in server. Then onboard the thing with thingType and firmwareVersion
// Notes that the APIAuthor should be already initialized as a Gateway or EndNode
func (a APIAuthor) UpdateMultipleTraitState(thingID string, request interface{}) error {
	path := fmt.Sprintf("/targets/thing:%s/states", thingID)
	url := a.App.ThingIFURL(path)

	req, err := a.newRequest("PUT", url, request)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/vnd.kii.MultipleTraitState+json")

	if _, err := executeRequest(req); err != nil {
		return err
	}
	return nil
}

// UpdateTraitState updates Thing state with specified trait/alias. You should first define trait, thingType,
// firmwareVersion and alias in server. Then onboard the thing with thingType and firmwareVersion
// Notes that the APIAuthor should be already initialized as a Gateway or EndNode
func (a APIAuthor) UpdateTraitState(thingID string, alias string, request interface{}) error {
	path := fmt.Sprintf("/targets/thing:%s/states/aliases/%s", thingID, alias)
	url := a.App.ThingIFURL(path)

	req, err := a.newRequest("PUT", url, request)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/vnd.kii.TraitState+json")

	if _, err := executeRequest(req); err != nil {
		return err
	}
	return nil
}

//InstallMqtt a MQTT installation to the Kii cloud for current logged in user.
func (a APIAuthor) InstallMqtt(development bool) (installationID string, err error) {
	url := a.App.CloudURL("/installations")
	req, err := a.newRequest("POST", url, map[string]interface{}{
		"deviceType":  "MQTT",
		"development": development,
	})
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-type", "application/vnd.kii.InstallationCreationRequest+json")

	bodyStr, err := executeRequest(req)
	if err != nil {
		return "", err
	}
	var resp interface{}
	err = json.Unmarshal(bodyStr, &resp)
	if err != nil {
		return "", err
	}
	return dproxy.New(resp).M("installationID").String()
}

// GetMqttEndpoint gets mqtt endpoint with specified installationID.
func (a APIAuthor) GetMqttEndpoint(installationID string) (endpoint *MqttEndpoint, err error) {
	path := fmt.Sprintf("/installations/%s/mqtt-endpoint", installationID)
	url := a.App.CloudURL(path)
	req, err := a.newRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	bodyStr, err := executeRequest(req)
	if err != nil {
		return nil, err
	}

	var ret MqttEndpoint
	if err := json.Unmarshal(bodyStr, &ret); err != nil {
		return nil, err
	}
	return &ret, nil

}

// PostObject creates a kii object with data
func (a APIAuthor) PostObject(bucket Bucket, data map[string]interface{}) (*CreateObjectResponse, error) {
	path := fmt.Sprintf("/%s/objects", bucket.Path())
	url := a.App.CloudURL(path)
	req, err := a.newRequest("POST", url, data)
	if err != nil {
		return nil, err
	}

	bodyStr, err := executeRequest(req)
	if err != nil {
		return nil, err
	}

	var ret CreateObjectResponse
	if err := json.Unmarshal(bodyStr, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}

// GetObject retrieves kii object with object ID
func (a APIAuthor) GetObject(bucket Bucket, objectID string) (interface{}, error) {
	path := fmt.Sprintf("/%s/objects/%s", bucket.Path(), objectID)
	url := a.App.CloudURL(path)
	req, err := a.newRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := executeRequest(req)
	if err != nil {
		return nil, err
	}

	var obj interface{}
	err = json.Unmarshal(resp, &obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

// DeleteObject deletes a kii object with object ID
func (a APIAuthor) DeleteObject(bucket Bucket, objectID string) error {
	path := fmt.Sprintf("/%s/objects/%s", bucket.Path(), objectID)
	url := a.App.CloudURL(path)
	req, err := a.newRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	_, err = executeRequest(req)
	if err != nil {
		return err
	}
	return nil
}

//DeleteBucket deletes bucket
func (a APIAuthor) DeleteBucket(bucket Bucket) error {
	url := a.App.CloudURL(bucket.Path())

	req, err := a.newRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	_, err = executeRequest(req)
	if err != nil {
		return err
	}
	return nil
}

//QueryThings query things owned by user
func (a APIAuthor) QueryThings(request ThingQueryRequest) (*QueryThingsResponse, error) {
	if request.OwnerID == "" {
		return nil, errors.New("OwnerID must not be empty")
	}

	clause := Clause{
		"type":  "contains",
		"field": "userOwners",
		"value": request.OwnerID,
	}
	if request.Clause != nil {
		clause = AndClause(clause, request.Clause)
	}

	requestObj := map[string]interface{}{
		"thingQuery": map[string]interface{}{
			"clause": clause,
		},
	}

	if request.BestEffortLimit != 0 {
		requestObj["bestEffortLimit"] = strconv.Itoa(request.BestEffortLimit)
	}
	if request.NextPaginationKey != "" {
		requestObj["paginationKey"] = request.NextPaginationKey
	}

	url := a.App.CloudURL("/things/query")

	req, err := a.newRequest("POST", url, requestObj)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/vnd.kii.ThingQueryRequest+json")

	bodyStr, err := executeRequest(req)
	if err != nil {
		return nil, err
	}
	var ret QueryThingsResponse
	err = json.Unmarshal(bodyStr, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

// ResetThingPassword reset password of existing thing.
func (a APIAuthor) ResetThingPassword(thingID, newPassword string) error {
	path := fmt.Sprintf("/things/%s/password", thingID)
	url := a.App.CloudURL(path)

	req, err := a.newRequest("PUT", url, map[string]string{
		"newPassword": newPassword,
	})
	if err != nil {
		return err
	}
	// a.newRequest() don't set Content-Type for nil body. So we must set it
	// explicitly.
	req.Header.Set("Content-Type", "application/vnd.kii.ChangeThingPasswordRequest+json")

	if _, err := executeRequest(req); err != nil {
		return err
	}
	return nil
}
