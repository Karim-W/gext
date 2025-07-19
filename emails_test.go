//
//  emails_test.go
//  gext
//
//  Created by karim-w on 16/07/2025.
//

package gext_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/karim-w/gext"
)

func TestEmails(t *testing.T) {
	e := "sample@example.com"

	email, err := gext.NewEmail(e)
	if err != nil {
		t.Fatalf("Failed to create email: %v", err)
	}

	if email.String() != "45e67126a4c44c6ae030279e21437c79" {
		t.Errorf("Expected email key to be '45e67126a4c44c6ae030279e21437c79', got '%s'", email.String())
	}

	if email.Key() != "45e67126a4c44c6ae030279e21437c79" {
		t.Errorf("Expected email key to be '45e67126a4c44c6ae030279e21437c79', got '%s'", email.Key())
	}

	if email.IsEmpty() {
		t.Error("Expected email to not be empty")
	}

	fmt.Println("Email:", email)
}

func TestEmailUnparse(t *testing.T) {
	e := "BAD_EMAIL"
	_, err := gext.NewEmail(e)
	if err == nil {
		t.Fatalf("Expected error for invalid email, got none")
	}

	if err != gext.Err_FailedToParseEmail {
		t.Errorf("Expected error to be '%v', got '%v'", gext.Err_FailedToParseEmail, err)
	}

	t.Logf("Expected error for invalid email: %v", err)
}

func TestEmailMarshalJSON(t *testing.T) {
	e, err := gext.NewEmail("sample@example.com")
	if err != nil {
		t.Fatalf("Failed to create email: %v", err)
	}

	bytes, err := json.Marshal(e)
	if err != nil {
		t.Fatalf("Failed to marshal email: %v", err)
	}

	var result string
	if err := json.Unmarshal(bytes, &result); err != nil {
		t.Fatalf("Failed to unmarshal email: %v", err)
	}

	t.Logf("Marshalled email: %s", result)
}

func TestEmailUnmarshalJSON(t *testing.T) {
	byts := []byte(`{
	"email":"sample@example.com"
	}`)

	var e struct {
		Email gext.Email `json:"email"`
	}

	if err := json.Unmarshal(byts, &e); err != nil {
		t.Fatalf("Failed to unmarshal email: %v", err)
	}

	t.Logf("Unmarshalled email: %s", e.Email)

	if e.Email.String() != "45e67126a4c44c6ae030279e21437c79" {
		t.Errorf("Expected email key to be '45e67126a4c44c6ae030279e21437c79', got '%s'", e.Email.String())
	}
}

func TestEmailUnmarshalJSONInvalid(t *testing.T) {
	byts := []byte(`{
	"email":"BAD_EMAIL"
	}`)

	var e struct {
		Email gext.Email `json:"email"`
	}

	err := json.Unmarshal(byts, &e)
	if err == nil {
		t.Fatalf("Expected error for invalid email, got none")
	}

	if err != gext.Err_FailedToParseEmail {
		t.Errorf("Expected error to be '%v', got '%v'", gext.Err_FailedToParseEmail, err)
	}

	t.Logf("Expected error for invalid email: %v", err)
}

func TestEmailUnmarshalJSON_Unmarshall_err(t *testing.T) {
	byts := []byte(`{
	"email":123
	}`)

	var e struct {
		Email gext.Email `json:"email"`
	}

	err := json.Unmarshal(byts, &e)
	if err == nil {
		t.Fatalf("Expected error for invalid email, got none")
	}

	if err != gext.Err_FailedToParseEmail {
		t.Errorf("Expected error to be '%v', got '%v'", gext.Err_FailedToParseEmail, err)
	}

	t.Logf("Expected error for invalid email: %v", err)
}

func TestEqual(t *testing.T) {
	email_1, err := gext.NewEmail("sample@example.com")
	if err != nil {
		t.Fatalf("Failed to create email: %v", err)
	}

	email_2, err := gext.NewEmail("sample@example.com")
	if err != nil {
		t.Fatalf("Failed to create email: %v", err)
	}

	if !email_1.Equal(email_2) {
		t.Errorf("Expected emails to be equal, but they are not")
	}

	t.Logf("Emails are equal: %s == %s", email_1.String(), email_2.String())
}

func TestEmail_Scan(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		wantErr   bool
		errString string
	}{
		{
			name:    "valid string input",
			input:   "user@example.com",
			wantErr: false,
		},
		{
			name:    "valid []byte input",
			input:   []byte("user@example.com"),
			wantErr: false,
		},
		{
			name:      "invalid []byte input",
			input:     []byte("not-an-email"),
			wantErr:   true,
			errString: "failed to parse email",
		},
		{
			name:      "invalid type input",
			input:     123,
			wantErr:   true,
			errString: "failed to scan email",
		},
		{
			name:      "invalid email format",
			input:     "not-an-email",
			wantErr:   true,
			errString: "failed to parse email",
		},
		{
			name:      "nil input",
			input:     nil,
			wantErr:   false, // Expecting no error, should set to empty TestEmail_Scan
			errString: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var e gext.Email
			err := e.Scan(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				} else if !strings.Contains(err.Error(), tt.errString) {
					t.Errorf("expected error containing %q, got %v", tt.errString, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestEmail_Value(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		email   string
		want    interface{}
		wantErr bool
	}{
		{
			name:    "valid email",
			email:   "example@example.com",
			want:    "example@example.com",
			wantErr: false,
		},
		{
			name:    "empty email",
			email:   "",
			want:    nil, // Expecting nil for empty email
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, _ := gext.NewEmail(tt.email)
			got, gotErr := e.Value()
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Value() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Value() succeeded unexpectedly")
			}
			if tt.want == nil && got == nil {
				return // Both are nil, which is expected
			}

			if !strings.EqualFold(got.(string), tt.want.(string)) {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}
