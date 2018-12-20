package model

import "testing"

func TestCheckEmailDomain(t *testing.T) {
	tests := []struct {
		email   string
		domains []string
		valid   bool
	}{
		{
			"alexandre@google.com", []string{}, true,
		},
		{
			"alexandre@google.com", []string{"facebook.com"}, false,
		},
		{
			"alexandre@google.com", []string{"facebook.com", "google.com"}, true,
		},
	}
	for _, test := range tests {
		u := User{Email: test.email}
		ok := u.CheckEmailDomain(test.domains)
		if ok != test.valid {
			t.Errorf("Expected %v check to be %v, got: %v", test.email, test.valid, ok)
		}
	}
}
