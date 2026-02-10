package dto

type SendEmailRequest struct {
	To         string `json:"to" validate:"required,email"`
	From       string `json:"from"`
	Name       string `json:"name" validate:"max=100"`
	Subject    string `json:"subject" validate:"required,min=1,max=200"`
	Saudacao   string `json:"saudacao" validate:"max=500"`
	Body       string `json:"body" validate:"required,min=1"`
	Assinatura string `json:"assinatura" validate:"max=1000"`
}

type SendEmailResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}
