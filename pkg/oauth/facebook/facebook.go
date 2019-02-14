package facebook

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/tarektouati/postgrest-auth/pkg/config"
	"github.com/tarektouati/postgrest-auth/pkg/model"
	"github.com/tarektouati/postgrest-auth/pkg/oauth"
)

type facebookProvider struct {
}

//Facebookuser struct of facebook user
type Facebookuser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Link          string `json:"link"`
	Picture       string `json:"picture"`
}

// New retrive facebook user info based on the token provided
func New() oauth.Provider {
	return &facebookProvider{}
}

// GetUserInfo retrive facebook user info based on the token provided
func (provider *facebookProvider) GetUserInfo(payload *oauth.Oauth2Payload, oauthStateString string, config *config.Config, db *sql.DB) (model.User, string, error) {
	var facebookUser Facebookuser
	if payload.State != oauthStateString {
		return model.User{}, "", fmt.Errorf("invalid oauth state")
	}
	response, err := http.Get(fmt.Sprintf("https://graph.facebook.com/me?fields=email&access_token=%v", payload.Token))
	if err != nil {
		return model.User{}, "", fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return model.User{}, "", fmt.Errorf("failed reading response body: %s", err.Error())
	}

	if err := json.Unmarshal(content, &facebookUser); err != nil {
		return model.User{}, "", fmt.Errorf("An error occurred, maybe your haven't check the right scopes  %s", err.Error())
	}

	userRandomPassword := oauth.GenreatePassword(12)
	user := model.User{
		Email:     facebookUser.Email,
		Confirmed: true,
		Password:  userRandomPassword,
	}
	if err := user.HashPassword(); err != nil {
		return model.User{}, "", fmt.Errorf("An error occurred while hashing your password  %s", err.Error())
	}
	token, err := oauth.CreateUser(&user, config, db)
	if err != nil {
		return model.User{}, "", err
	}

	return user, token, nil
}
