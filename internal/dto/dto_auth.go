package dto

type Register struct {
	Name     string `json:"name" validate:"required,min=5"`
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=5"`
	Password string `json:"password" validate:"required,strongPassword"`
}

type Login struct {
	Identity string `json:"identity" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type VerifyNewAccount struct {
	Code  string `json:"code" validate:"required,min=8,max=8"`
	Token string `json:"token" validate:"required"`
}

type ForgotPassword struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPassword struct {
	Password        string `json:"password" validate:"required,strongPassword"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
	Token           string `json:"token" validate:"required"`
}

type ResendVerification struct {
	Email string `json:"email" validate:"required,email"`
}
