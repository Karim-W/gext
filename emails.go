//
//  emails.go
//  gext
//
//  Created by karim-w on 16/07/2025.
//

package gext

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/mail"
	"strings"
)

var Err_FailedToParseEmail = fmt.Errorf("failed to parse email address")

type Email struct {
	email string
	key   string
}

// NewEmail creates a new Email instance from a string
// It returns an error if the email address is invalid.
//
// example:
//
// email,err := gext.NewEmail(string_mail)
//
//	if err != nil {
//	  // handle error
//	}
func NewEmail(email string) (Email, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	_, err := mail.ParseAddress(email)
	if err != nil {
		return Email{}, Err_FailedToParseEmail
	}

	e := Email{
		email: email,
	}

	e.init()

	return e, nil
}

func (e *Email) init() {
	hash := md5.Sum([]byte(e.email))
	e.key = hex.EncodeToString(hash[:])
}

// unmarshalJSON unmarshals the email from a JSON string
func (e *Email) UnmarshalJSON(data []byte) error {
	var e_string string
	if err := json.Unmarshal(data, &e_string); err != nil {
		return Err_FailedToParseEmail
	}

	email, err := NewEmail(e_string)
	if err != nil {
		return Err_FailedToParseEmail
	}

	*e = email

	return nil
}

// MarshalJSON marshals the email to a JSON string
func (e Email) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.email + "\""), nil
}

// String returns the email as a string
func (e Email) String() string {
	return e.key
}

// Key returns the key of the email
func (e Email) Key() string {
	return e.key
}

// Equal checks if two emails are equal
func (e Email) Equal(other Email) bool {
	return e.key == other.key
}

// IsEmpty checks if the email is empty
func (e Email) IsEmpty() bool {
	return e.email == ""
}

// Scan implements the sql.Scanner interface for Email
func (e *Email) Scan(value interface{}) error {
	if value == nil {
		*e = Email{} // Set to empty email
		return nil
	}

	str, ok := value.(string)
	if ok {
		email, err := NewEmail(str)
		if err != nil {
			return err
		}

		*e = email

		return nil
	}

	val, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan email: %v", value)
	}

	email, err := NewEmail(string(val))
	if err != nil {
		return err
	}
	*e = email

	return nil
}

func (e Email) Value() (interface{}, error) {
	if e.IsEmpty() {
		return nil, nil // Return nil to represent a SQL NULL
	}

	return e.email, nil
}
