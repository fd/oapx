package oauth

import (
	"github.com/RangelReale/osin"
)

var config = func() *osin.ServerConfig {
	c := osin.NewServerConfig()
	c.ErrorStatusCode = 401
	c.AllowClientSecretInParams = true
	c.AllowGetAccessRequest = true
	c.AllowedAccessTypes = osin.AllowedAccessType{
		osin.AUTHORIZATION_CODE,
		osin.REFRESH_TOKEN,
	}
	return c
}()
