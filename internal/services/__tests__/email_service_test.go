package services_test

import (
	"errors"
	"my-go-api/internal/services"
	mockutils "my-go-api/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type EmailServiceTestSuite struct {
	suite.Suite
	ctrl      *gomock.Controller
	mockUtils *mockutils.MockIUtils
	services  services.IEmailService
}

func (suite *EmailServiceTestSuite) SetupTest() {
	appUri := "http://localhost:5000"
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockUtils = mockutils.NewMockIUtils(suite.ctrl)
	suite.services = services.NewEmailService(appUri, suite.mockUtils)
}

func (suite *EmailServiceTestSuite) TearDownTest() {
	suite.ctrl.Finish()
}

func (suite *EmailServiceTestSuite) TestSendVerificationEmail() {
	suite.Run("It should send the email", func() {
		suite.mockUtils.EXPECT().SendEmailWithGmail(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		err := suite.services.SendVerificationEmail(services.SendEmailVerificationParams{
			Name:  "ari",
			Email: "ari@mail.com",
			Code:  "8888",
		})
		assert.NoError(suite.T(), err)
	})

	suite.Run("It should fail to send the email", func() {
		suite.mockUtils.EXPECT().SendEmailWithGmail(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("error"))
		err := suite.services.SendVerificationEmail(services.SendEmailVerificationParams{
			Name:  "ari",
			Email: "ari@mail.com",
			Code:  "8888",
		})
		assert.Error(suite.T(), err)
	})
}

func TestEmailServiceTestSuite(t *testing.T) {
	suite.Run(t, new(EmailServiceTestSuite))
}
