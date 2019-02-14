package mail

import (
	"github.com/alexandrevilain/postgrest-auth/pkg/config"
	"github.com/matcornic/hermes"
)

// EmailGenerator is the struct keeping the base config of hermes
type EmailGenerator struct {
	hermes hermes.Hermes
}

// NewEmailGenerator creates the base config of hermes
func NewEmailGenerator(config *config.App) *EmailGenerator {
	h := hermes.Hermes{
		Product: hermes.Product{
			Name:      config.Name,
			Link:      config.Link,
			Logo:      config.Logo,
			Copyright: "Copyright © 2017 " + config.Name + ". All rights reserved.",
		},
	}
	h.Theme = new(hermes.Flat)
	return &EmailGenerator{
		hermes: h,
	}
}

// GenerateConfirmEmail generate a custom confirm email
func (g *EmailGenerator) GenerateConfirmEmail(fullname string, link string) (string, error) {
	email := hermes.Email{
		Body: hermes.Body{
			Name: fullname,
			Intros: []string{
				"Welcome to " + g.hermes.Product.Name + "! We're very excited to have you on board.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "To get started with " + g.hermes.Product.Name + ", please click here:",
					Button: hermes.Button{
						Text: "Confirm your account",
						Link: link,
					},
				},
			},
			Outros: []string{
				"Need help, or have questions? Just reply to this email, we'd love to help.",
			},
		},
	}
	return g.hermes.GenerateHTML(email)
}

// GenerateRestePasswordEmail generate a custom reset password email
func (g *EmailGenerator) GenerateRestePasswordEmail(fullname string, link string) (string, error) {
	email := hermes.Email{
		Body: hermes.Body{
			Name: fullname,
			Intros: []string{
				"You have received this email because a password reset request for " + g.hermes.Product.Name + " account was received.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Click the button below to reset your password:",
					Button: hermes.Button{
						Color: "#DC4D2F",
						Text:  "Reset your password",
						Link:  link,
					},
				},
			},
			Outros: []string{
				"If you did not request a password reset, no further action is required on your part.",
			},
			Signature: "Thanks",
		},
	}
	return g.hermes.GenerateHTML(email)
}
