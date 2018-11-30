package models

import (
	"database/sql"
	"time"

	"github.com/dchest/passwordreset"
	jwt "github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user of our auth system
type User struct {
	ID                 string `json:"id"`
	Email              string `json:"email"`
	Password           string `json:"password"`
	Confirmed          bool
	ConfirmToken       sql.NullString
	ResetPasswordToken sql.NullString
}

// FindByEmail allows us to find a user by its email (used for authentication)
func (u *User) FindByEmail(db *sql.DB) error {
	return db.QueryRow("SELECT id, password, confirmed, confirmToken, resetPasswordToken FROM auth.users WHERE email = $1", u.Email).Scan(&u.ID, &u.Password, &u.Confirmed, &u.ConfirmToken, &u.ResetPasswordToken)
}

// FindByID allows us to find a user by its id (used for authentication)
func (u *User) FindByID(db *sql.DB) error {
	return db.QueryRow("SELECT email, password, confirmed, confirmToken, resetPasswordToken FROM auth.users WHERE id = $1", u.ID).Scan(&u.Email, &u.Password, &u.Confirmed, &u.ConfirmToken, &u.ResetPasswordToken)
}

// Create allow us to create new user in database
func (u *User) Create(db *sql.DB) error {
	u.ID = uuid.NewV4().String()
	u.ConfirmToken = sql.NullString{String: uuid.NewV4().String(), Valid: true}
	return db.QueryRow("INSERT INTO auth.users(id, email, password, confirmToken) VALUES($1, $2, $3, $4) RETURNING id", u.ID, u.Email, u.Password, u.ConfirmToken).Scan(&u.ID)
}

// CreateResetToken create a random reset token
func (u *User) CreateResetToken(db *sql.DB, secret string) error {
	u.ResetPasswordToken = sql.NullString{String: passwordreset.NewToken(u.Email, 1*time.Hour, []byte(u.Password), []byte(secret)), Valid: true}
	_, err := db.Query("UPDATE auth.users SET resetPasswordToken = $1 WHERE id = $2", u.ResetPasswordToken, u.ID)
	return err
}

// ValidateResetToken is used to validate a reset token
func (u *User) ValidateResetToken(db *sql.DB, secret, token string) error {
	// The library used need a function to get the password from the user's login (here email)
	pwdval := func(email string) ([]byte, error) {
		u.Email = email
		err := u.FindByEmail(db)
		if err != nil {
			return nil, err
		}
		return []byte(u.Password), nil
	}
	_, err := passwordreset.VerifyToken(token, pwdval, []byte(secret))
	return err
}

// UpdatePassword edits the user's password in memory and in db
func (u *User) UpdatePassword(db *sql.DB, password string) error {
	u.Password = password
	u.HashPassword()
	_, err := db.Query("UPDATE auth.users SET password = $1 WHERE id = $2 ", u.Password, u.ID)
	return err
}

// UpdateStatus edits the user's confirmation status
func (u *User) UpdateStatus(db *sql.DB, confirm bool) error {
	u.Confirmed = confirm
	u.ConfirmToken = sql.NullString{Valid: false}
	_, err := db.Query("UPDATE auth.users SET confirmed = $1, confirmToken = $2 WHERE id = $3", u.Confirmed, u.ConfirmToken, u.ID)
	return err
}

// HashPassword hashes the user's password using bcrypt
func (u *User) HashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}

// CheckPassword checks if the provided password matches the actual user's password hash
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// CreateJWTToken creates a new JWT token for the user
func (u *User) CreateJWTToken(role, secret string, exp int) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	// Create a map to store our claims
	claims := token.Claims.(jwt.MapClaims)
	claims["userid"] = u.ID
	claims["email"] = u.Email
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(exp)).Unix()

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// GetMapRepresentation return the json representation of the user without secret informations
func (u *User) GetMapRepresentation() map[string]interface{} {
	return map[string]interface{}{
		"id":    u.ID,
		"email": u.Email,
	}
}
