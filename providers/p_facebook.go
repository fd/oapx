package providers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"code.google.com/p/goauth2/oauth"
)

var p_facebook = Provider{
	Code:             "fb",
	AuthURL:          "https://www.facebook.com/dialog/oauth",
	TokenURL:         "https://graph.facebook.com/oauth/access_token",
	Scopes:           []string{"basic_info", "email"},
	get_profile_func: get_profile_facebook,
	IconName:         "fa-facebook",
	IconColor:        "#3B5998",
}

func get_profile_facebook(transport *oauth.Transport) (Profile, error) {
	resp, err := transport.Client().Get("https://graph.facebook.com/me?fields=id,name,first_name,last_name,email,username,picture.width(256).height(256)")
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

	profile := &FacebookProfile{}

	err = json.Unmarshal(data, &profile)
	if err != nil {
		return nil, err
	}

	profile.raw = data

	return profile, nil
}

type FacebookProfile struct {
	DId       string `json:"id"`
	DName     string `json:"name"`
	DEmail    string `json:"email"`
	DUsername string `json:"username"`
	DPicture  struct {
		Data struct {
			URL string `json:"url"`
		} `json:"data"`
	} `json:"picture"`
	raw []byte
}

func (p *FacebookProfile) RemoteId() string {
	return "fb:" + p.DId
}

func (p *FacebookProfile) Identifiers() []string {
	return []string{"fb:" + p.DUsername, "fb:" + p.DId}
}

func (p *FacebookProfile) Selectors() []string {
	return []string{p.DEmail}
}

func (p *FacebookProfile) Name() string {
	return p.DName
}

func (p *FacebookProfile) Email() string {
	return p.DEmail
}

func (p *FacebookProfile) PictureURL() string {
	return p.DPicture.Data.URL
}

func (p *FacebookProfile) RawData() []byte {
	return p.raw
}
