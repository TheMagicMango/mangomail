package usecase

import (
	"fmt"

	"github.com/TheMagicMango/mangomail/internal/domain/event"
	"github.com/TheMagicMango/mangomail/pkg/events"
	"github.com/resend/resend-go/v2"
)

type SendEmailInputDTO struct {
	To           string
	Name         string
	Subject      string
	Saudacao     string
	Body         string
	Assinatura   string
	From         string
	HTMLTemplate string
}

type SendEmailUseCase struct {
	EventDispatcher events.EventDispatcherInterface
}

func NewSendEmailUseCase(eventDispatcher events.EventDispatcherInterface) *SendEmailUseCase {
	return &SendEmailUseCase{
		EventDispatcher: eventDispatcher,
	}
}

func (uc *SendEmailUseCase) Execute(input SendEmailInputDTO) error {
	placeholders := map[string]interface{}{
		"name":       input.Name,
		"subject":    input.Subject,
		"saudacao":   input.Saudacao,
		"body":       input.Body,
		"assinatura": input.Assinatura,
	}

	html := replacePlaceholders(input.HTMLTemplate, placeholders)
	subject := replacePlaceholders(input.Subject, placeholders)

	emailReq := &resend.SendEmailRequest{
		From:    input.From,
		To:      []string{input.To},
		Subject: subject,
		Html:    html,
	}

	emailSentEvent := event.NewEmailSent()
	emailSentEvent.SetPayload(emailReq)

	if err := uc.EventDispatcher.Dispatch(emailSentEvent); err != nil {
		return fmt.Errorf("failed to dispatch email event: %w", err)
	}

	return nil
}
