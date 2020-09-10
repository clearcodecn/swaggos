package swaggos

import "github.com/go-openapi/spec"

// Oauth2 create a oauth2 header
func (y *Swaggos) Oauth2Password(tokenURL string, scopes []string) *Swaggos {
	return y.Oauth2Config(Oauth2Config{
		Flow:     PasswordFlow,
		TokenURL: tokenURL,
		Scopes:   scopes,
	})
}

func (y *Swaggos) Oauth2Implicit(authURL string, scopes []string) *Swaggos {
	return y.Oauth2Config(Oauth2Config{
		AuthorizationUrl: authURL,
		Flow:             ImplicitFlow,
		Scopes:           scopes,
	})
}

func (y *Swaggos) Oauth2Client(tokURL string, scopes []string) *Swaggos {
	return y.Oauth2Config(Oauth2Config{
		TokenURL: tokURL,
		Flow:             ApplicationFlow,
		Scopes:           scopes,
	})
}

func (y *Swaggos) Oauth2AccessCode(authURL string, tokURL string, scopes []string) *Swaggos {
	return y.Oauth2Config(Oauth2Config{
		TokenURL:         tokURL,
		AuthorizationUrl: authURL,
		Flow:             AccessCodeFlow,
		Scopes:           scopes,
	})
}

type Oauth2Flow string

const (
	AccessCodeFlow  Oauth2Flow = "accessCode"
	ImplicitFlow               = "implicit"
	PasswordFlow               = "password"
	ApplicationFlow            = "application"
)

type Oauth2Config struct {
	Flow             Oauth2Flow
	AuthorizationUrl string
	TokenURL         string
	Scopes           []string
}

func (y *Swaggos) Oauth2Config(config Oauth2Config) *Swaggos {
	var schema *spec.SecurityScheme
	switch config.Flow {
	case AccessCodeFlow:
		if config.AuthorizationUrl == "" || config.TokenURL == "" {
			panic("AuthorizationUrl or TokenURL is empty")
		}
		schema = spec.OAuth2AccessToken(config.AuthorizationUrl, config.TokenURL)
	case ImplicitFlow:
		if config.AuthorizationUrl == "" {
			panic("AuthorizationUrl or TokenURL is empty")
		}
		schema = spec.OAuth2Implicit(config.AuthorizationUrl)
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
	}

	return y.addAuth("Oauth2", schema, map[string][]string{
		"Oauth2": config.Scopes,
	})
}
