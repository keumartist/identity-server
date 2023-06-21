package mail

type MailService interface {
	SendVerificationEmail(email, verificationCode string) error
}
