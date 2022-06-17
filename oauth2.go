package swaggos

import (
	"fmt"

	"github.com/go-openapi/spec"
)

// Oauth2Password setup swagger oauth2 by password
func (swaggos *Swaggos) Oauth2Password(tokenURL string, scopes []string) *Swaggos {
	return swaggos.Oauth2Config(Oauth2Config{
		Flow:     PasswordFlow,
		TokenURL: tokenURL,
		Scopes:   scopes,
	})
}

// Oauth2Implicit setup swagger oauth2 by implicit
func (swaggos *Swaggos) Oauth2Implicit(authURL string, scopes []string) *Swaggos {
	return swaggos.Oauth2Config(Oauth2Config{
		AuthorizationURL: authURL,
		Flow:             ImplicitFlow,
		Scopes:           scopes,
	})
}

// Oauth2Client setup swagger oauth2 by client
func (swaggos *Swaggos) Oauth2Client(tokURL string, scopes []string) *Swaggos {
	return swaggos.Oauth2Config(Oauth2Config{
		TokenURL: tokURL,
		Flow:     ApplicationFlow,
		Scopes:   scopes,
	})
}

// Oauth2AccessCode setup swagger oauth2 by access code
func (swaggos *Swaggos) Oauth2AccessCode(authURL string, tokURL string, scopes []string) *Swaggos {
	return swaggos.Oauth2Config(Oauth2Config{
		TokenURL:         tokURL,
		AuthorizationURL: authURL,
		Flow:             AccessCodeFlow,
		Scopes:           scopes,
	})
}

// Oauth2Flow is the type of oauth2
type Oauth2Flow string

const (
	// AccessCodeFlow accessCode flow
	AccessCodeFlow Oauth2Flow = "accessCode"
	// ImplicitFlow implicit flow
	ImplicitFlow = "implicit"
	// PasswordFlow password flow
	PasswordFlow = "password"
	// ApplicationFlow application flow
	ApplicationFlow = "application"
)

// Oauth2Config is config for oauth2
type Oauth2Config struct {
	Flow             Oauth2Flow
	AuthorizationURL string
	TokenURL         string
	Scopes           []string
}

// Oauth2Config setup oauth2 access type
func (swaggos *Swaggos) Oauth2Config(config Oauth2Config) *Swaggos {
	var schema *spec.SecurityScheme
	switch config.Flow {
	case AccessCodeFlow:
		if config.AuthorizationURL == "" || config.TokenURL == "" {
			panic("AuthorizationURL or TokenURL is empty")
		}
		schema = spec.OAuth2AccessToken(config.AuthorizationURL, config.TokenURL)
	case ImplicitFlow:
		if config.AuthorizationURL == "" {
			panic("AuthorizationURL or TokenURL is empty")
		}
		schema = spec.OAuth2Implicit(config.AuthorizationURL)
	case PasswordFlow:
		if config.TokenURL == "" {
			panic("TokenURL or TokenURL is empty")
		}
		schema = spec.OAuth2Password(config.TokenURL)
	case ApplicationFlow:
		if config.TokenURL == "" {
			panic("TokenURL or TokenURL is empty")
		}
		schema = spec.OAuth2Application(config.TokenURL)
	default:
		panic("invalid oauth2 flow")
	}
	for _, scope := range config.Scopes {
		schema.AddScope(scope, fmt.Sprintf("scope for: %s", scope))
	}
	return swaggos.addAuth("Oauth2", schema, map[string][]string{
		"Oauth2": config.Scopes,
	})
}
