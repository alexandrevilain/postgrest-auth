package oauth

// Ici tu dois faire une interface qui est générique pour tous les providers OAuth2
// Par exemple, chaque provider doit fournir la method ToModelUser
// Genre func (g *GoogleProvider) ToModelUser() (models.User, error)

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/tarektouati/postgrest-auth/pkg/config"
	"github.com/tarektouati/postgrest-auth/pkg/model"
)

// Oauth2Payload payload to retrive for google
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

// CreateUser transform google user to db user and return token
func CreateUser(user *model.User, config *config.Config, db *sql.DB) (string, error) {
	if err := user.Create(db); err != nil {
		return "", fmt.Errorf("An error occurred while creating your account %s", err.Error())
	}
	jwt, err := user.CreateJWTToken(config.DB.Roles.User, config.JWT.Secret, config.JWT.Exp)
	if err != nil {
		return "", fmt.Errorf("An error occurred while creating your jwt token %s", err.Error())
	}
	return jwt, nil
}

//Providers give you all providers functions for auth2
type Provider interface {
	GetUserInfo(payload *Oauth2Payload, oauthStateString string, config *config.Config, db *sql.DB) (model.User, string, error)
}
