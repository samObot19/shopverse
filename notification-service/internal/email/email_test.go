package email

import (
	"testing"
)

func TestSendEmail(t *testing.T) {
	tests := []struct {
		name     string
		to       string
		subject  string
		body     string
		wantErr  bool
	}{
		{
			name:    "Valid email",
			to:      "test@example.com",
			subject: "Test Subject",
			body:    "This is a test email.",
			wantErr: false,
		},
		{
			name:    "Invalid email",
			to:      "invalid-email",
			subject: "Test Subject",
			body:    "This is a test email.",
			wantErr: true,
		},
		{
			name:    "Empty recipient",
			to:      "",
			subject: "Test Subject",
			body:    "This is a test email.",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SendEmail(tt.to, tt.subject, tt.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}