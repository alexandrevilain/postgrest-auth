package mail

import (
	"github.com/alexandrevilain/postgrest-auth/pkg/config"
	gomail "gopkg.in/gomail.v2"
)

// sender is the config for sending emails
type sender struct {
	from   string
	dialer *gomail.Dialer
}

// NewSender creates a new instance of sender to manage mail sending of the app
func newSender(config *config.Email) *sender {
	return &sender{
		from:   config.From,
		dialer: gomail.NewPlainDialer(config.Host, config.Port, config.Auth.User, config.Auth.Pass),
	}
}

// sendEmail send email using the config of the sender struct
func (s *sender) sendEmail(to, subject, content string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", content)

	return s.dialer.DialAndSend(m)
}
