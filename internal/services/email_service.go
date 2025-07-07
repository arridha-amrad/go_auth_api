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

type IEmailService interface {
	SendVerificationEmail(params SendEmailVerificationParams) error
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
