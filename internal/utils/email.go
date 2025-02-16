package utils

import (
	"fmt"
	"os"

	"github.com/resend/resend-go/v2"
)

func sendEmail(subject, toEmail, htmlContent string) error {
	apiKey := os.Getenv("RESEND_API_KEY") // Replace with your actual API key
	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From:    "Acme <onboarding@resend.dev>",
		To:      []string{toEmail},
		Html:    htmlContent,
		Subject: subject,
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	fmt.Println("Email sent successfully, ID:", sent.Id)
	return nil
}
