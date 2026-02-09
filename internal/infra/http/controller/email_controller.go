package controller

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/TheMagicMango/mangomail/internal/infra/http/dto"
	"github.com/TheMagicMango/mangomail/internal/usecase"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type EmailController struct {
	SendEmailUseCase *usecase.SendEmailUseCase
	From             string
	HTMLTemplate     string
}

func NewEmailController(uc *usecase.SendEmailUseCase, from string, htmlTemplate string) *EmailController {
	return &EmailController{
		SendEmailUseCase: uc,
		From:             from,
		HTMLTemplate:     htmlTemplate,
	}
}

func (c *EmailController) SendEmail(w http.ResponseWriter, r *http.Request) {
	var req dto.SendEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "Invalid JSON request body",
		})
		return
	}

	// Validate request using validator
	if err := validate.Struct(&req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorMsg := formatValidationErrors(validationErrors)
		respondJSON(w, http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   errorMsg,
		})
		return
	}

	err := c.SendEmailUseCase.Execute(usecase.SendEmailInputDTO{
		To:           req.To,
		Name:         req.Name,
		Subject:      req.Subject,
		Saudacao:     req.Saudacao,
		Body:         req.Body,
		Assinatura:   req.Assinatura,
		From:         c.From,
		HTMLTemplate: c.HTMLTemplate,
	})
	if err != nil {
		slog.Error("Failed to send email", "to", req.To, "error", err)
		respondJSON(w, http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Error:   "Failed to send email",
		})
		return
	}

	respondJSON(w, http.StatusOK, dto.SendEmailResponse{
		Success: true,
		Message: "Email sent successfully",
	})
}

func formatValidationErrors(errs validator.ValidationErrors) string {
	if len(errs) == 0 {
		return "Validation failed"
	}

	// Return the first error in a user-friendly format
	err := errs[0]
	field := err.Field()

	switch err.Tag() {
	case "required":
		return "Field '" + field + "' is required"
	case "email":
		return "Field '" + field + "' must be a valid email address"
	case "min":
		return "Field '" + field + "' is too short"
	case "max":
		return "Field '" + field + "' is too long"
	default:
		return "Field '" + field + "' is invalid"
	}
}

func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("Failed to encode JSON response", "error", err)
	}
}
