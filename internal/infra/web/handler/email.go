// (c) Magic Mango and individual authors
// SPDX-License-Identifier: Apache-2.0

package handler

import (
	"encoding/json"
	"github.com/TheMagicMango/mangomail/internal/usecase"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type EmailHandler struct {
	SendEmailUseCase *usecase.SendEmailUseCase
}

func NewEmailHandler(sendEmailUseCase *usecase.SendEmailUseCase) *EmailHandler {
	return &EmailHandler{
		SendEmailUseCase: sendEmailUseCase,
	}
}

// SendEmail godoc
// @Summary      Send an email
// @Description  Send a single email using the Resend API
// @Tags         email
// @Accept       json
// @Produce      json
// @Param        x-api-key  header    string                     true  "API Key for authentication"
// @Param        request    body      usecase.SendEmailInputDTO  true  "Email request payload"
// @Success      200        "Email sent successfully"
// @Failure      400        {object}  BadRequestError
// @Failure      401        {object}  middleware.UnauthorizedError
// @Failure      500        {object}  InternalServerError
// @Security     ApiKeyAuth
// @Router       /api/v1/send [post]
func (h *EmailHandler) SendEmail(w http.ResponseWriter, r *http.Request) {
	var input usecase.SendEmailInputDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(BadRequestError{
			Message: "Invalid JSON payload",
		})
		return
	}

	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ValidationError{
			Message: err.Error(),
		})
		return
	}

	err := h.SendEmailUseCase.Execute(input)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(InternalServerError{
			Message: "Failed to send email",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
}
