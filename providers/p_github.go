package providers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"code.google.com/p/goauth2/oauth"
)

var p_github = Provider{
	Code:             "gh",
	AuthURL:          "https://github.com/login/oauth/authorize",
	TokenURL:         "https://github.com/login/oauth/access_token",
	Scopes:           []string{"user"},
	get_profile_func: get_profile_github,
	IconName:         "fa-github",
	IconColor:        "#000000",
}

func get_profile_github(transport *oauth.Transport) (Profile, error) {
	resp, err := transport.Client().Get("https://api.github.com/user")
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

	profile := &GithubProfile{}

	err = json.Unmarshal(data, &profile)
	if err != nil {
		return nil, err
	}

	profile.raw = data

	return profile, nil
}

type GithubProfile struct {
	DId      int64  `json:"id"`
	DLogin   string `json:"login"`
	DName    string `json:"name"`
	DEmail   string `json:"email"`
	DPicture string `json:"avatar_url"`
	raw      []byte
}

func (p *GithubProfile) RemoteId() string {
	return "gh:" + strconv.FormatInt(p.DId, 10)
}

func (p *GithubProfile) Identifiers() []string {
	return []string{p.RemoteId(), "gh:" + p.DLogin}
}

func (p *GithubProfile) Selectors() []string {
	return []string{p.Email()}
}

func (p *GithubProfile) Name() string {
	return p.DName
}

func (p *GithubProfile) Email() string {
	return p.DEmail
}

func (p *GithubProfile) PictureURL() string {
	if p.DPicture == "" {
		return ""
	}
	return strings.TrimSuffix(p.DPicture, "?") + "?s=256"
}

func (p *GithubProfile) RawData() []byte {
	return p.raw
}
