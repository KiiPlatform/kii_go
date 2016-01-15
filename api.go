package kii

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Logger set logger for kii module.  Default is discard logger, no logs are
// outputed.  Please set a valid logger if you want to make kii module put
// logs.
//
//	kii.Logger = log.New(os.Stderr, "", log.LstdFlags)
var Logger = log.New(ioutil.Discard, "", 0)

// LayoutPosition represents Layout position of the Thing.
type LayoutPosition int

const (
	// ENDNODE represents layout position of endnodes.
	ENDNODE LayoutPosition = iota
	// STANDALONE represents layout position of standalone.
	STANDALONE
	// GATEWAY represents layout position of gateway.
	GATEWAY
)

// Obtain Layout postion of the Thing in string.
func (lp LayoutPosition) String() string {
	switch lp {
	case ENDNODE:
		return "END_NODE"
	case STANDALONE:
		return "STANDALONE"
	case GATEWAY:
		return "GATEWAY"
	default:
		return fmt.Sprintf("!LayoutPosition(%d)", lp)
	}
}

// OnboardGatewayRequest for requesting Gateway Onboard.
type OnboardGatewayRequest struct {
	VendorThingID   string                 `json:"vendorThingID"`
	ThingPassword   string                 `json:"thingPassword"`
	ThingType       string                 `json:"thingType"`
	LayoutPosition  string                 `json:"layoutPosition"`
	ThingProperties map[string]interface{} `json:"thingProperties"`
}

var _ contentTyper = (*OnboardGatewayRequest)(nil)

func (r *OnboardGatewayRequest) contentType() string {
	return "application/vnd.kii.onboardingWithVendorThingIDByThing+json"
}

// OnboardGatewayResponse for receiving response of Gateway Onboard.
type OnboardGatewayResponse struct {
	ThingID      string       `json:"thingID"`
	AccessToken  string       `json:"accessToken"`
	MqttEndpoint MqttEndpoint `json:"mqttEndpoint"`
}

// MqttEndpoint represents MQTT endpoint.
type MqttEndpoint struct {
	InstallationID string `json:"installationID"`
	Host           string `json:"host"`
	MqttTopic      string `json:"mqttTopic"`
	Username       string `json:"userName"`
	Password       string `json:"password"`
	PortSSL        int    `json:"portSSL"`
	PortTCP        int    `json:"portTCP"`
}

// EndNodeTokenRequest for requesting end node token
type EndNodeTokenRequest struct {
	ExpiresIn string `json:"expires_in,omitempty"`
}

// EndNodeTokenResponse for receiving response of end node token
type EndNodeTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	ThingID      string `json:"id"`
	RefreshToken string `json:"refresh_token"`
}

// RegisterThingRequest reresents predefined fileds for requesting Thing Registration.
type RegisterThingRequest struct {
	VendorThingID   string `json:"_vendorThingID"`
	ThingPassword   string `json:"_password"`
	ThingType       string `json:"_thingType,omitempty"`
	LayoutPosition  string `json:"_layoutPosition,omitempty"`
	Vendor          string `json:"_vendor,omitempty"`
	FirmwareVersion string `json:"_firmwareVersion,omitempty"`
	Lot             string `json:"_lot,omitempty"`
	StringField1    string `json:"_stringField1,omitempty"`
	StringField2    string `json:"_stringField2,omitempty"`
	StringField3    string `json:"_stringField3,omitempty"`
	StringField4    string `json:"_stringField4,omitempty"`
	StringField5    string `json:"_stringField5,omitempty"`
	NumberField1    int64  `json:"_numberField1,omitempty"`
	NumberField2    int64  `json:"_numberField2,omitempty"`
	NumberField3    int64  `json:"_numberField3,omitempty"`
	NumberField4    int64  `json:"_numberField4,omitempty"`
	NumberField5    int64  `json:"_numberField5,omitempty"`
}

// RegisterThingResponse for receiving response of end node token
type RegisterThingResponse struct {
	ThingID        string `json:"_thingID"`
	VendorThingID  string `json:"_vendorThingID"`
	ThingType      string `json:"_thingType"`
	LayoutPosition string `json:"_layoutPosition"`
	Created        int    `json:"_created"`
	Disabled       bool   `json:"_disabled"`
}

// UserRegisterRequest for request registration of KiiUser.
// At least one of LoginName, EmailAddress or PhoneNumber must be provided.
type UserRegisterRequest struct {
	LoginName           string `json:"loginName,omitempty"`
	DisplayName         string `json:"displayName,omitempty"`
	Country             string `json:"country,omitempty"`
	Locale              string `json:"locale,omitempty"`
	EmailAddress        string `json:"emailAddress,omitempty"`
	PhoneNumber         string `json:"phoneNumber,omitempty"`
	PhoneNumberVerified bool   `json:"phoneNumberVerified,omitempty"`
	Password            string `json:"password"`
}

// UserRegisterResponse for receiving registration of KiiUser.
type UserRegisterResponse struct {
	UserID              string `json:"userID"`
	LoginName           string `json:"loginName"`
	DisplayName         string `json:"displayName"`
	Country             string `json:"country"`
	Locale              string `json:"locale"`
	EmailAddress        string `json:"emailAddress"`
	PhoneNumber         string `json:"phoneNumber"`
	PhoneNumberVerified bool   `json:"phoneNumberVerified"`
	HasPassword         bool   `json:"_hasPassword"`
}

// UserLoginRequest for requesting login of KiiUser
type UserLoginRequest struct {
	UserName     string `json:"username"`
	Password     string `json:"password"`
	ExpiresAt    string `json:"expiresAt,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	GrantType    string `json:"grant_type,omitempty"`
}

// UserLoginResponse for receiving response of login
type UserLoginResponse struct {
	ID           string `json:"id"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
}

// PostCommandRequest for posting command
// Issuer can be group or user.
// If user, must be "user:<user-id>".
type PostCommandRequest struct {
	Issuer           string                   `json:"issuer"`
	Actions          []map[string]interface{} `json:"actions"`
	Schema           string                   `json:"schema"`
	SchemaVersion    int                      `json:"schemaVersion"`
	FiredByTriggerID string                   `json:"firedByTriggerID,omitempty"`
	Title            string                   `json:"title,omitempty"`
	Description      string                   `json:"description,omitempty"`
	Metadata         map[string]interface{}   `json:"metadata,omitempty"`
}

// PostCommandResponse for receiving response of posting command
type PostCommandResponse struct {
	CommandID string `json:"commandID"`
}

// OnboardByOwnerRequest for requesting Onboard by Thing Owner.
type OnboardByOwnerRequest struct {
	ThingID        string `json:"thingID"`
	ThingPassword  string `json:"thingPassword"`
	Owner          string `json:"owner"`
	LayoutPosition string `json:"layoutPosition,omitempty"` // pattern: GATEWAY|STANDALONE|ENDNODE, STANDALONE by default
}

// OnboardEndnodeRequestCommon is the command fields for
// OnboardEndnodeWithGatewayThingIDRequest and OnboardEndnodeWithGatewayVendorThingIDRequest
type OnboardEndnodeRequestCommon struct {
	EndNodeVendorThingID   string `json:"endNodeVendorThingID"`
	EndNodePassword        string `json:"endNodePassword"`
	Owner                  string `json:"owner"`
	EndNodeThingProperties string `json:"endNodeThingProperties,omitempty"`
	EndNodeThingType       string `json:"endNodeThingType,omitempty"`
}

// OnboardEndnodeWithGatewayThingIDRequest for requesting Onboard with thingID of gateway
type OnboardEndnodeWithGatewayThingIDRequest struct {
	GatewayThingID string `json:"gatewayThingID"`
	OnboardEndnodeRequestCommon
}

// OnboardEndnodeWithGatewayVendorThingIDRequest for requesting Onboard with vendorThingID of gateway
type OnboardEndnodeWithGatewayVendorThingIDRequest struct {
	GatewayVendorThingID string `json:"gatewayVendorThingID"`
	OnboardEndnodeRequestCommon
}

// OnboardEndnodeResponse for receiving response of onboarding endnode with gateway
type OnboardEndnodeResponse struct {
	AccessToken    string `json:"accessToken"`
	EndNodeThingID string `json:"endNodeThingID"`
}

// UpdateCommandResultsRequest for updating command results
type UpdateCommandResultsRequest struct {
	ActionResults []map[string]interface{} `json:"actionResults"`
}

// EndNode represents end-node
type EndNode struct {
	ThingID       string `json:"thingID"`
	VendorThingID string `json:"vendorThingID"`
}

// ListEndNodesResponse for receiving response of list request
type ListEndNodesResponse struct {
	Results           []EndNode `json:"results"`
	NextPaginationKey string    `json:"nextPaginationKey"`
}

// ListRequest consist of parameters when request list of
// data(like end-nodes) from Kii Cloud
type ListRequest struct {
	BestEffortLimit   int
	NextPaginationKey string
}

// AnonymousLogin logins as Anonymous user.
// When there's no error, APIAuthor is returned.
func AnonymousLogin(app App) (*APIAuthor, error) {
	type AnonymousLoginRequest struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		GrantType    string `json:"grant_type"`
	}
	type AnonymousLoginResponse struct {
		ID          string `json:"id"`
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		TokenType   string `json:"token_type"`
	}
	reqObj := AnonymousLoginRequest{
		ClientID:     app.AppID,
		ClientSecret: app.AppKey,
		GrantType:    "client_credentials",
	}
	req, err := newRequest("POST", app.CloudURL("/oauth2/token"), &reqObj)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// TODO: return as an error.
	}

	bodyStr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// FIXME: should be removed after debug?
	Logger.Println("body: " + string(bodyStr))

	var respObj AnonymousLoginResponse
	err = json.Unmarshal(bodyStr, &respObj)
	if err != nil {
		return nil, err
	}
	return &APIAuthor{
		Token: respObj.AccessToken,
		App:   app,
	}, nil
}
