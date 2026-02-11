// (c) Magic Mango and individual authors
// SPDX-License-Identifier: Apache-2.0

package usecase

import (
	"fmt"

	"github.com/TheMagicMango/mangomail/internal/domain/event"
	"github.com/TheMagicMango/mangomail/pkg/events"
	"github.com/resend/resend-go/v2"
)

type SendEmailInputDTO struct {
	To          string   `json:"to" validate:"required"`
	From        string   `json:"from" validate:"required"`
	Html        string   `json:"html" validate:"required"`
	Subject     string   `json:"subject" validate:"required"`
	ReplyTo     string   `json:"reply_to"`
	Attachments []string `json:"attachments"`
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
	var attachments []*resend.Attachment
	if len(input.Attachments) > 0 {
		attachments = make([]*resend.Attachment, len(input.Attachments))
		for i, url := range input.Attachments {
			attachments[i] = &resend.Attachment{Path: url}
		}
	}

	email := &resend.SendEmailRequest{
		From:        input.From,
		To:          []string{input.To},
		ReplyTo:     input.ReplyTo,
		Subject:     input.Subject,
		Html:        input.Html,
		Attachments: attachments,
	}

	emailSentEvent := event.NewEmailSent()
	emailSentEvent.SetPayload(email)

	errs := uc.EventDispatcher.Dispatch(emailSentEvent)
	if len(errs) > 0 {
		return fmt.Errorf("failed to dispatch email event: %w", errs[0])
	}

	return nil
}
