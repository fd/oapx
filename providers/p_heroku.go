package providers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"code.google.com/p/goauth2/oauth"
)

var p_heroku = Provider{
	Code:             "hk",
	AuthURL:          "https://id.heroku.com/oauth/authorize",
	TokenURL:         "https://id.heroku.com/oauth/token",
	Scopes:           []string{"identity"},
	get_profile_func: get_profile_heroku,
	exchange_func:    exchange_heroku,
	IconName:         "fa-cloud",
	IconColor:        "#545ab6",
}

func exchange_heroku(transport *oauth.Transport, code string) (*oauth.Token, error) {
	resp, err := http.PostForm("https://id.heroku.com/oauth/token", url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"client_secret": {transport.ClientSecret},
	})
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

	var respData struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
	}

	err = json.Unmarshal(data, &respData)
	if err != nil {
		return nil, err
	}

	token := &oauth.Token{}
	token.AccessToken = respData.AccessToken
	token.RefreshToken = respData.RefreshToken
	token.Expiry = time.Now().Add(time.Duration(respData.ExpiresIn) * time.Second)
	transport.Token = token

	return token, nil
}

func get_profile_heroku(transport *oauth.Transport) (Profile, error) {
	resp, err := transport.Client().Get("https://api.heroku.com/account")
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

	profile := &HerokuProfile{}

	err = json.Unmarshal(data, &profile)
	if err != nil {
		return nil, err
	}

	profile.raw = data

	return profile, nil
}

type HerokuProfile struct {
	DId    string `json:"id"`
	DName  string `json:"name"`
	DEmail string `json:"email"`
	raw    []byte
}

func (p *HerokuProfile) RemoteId() string {
	return "hk:" + p.DId
}

func (p *HerokuProfile) Identifiers() []string {
	return []string{"hk:" + p.DId}
}

func (p *HerokuProfile) Selectors() []string {
	return []string{p.DEmail}
}

func (p *HerokuProfile) Name() string {
	return p.DName
}

func (p *HerokuProfile) Email() string {
	return p.DEmail
}

func (p *HerokuProfile) PictureURL() string {
	return ""
}

func (p *HerokuProfile) RawData() []byte {
	return p.raw
}
