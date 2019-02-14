package oauth

import (
	"math/rand"
	"time"

	"github.com/alexandrevilain/postgrest-auth/pkg/model"
)

// Oauth2Payload is payload struct to retrive from provider login
type Oauth2Payload struct {
	State string `json:"state"`
	Token string `json:"token"`
}
type provider struct {
}

// GenreatePassword generate random password
func GenreatePassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(
		rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

//Provider give you all providers functions for oauth2
type Provider interface {
	GetUserInfo(payload *Oauth2Payload, oauthStateString string) (model.User, error)
}
