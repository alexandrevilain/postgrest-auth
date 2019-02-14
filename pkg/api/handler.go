package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/alexandrevilain/postgrest-auth/pkg/oauth"

	"github.com/alexandrevilain/postgrest-auth/pkg/config"
	"github.com/alexandrevilain/postgrest-auth/pkg/mail"
	"github.com/alexandrevilain/postgrest-auth/pkg/model"
	"github.com/alexandrevilain/postgrest-auth/pkg/oauth/facebook"
	"github.com/alexandrevilain/postgrest-auth/pkg/oauth/google"
	"github.com/labstack/echo"
)

type handler struct {
	db         *sql.DB
	config     *config.Config
	emailQueue chan mail.EmailSendRequest
	emails     *mail.EmailGenerator
}

func (h *handler) signin(c echo.Context) error {
	var user model.User
	if err := c.Bind(&user); err != nil {
		return err
	}
	submittedPassword := user.Password
	err := user.FindByEmail(h.db)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Unable to find your account")
	}
	ok := user.CheckPassword(submittedPassword)
	if !ok {
		return echo.NewHTTPError(http.StatusNotFound, "Unable to find your account")
	}
	// Check for email confirmation
	if !user.Confirmed {
		return echo.NewHTTPError(http.StatusUnauthorized, "Please confirm your account")
	}
	jwt, err := user.CreateJWTToken(h.config.DB.Roles.User, h.config.JWT.Secret, h.config.JWT.Exp)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "An error occurred while creating your token")
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"user":  user.GetMapRepresentation(),
		"token": jwt,
	})
}

func (h *handler) signup(c echo.Context) error {
	var user model.User
	if err := c.Bind(&user); err != nil {
		return err
	}
	ok := user.CheckEmailDomain(h.config.API.AllowedDomains)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "You're not allowed to create an account with the provied email address")
	}
	if err := user.HashPassword(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "An error occurred while hashing your password")
	}

	if err := user.Create(h.db); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "An error occurred while creating your account")
	}

	token, err := user.ConfirmToken.Value()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "An error occurred while creating your account")
	}

	confirmLink := fmt.Sprintf(h.config.Links.Confirm, user.ID, token)
	email, err := h.emails.GenerateConfirmEmail(user.Email, confirmLink)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "An error occurred while creating your account")
	}

	h.emailQueue <- mail.EmailSendRequest{
		To:      user.Email,
		Title:   "Please confirm your account",
		Content: email,
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id":      user.ID,
		"success": true,
	})
}

func (h *handler) confirmAccount(c echo.Context) error {
	var user model.User
	user.ID = c.Param("id")
	if err := user.FindByID(h.db); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Unable to find your account")
	}
	token, err := user.ConfirmToken.Value()
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, "Your email confirmation token is not valid")
	}
	if token != c.QueryParam("token") {
		return echo.NewHTTPError(http.StatusForbidden, "Your email confirmation token is not valid")
	}

	if err := user.UpdateStatus(h.db, true); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "An error occurred while updating your email confirmation")
	}

	return c.JSON(http.StatusCreated, map[string]bool{
		"success": true,
	})
}

// When a user ask for passwork reset
func (h *handler) sendPasswordReset(c echo.Context) error {
	var user model.User
	if err := c.Bind(&user); err != nil {
		return err
	}
	if err := user.FindByEmail(h.db); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Unable to find your account")
	}
	if err := user.CreateResetToken(h.db, h.config.API.ResetToken); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "An error occurred while creating your reset password")
	}

	token, err := user.ResetPasswordToken.Value()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "An error occurred while creating your reset password")
	}

	resetLink := fmt.Sprintf(h.config.Links.Reset, token)
	email, err := h.emails.GenerateRestePasswordEmail(user.Email, resetLink)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "An error occurred while creating your reset password")
	}

	h.emailQueue <- mail.EmailSendRequest{
		To:      user.Email,
		Title:   "Here is your reset link",
		Content: email,
	}

	return c.JSON(http.StatusCreated, map[string]bool{
		"success": true,
	})
}

type resetRequest struct {
	Password string `json:"password"`
}

// When a user submit a password reset
// The reset token in query string and the new password in the request's body
func (h *handler) resetPassword(c echo.Context) error {
	// Get new password
	var req resetRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	var user model.User
	err := user.ValidateResetToken(h.db, h.config.API.ResetToken, c.Param("token"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Wrong reset token")
	}

	if err := user.UpdatePassword(h.db, req.Password); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "An error occurred while updating your password")
	}

	return c.JSON(http.StatusCreated, map[string]bool{
		"success": true,
	})
}

func (h *handler) signinWithProvider(c echo.Context) error {
	payload := new(oauth.Oauth2Payload)

	if err := c.Bind(payload); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "An error occurred with your payload")
	}
	provider := c.Param("provider")
	if provider == "" {
		return echo.NewHTTPError(http.StatusInternalServerError, "An error occurred with your provider route")
	}
	var p oauth.Provider
	switch provider {
	case "google":
		p = google.New()
	case "facebook":
		p = facebook.New()
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("%s provider is not supported", provider))
	}
	user, err := p.GetUserInfo(payload, h.config.OAuth2.State)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if err := user.Create(h.db); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("An error occurred while creating your account %s", err.Error()))
	}
	token, err := user.CreateJWTToken(h.config.DB.Roles.User, h.config.JWT.Secret, h.config.JWT.Exp)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("An error occurred while creating your jwt token %s", err.Error()))
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"user":  user.GetMapRepresentation(),
		"token": token,
	})

}
