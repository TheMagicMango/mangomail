// (c) Magic Mango and individual authors
// SPDX-License-Identifier: Apache-2.0

package web

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	_ "github.com/TheMagicMango/mangomail/api"
	"github.com/TheMagicMango/mangomail/configs"
	"github.com/TheMagicMango/mangomail/internal/domain/event"
	event_handler "github.com/TheMagicMango/mangomail/internal/domain/event/handler"
	"github.com/TheMagicMango/mangomail/internal/infra/web/handler"
	"github.com/TheMagicMango/mangomail/internal/infra/web/middleware"
	"github.com/TheMagicMango/mangomail/internal/usecase"
	"github.com/TheMagicMango/mangomail/pkg/events"
	"github.com/TheMagicMango/mangomail/pkg/service"
	"github.com/resend/resend-go/v2"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Service struct {
	service.Service
	HTTPServerFunc  func() error
	eventDispatcher events.EventDispatcherInterface
}

type CreateInfo struct {
	service.CreateInfo
	Config          configs.MangomailConfig
	EventDispatcher events.EventDispatcherInterface
}

func Create(ctx context.Context, createInfo *CreateInfo) (*Service, error) {
	var err error
	if err = ctx.Err(); err != nil {
		return nil, err
	}

	s := &Service{}
	createInfo.Impl = s

	err = service.Create(ctx, &createInfo.CreateInfo, &s.Service)
	if err != nil {
		return nil, err
	}

	s.eventDispatcher = createInfo.EventDispatcher
	if s.eventDispatcher == nil {
		return nil, fmt.Errorf("event dispatcher on web service create is nil")
	}

	resendClient := resend.NewClient(createInfo.Config.ResendApiKey.Value)
	emailSentEventHandler := event_handler.NewEmailSentHandler(resendClient)
	if err := s.eventDispatcher.Register(event.NewEmailSent().GetName(), emailSentEventHandler); err != nil {
		return nil, fmt.Errorf("failed to register email handler: %w", err)
	}

	// Setup use cases and handlers
	sendEmailUseCase := usecase.NewSendEmailUseCase(s.eventDispatcher)
	emailHandler := handler.NewEmailHandler(sendEmailUseCase)

	// Setup HTTP router
	r := http.NewServeMux()

	// Swagger documentation
	r.Handle("GET /api/v1/docs/", httpSwagger.WrapHandler)

	// Protected routes with authentication
	r.Handle("POST /api/v1/send", middleware.RequireApiKey(createInfo.Config.ApiKey.Value)(
		http.HandlerFunc(emailHandler.SendEmail),
	))

	// Apply CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-API-Key"},
		AllowCredentials: false,
	})
	handler := c.Handler(r)

	// Create HTTP server
	server := &http.Server{
		Addr:         createInfo.Config.HttpAddress,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
		ErrorLog:     slog.NewLogLogger(s.Logger.Handler(), slog.LevelError),
	}

	s.HTTPServerFunc = func() error {
		maxRetries := createInfo.Config.HttpMaxRetries
		retryInterval := createInfo.Config.HttpRetryInterval
		s.Logger.Info("Listening", "address", createInfo.Config.HttpAddress)
		var err error = nil
		for retry := 0; retry <= maxRetries; retry++ {
			switch err = server.ListenAndServe(); err {
			case http.ErrServerClosed:
				return nil
			default:
				s.Logger.Error("http",
					"error", err,
					"try", retry+1,
					"maxRetries", maxRetries,
					"error", err)
			}
			time.Sleep(retryInterval)
		}
		return err
	}

	return s, nil
}

func (s *Service) Alive() bool     { return true }
func (s *Service) Ready() bool     { return true }
func (s *Service) Reload() []error { return nil }
func (s *Service) Tick() []error   { return nil }

func (s *Service) Serve() error {
	go s.HTTPServerFunc()
	return s.Service.Serve()
}

func (s *Service) String() string {
	return s.Name
}

func (s *Service) Stop(_ bool) []error {
	return nil
}
