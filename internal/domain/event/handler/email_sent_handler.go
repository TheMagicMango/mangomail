package handler

import (
	"encoding/json"
	"log/slog"
	"sync"

	"github.com/TheMagicMango/mangomail/pkg/events"
	"github.com/resend/resend-go/v2"
)

type EmailSentHandler struct {
	ResendClient *resend.Client
}

func NewEmailSentHandler(client *resend.Client) *EmailSentHandler {
	return &EmailSentHandler{
		ResendClient: client,
	}
}

func (h *EmailSentHandler) Handle(ev events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()

	payload := ev.GetPayload()

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		slog.Error("Failed to marshal email", "error", err)
		return
	}

	var email resend.SendEmailRequest
	err = json.Unmarshal(payloadBytes, &email)
	if err != nil {
		slog.Error("Failed to unmarshal email", "error", err)
		return
	}

	sent, err := h.ResendClient.Emails.Send(&email)
	if err != nil {
		slog.Error("Failed to send email",
			"to", email.To,
			"subject", email.Subject,
			"error", err)
		return
	}

	slog.Info("Email sent successfully",
		"to", email.To,
		"subject", email.Subject,
		"id", sent.Id)
}
