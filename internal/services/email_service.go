package services

import (
	"fmt"
	"my-go-api/internal/utils"
)

type SendEmailVerificationParams struct {
	Name  string
	Email string
	Code  string
}
type SendPasswordResetParams struct {
	Name  string
	Email string
	Token string
}

type IEmailService interface {
	SendVerificationEmail(params SendEmailVerificationParams) error
	SendPasswordResetRequest(params SendPasswordResetParams) error
}

type emailService struct {
	appUri  string
	utility utils.IUtils
}

func NewEmailService(appUri string, utility utils.IUtils) IEmailService {
	return &emailService{
		appUri:  appUri,
		utility: utility,
	}
}

func (s *emailService) SendVerificationEmail(params SendEmailVerificationParams) error {
	var subject = "Email verification"

	var emailBody = fmt.Sprintf(`
	Hello %s.
	This is your verification code: %s`,
		params.Name, params.Code)

	err := s.utility.SendEmailWithGmail(subject, emailBody, params.Email)
	if err != nil {
		return err
	}

	return nil
}

func (s *emailService) SendPasswordResetRequest(params SendPasswordResetParams) error {
	var subject = "Reset Password"
	link := fmt.Sprintf("%s/reset-password/%s", s.appUri, params.Token)

	var emailBody = fmt.Sprintf(`
	Hello %s.
	You receive this email because you sent a request to update your password.
	You can ignore this email if you didn't.
	Please follow this link to update your password
	%s
	`,
		params.Name, link)

	err := s.utility.SendEmailWithGmail(subject, emailBody, params.Email)
	if err != nil {
		return err
	}

	return nil
}
