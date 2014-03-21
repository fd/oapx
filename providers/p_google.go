package providers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"code.google.com/p/goauth2/oauth"
)

var p_google = Provider{
	Code:     "go",
	AuthURL:  "https://accounts.google.com/o/oauth2/auth",
	TokenURL: "https://accounts.google.com/o/oauth2/token",
	Scopes: []string{
		"openid",
		"profile",
		"email",
	},
	get_profile_func: get_profile_google,
	IconName:         "fa-google-plus",
	IconColor:        "#F90101",
}

func get_profile_google(transport *oauth.Transport) (Profile, error) {
	resp, err := transport.Client().Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("http status: %d", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	profile := &GoogleProfile{}

	err = json.Unmarshal(data, &profile)
	if err != nil {
		return nil, err
	}

	profile.raw = data

	return profile, nil
}

type GoogleProfile struct {
	DId        string `json:"id"`
	DName      string `json:"name"`
	DEmail     string `json:"email"`
	DPicture   string `json:"picture"`
	DAppDomain string `json:"hd"`
	raw        []byte
}

func (p *GoogleProfile) RemoteId() string {
	return "go:" + p.DId
}

func (p *GoogleProfile) Identifiers() []string {
	return []string{"go:" + p.DId}
}

func (p *GoogleProfile) Selectors() []string {
	s := []string{p.Email()}
	if p.DAppDomain != "" {
		s = append(s, "go:app:"+p.DAppDomain)
	}
	return s
}

func (p *GoogleProfile) Name() string {
	return p.DName
}

func (p *GoogleProfile) Email() string {
	return p.DEmail
}

func (p *GoogleProfile) PictureURL() string {
	if p.DPicture == "" {
		return ""
	}
	return p.DPicture + "?sz=256"
}

func (p *GoogleProfile) RawData() []byte {
	return p.raw
}
