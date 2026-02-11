// (c) Magic Mango and individual authors
// SPDX-License-Identifier: Apache-2.0

package usecase

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/TheMagicMango/mangomail/pkg/events"
	"github.com/resend/resend-go/v2"
)

type SendCampaignInputDTO struct {
	From        string   `json:"from" validate:"required"`
	Subject     string   `json:"subject" validate:"required"`
	ReplyTo     string   `json:"reply_to"`
	Attachments []string `json:"attachments"`
}

type SendCampaignUseCase struct {
	EmailSentEvent  events.EventInterface
	EventDispatcher events.EventDispatcherInterface
}

func NewSendCampaignUseCase(
	emailSentEvent events.EventInterface,
	eventDispatcher events.EventDispatcherInterface,
) *SendCampaignUseCase {
	return &SendCampaignUseCase{
		EmailSentEvent:  emailSentEvent,
		EventDispatcher: eventDispatcher,
	}
}

func (uc *SendCampaignUseCase) Execute(input SendCampaignInputDTO, htmlContent string, csvRows []map[string]interface{}, rateLimit int) error {
	startTime := time.Now()

	// Prepare attachments
	var attachments []*resend.Attachment
	if len(input.Attachments) > 0 {
		attachments = make([]*resend.Attachment, len(input.Attachments))
		for i, url := range input.Attachments {
			attachments[i] = &resend.Attachment{Path: url}
		}
	}

	// Send emails in batches
	batchSize := rateLimit
	totalRows := len(csvRows)

	for i := 0; i < totalRows; i += batchSize {
		end := i + batchSize
		if end > totalRows {
			end = totalRows
		}

		batch := csvRows[i:end]

		for _, row := range batch {
			email, ok := row["email"].(string)
			if !ok || email == "" {
				slog.Warn("Skipping row with invalid email", "row", row)
				continue
			}

			// Replace placeholders inline
			subject := input.Subject
			html := htmlContent
			for key, value := range row {
				placeholder := fmt.Sprintf("{{%s}}", key)
				valueStr := fmt.Sprint(value)
				subject = strings.ReplaceAll(subject, placeholder, valueStr)
				html = strings.ReplaceAll(html, placeholder, valueStr)
			}

			emailReq := &resend.SendEmailRequest{
				From:        input.From,
				To:          []string{email},
				Subject:     subject,
				Html:        html,
				Attachments: attachments,
			}

			if input.ReplyTo != "" {
				emailReq.ReplyTo = input.ReplyTo
			}

			uc.EmailSentEvent.SetPayload(emailReq)
			errs := uc.EventDispatcher.Dispatch(uc.EmailSentEvent)
			for _, err := range errs {
				slog.Error("Failed to send email", "email", email, "error", err)
			}
		}

		if end < totalRows {
			slog.Info("Waiting before next batch", "sent", end, "total", totalRows, "delay", time.Second)
			time.Sleep(time.Second)
		}
	}

	duration := time.Since(startTime)
	slog.Info("Campaign completed", "duration", duration)

	return nil
}
