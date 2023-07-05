package mail

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type MailServiceImpl struct {
	sesService *ses.SES
	sender     string
}

func NewMailService() MailService {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION"))},
	)

	if err != nil {
		log.Println(err)
	}

	svc := ses.New(sess)

	return &MailServiceImpl{
		sesService: svc,
		sender:     os.Getenv("SENDER_EMAIL"),
	}
}

func (s *MailServiceImpl) SendVerificationEmail(email, verificationCode string) error {
	return s.sendWithSES(email, verificationCode)
}

func (s *MailServiceImpl) sendWithSES(email, verificationCode string) error {
	recipient := email
	subject := "Email verification code from artsso"
	htmlBody := fmt.Sprintf("<h1>Verification code</h1><p>%s</p>", verificationCode)
	textBody := ""
	charSet := "UTF-8"

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(recipient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(charSet),
					Data:    aws.String(htmlBody),
				},
				Text: &ses.Content{
					Charset: aws.String(charSet),
					Data:    aws.String(textBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(charSet),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(s.sender),
	}

	_, err := s.sesService.SendEmail(input)

	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("Email Sent to address: " + recipient)

	return nil
}
