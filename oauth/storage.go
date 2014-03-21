package oauth

import (
	"database/sql"
	"fmt"

	"github.com/RangelReale/osin"
	"github.com/jmoiron/sqlx"

	"github.com/fd/oauth2-proxy/data"
)

type Storage struct {
	db *sqlx.DB
}

// GetClient loads the client by id (client_id)
func (s *Storage) GetClient(client_id string) (*osin.Client, error) {
	application, err := data.GetApplicationWithClientId(s.db, client_id)
	if err != nil {
		return nil, err
	}

	if application == nil {
		return nil, fmt.Errorf("no such application")
	}

	client := &osin.Client{
		Id:          application.ClientId,
		Secret:      application.ClientSecret,
		RedirectUri: application.RedirectURI,
		UserData:    application,
	}

	return client, nil
}

func (s *Storage) getClientWithId(id int64) (*osin.Client, error) {
	application, err := data.GetApplicationWithId(s.db, id)
	if err != nil {
		return nil, err
	}

	if application == nil {
		return nil, fmt.Errorf("no such application")
	}

	client := &osin.Client{
		Id:          application.ClientId,
		Secret:      application.ClientSecret,
		RedirectUri: application.RedirectURI,
		UserData:    application,
	}

	return client, nil
}

// SaveAuthorize saves authorize data.
func (s *Storage) SaveAuthorize(ad *osin.AuthorizeData) error {
	application := ad.Client.UserData.(*data.Application)

	identity := ad.UserData.(*data.Identity)

	authorization := &data.Authorization{
		IdentityId:    identity.Id,
		ApplicationId: application.Id,
		Code:          ad.Code,
		ExpiresIn:     ad.ExpiresIn,
		State:         ad.State,
		Scope:         ad.Scope,
		RedirectURI:   ad.RedirectUri,
	}

	return data.CreateAuthorization(s.db, authorization)
}

// LoadAuthorize looks up AuthorizeData by a code.
// Client information MUST be loaded together.
// Optionally can return error if expired.
func (s *Storage) LoadAuthorize(code string) (*osin.AuthorizeData, error) {
	authorization, err := data.GetAuthorizationWithCode(s.db, code)
	if err != nil {
		return nil, err
	}

	if authorization == nil {
		return nil, fmt.Errorf("no such authorization")
	}

	client, err := s.getClientWithId(authorization.ApplicationId)
	if err != nil {
		return nil, err
	}

	ad := &osin.AuthorizeData{
		Client:      client,
		Code:        authorization.Code,
		ExpiresIn:   authorization.ExpiresIn,
		Scope:       authorization.Scope,
		RedirectUri: authorization.RedirectURI,
		State:       authorization.State,
		CreatedAt:   authorization.CreatedAt,
		UserData:    authorization,
	}

	return ad, nil
}

// RemoveAuthorize revokes or deletes the authorization code.
func (s *Storage) RemoveAuthorize(code string) error {
	return data.DestroyAuthorizationWithCode(s.db, code)
}

// SaveAccess writes AccessData.
// If RefreshToken is not blank, it must save in a way that can be loaded using LoadRefresh.
func (s *Storage) SaveAccess(ad *osin.AccessData) error {
	application := ad.Client.UserData.(*data.Application)

	authorization := ad.AuthorizeData.UserData.(*data.Authorization)

	access_token := &data.AccessToken{
		IdentityId:    authorization.IdentityId,
		ApplicationId: application.Id,
		AccessToken:   ad.AccessToken,
		RefreshToken:  sql.NullString{String: ad.RefreshToken},
		ExpiresIn:     ad.ExpiresIn,
		Scope:         ad.Scope,
		RedirectURI:   ad.RedirectUri,
	}

	return data.CreateAccessToken(s.db, access_token)
}

// LoadAccess retrieves access data by token. Client information MUST be loaded together.
// AuthorizeData and AccessData DON'T NEED to be loaded if not easily available.
// Optionally can return error if expired.
func (s *Storage) LoadAccess(token string) (*osin.AccessData, error) {
	access_token, err := data.GetAccessTokenWithAccessToken(s.db, token)
	if err != nil {
		return nil, err
	}

	if access_token == nil {
		return nil, fmt.Errorf("no such access token")
	}

	client, err := s.getClientWithId(access_token.ApplicationId)
	if err != nil {
		return nil, err
	}

	ad := &osin.AccessData{
		Client:       client,
		AccessToken:  access_token.AccessToken,
		RefreshToken: access_token.RefreshToken.String,
		ExpiresIn:    access_token.ExpiresIn,
		Scope:        access_token.Scope,
		RedirectUri:  access_token.RedirectURI,
		CreatedAt:    access_token.CreatedAt,
		UserData:     access_token,
	}

	return ad, nil
}

// RemoveAccess revokes or deletes an AccessData.
func (s *Storage) RemoveAccess(token string) error {
	return data.DestroyAccessTokenWithAccessToken(s.db, token)
}

// LoadRefresh retrieves refresh AccessData. Client information MUST be loaded together.
// AuthorizeData and AccessData DON'T NEED to be loaded if not easily available.
// Optionally can return error if expired.
func (s *Storage) LoadRefresh(token string) (*osin.AccessData, error) {
	access_token, err := data.GetAccessTokenWithRefreshToken(s.db, token)
	if err != nil {
		return nil, err
	}

	if access_token == nil {
		return nil, fmt.Errorf("no such refresh token")
	}

	client, err := s.getClientWithId(access_token.ApplicationId)
	if err != nil {
		return nil, err
	}

	ad := &osin.AccessData{
		Client:       client,
		AccessToken:  access_token.AccessToken,
		RefreshToken: access_token.RefreshToken.String,
		ExpiresIn:    access_token.ExpiresIn,
		Scope:        access_token.Scope,
		RedirectUri:  access_token.RedirectURI,
		CreatedAt:    access_token.CreatedAt,
		UserData:     access_token,
	}

	return ad, nil
}

// RemoveRefresh revokes or deletes refresh AccessData.
func (s *Storage) RemoveRefresh(token string) error {
	return data.DestroyAccessTokenWithRefreshToken(s.db, token)
}
