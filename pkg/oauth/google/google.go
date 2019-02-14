package google

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/alexandrevilain/postgrest-auth/pkg/config"
	"github.com/alexandrevilain/postgrest-auth/pkg/model"
	"github.com/alexandrevilain/postgrest-auth/pkg/oauth"
)

type googleProvider struct {
}

//Googleuser struct of google user
type Googleuser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Link          string `json:"link"`
	Picture       string `json:"picture"`
}

// New retrive google user info based on the token provided
func New() oauth.Provider {
	return &googleProvider{}
}

// GetGoogleUserInfo retrive google user info based on the token provided
func (provider *googleProvider) GetUserInfo(payload *oauth.Oauth2Payload, oauthStateString string, config *config.Config, db *sql.DB) (model.User, string, error) {
	var googleUser Googleuser
	if payload.State != oauthStateString {
		return model.User{}, "", fmt.Errorf("invalid oauth state")
	}
	response, err := http.Get(fmt.Sprintf("https://www.googleapis.com/oauth2/v2/userinfo?access_token=%v", payload.Token))
	if err != nil {
		return model.User{}, "", fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return model.User{}, "", fmt.Errorf("failed reading response body: %s", err.Error())
	}

	if err := json.Unmarshal(content, &googleUser); err != nil {
		return model.User{}, "", fmt.Errorf("An error occurred, maybe your haven't check the right scopes  %s", err.Error())
	}

	userRandomPassword := oauth.GenreatePassword(12)
	user := model.User{
		Email:     googleUser.Email,
		Confirmed: googleUser.VerifiedEmail,
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
