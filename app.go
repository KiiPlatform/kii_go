package kii

import "strings"

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

// rootURL returns app's endpoint root URL.
func (a *App) rootURL() string {
	return "https://" + a.HostName()
}

// CloudURL returns regular API URL for the app.
func (a *App) CloudURL(path string) string {
	return a.rootURL() + "/api/apps/" + a.AppID + path
}

// ThingIFURL returns Thing-IF API URL for the app.
func (a *App) ThingIFURL(path string) string {
	return a.rootURL() + "/thing-if/apps/" + a.AppID + path
}

func (a *App) newRequest(method, url string, body interface{}) (*request, error) {
	req, err := newRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Kii-AppID", a.AppID)
	req.Header.Set("X-Kii-AppKey", a.AppKey)
	return req, nil
}
