package providers

import (
	"code.google.com/p/goauth2/oauth"
)

type providerProfileFetcher func(transport *oauth.Transport) (Profile, error)

type Profile interface {
	RemoteId() string
	Identifiers() []string
	Selectors() []string
	Name() string
	Email() string
	PictureURL() string
	RawData() []byte
}
