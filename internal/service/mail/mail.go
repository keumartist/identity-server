package mail

type MailServiceImpl struct{}

func NewMailService() MailService {
	return &MailServiceImpl{}
}

func (s *MailServiceImpl) SendVerificationEmail(email, verificationCode string) error {
	return nil
}
