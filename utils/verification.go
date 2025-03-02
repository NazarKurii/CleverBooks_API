package utils

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/resend/resend-go/v2"
)

var apiKey = "133fa2956129aa76a13914f1d628efcb-9c3f0c68-2dad457b"
var domain = "cleverbooks.shop"
var apiURL = "https://api.textlocal.in/send/"
var sender = "cleverbooks@shop.com"
var subject = "Verification code"

func SendVerificationCode(recipientEmail string) (string, error) {

	verificationCode := strconv.Itoa(rand.Intn(89999) + 10000)

	apiKey := "re_fNMzU1ed_GKdb6rBBWu3Gokfpu2d1ooR5"

	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From: "Acme <onboarding@resend.dev>",
		To:   []string{recipientEmail},
		Html: fmt.Sprintf(`
						Hello from PumpkinBooks!
						<br> Verification code: <br>
						<strong>%s</strong>
					`, verificationCode),
		Subject: "Verification PumkinBooks",
		Cc:      []string{"cc@example.com"},
		Bcc:     []string{"bcc@example.com"},
		ReplyTo: "replyto@example.com",
	}

	_, err := client.Emails.Send(params)
	if err != nil {
		return "", err
	}

	return "00000", nil
}
