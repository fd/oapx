package providers

import (
	"fmt"
	"strings"

	"code.google.com/p/goauth2/oauth"
)

type Provider struct {
	Code             string
	ClientId         string
	ClientSecret     string
	AuthURL          string
	TokenURL         string
	Scopes           []string
	get_profile_func providerProfileFetcher
	exchange_func    exchangeFunc
	IconName         string
	IconColor        string
}

type exchangeFunc func(transport *oauth.Transport, code string) (*oauth.Token, error)

func New(typ, client_id, client_secret string) (*Provider, error) {
	provider, ok := providers[typ]
	if !ok {
		return nil, fmt.Errorf("Unknown provider %q", typ)
	}

	ptr := &Provider{}
	*ptr = provider
	ptr.ClientId = client_id
	ptr.ClientSecret = client_secret
	return ptr, nil
}

var providers = map[string]Provider{
	p_google.Code:   p_google,
	p_github.Code:   p_github,
	p_facebook.Code: p_facebook,
	p_heroku.Code:   p_heroku,
}

func (p *Provider) Config() *oauth.Config {
	return &oauth.Config{
		ClientId:       p.ClientId,
		ClientSecret:   p.ClientSecret,
		Scope:          strings.Join(p.Scopes, " "),
		AuthURL:        p.AuthURL,
		TokenURL:       p.TokenURL,
		RedirectURL:    "http://localhost:3000/callback",
		TokenCache:     nil,
		AccessType:     "offline",
		ApprovalPrompt: "auto",
	}
}

func (p *Provider) Transport(token *oauth.Token) *oauth.Transport {
	return &oauth.Transport{Config: p.Config(), Token: token}
}

func (p *Provider) GetProfile(transport *oauth.Transport) (Profile, error) {
	return p.get_profile_func(transport)
}

func (p *Provider) Exchange(transport *oauth.Transport, code string) (*oauth.Token, error) {
	if p.exchange_func != nil {
		return p.exchange_func(transport, code)
	}

	return transport.Exchange(code)
}
