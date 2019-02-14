package oauth

import (
	"github.com/alexandrevilain/postgrest-auth/pkg/model"
)

// Oauth2Payload is payload struct to retrive from provider login
type Oauth2Payload struct {
	State string `json:"state"`
	Token string `json:"token"`
}

//Provider give you all providers functions for oauth2
type Provider interface {
	GetUserInfo(payload *Oauth2Payload, oauthStateString string) (model.User, error)
}
