package kii

import (
	"fmt"
	"strings"
)

// App represents Application in Kii Cloud.
type App struct {
	AppID    string
	AppKey   string
	Location string
}

// HostName returns host name of the Application endpoint.
func (a *App) HostName() string {
	lowerLoc := strings.ToLower(a.Location)
	switch lowerLoc {
	case "jp":
		return "api-jp.kii.com"
	case "us":
		return "api.kii.com"
	case "cn":
		return "api-cn3.kii.com"
	case "sg":
		return "api-sg.kii.com"
	default:
		return lowerLoc
	}
}

// ThingIFBaseURL returns thing-if endpoint base url.
func (a *App) ThingIFBaseURL() string {
	return fmt.Sprintf("https://%s/thing-if/apps/%s", a.HostName(), a.AppID)
}

// KiiCloudBaseURL returns Kii Cloud endpoint base url.
func (a *App) KiiCloudBaseURL() string {
	return fmt.Sprintf("https://%s/api/apps/%s", a.HostName(), a.AppID)
}
