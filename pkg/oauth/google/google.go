package google

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

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

// New init provider with the googleProvider struct
func New() oauth.Provider {
	return &googleProvider{}
}

// GetUserInfo retrive google user infos based on the token provided
func (provider *googleProvider) GetUserInfo(payload *oauth.Oauth2Payload, oauthStateString string) (model.User, error) {
	var googleUser Googleuser
	var user model.User
	if payload.State != oauthStateString {
		return user, fmt.Errorf("invalid oauth state")
	}
	response, err := http.Get(fmt.Sprintf("https://www.googleapis.com/oauth2/v2/userinfo?access_token=%v", payload.Token))
	if err != nil {
		return user, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return user, fmt.Errorf("failed reading response body: %s", err.Error())
	}

	if err := json.Unmarshal(content, &googleUser); err != nil {
		return user, fmt.Errorf("An error occurred, maybe your haven't check the right scopes  %s", err.Error())
	}

	user.Email = googleUser.Email
	user.Confirmed = googleUser.VerifiedEmail

	return user, nil
}
