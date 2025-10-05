package lib

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"dev-go-apis/internal/views/emails"

	"github.com/resend/resend-go/v2"
)

var (
	client = resend.NewClient(RESEND_API_KEY)
)

func SendWelcomeEmail(email, token, name string) (bool, error) {
	ctx := context.Background()

	var buf bytes.Buffer
	err := emails.Welcome(name, fmt.Sprintf("%s/verify?token=%s", CLIENT_URL, token)).Render(context.Background(), &buf)
	if err != nil {
		return false, err
	}

	params := &resend.SendEmailRequest{
		From:    "Acus Dev <support@mail.acusdev.com>",
		To:      []string{email},
		Subject: "Welcome new user, please verify your email address",
		Html:    buf.String(),
	}

	sent, err := client.Emails.SendWithContext(ctx, params)
	if err != nil {
		return false, fmt.Errorf("error sending email: %s", err.Error())
	}

	log.Printf("Email sent ID: %s\n", sent.Id)

	return true, nil
}
