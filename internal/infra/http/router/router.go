package router

import (
	"encoding/json"
	"net/http"

	"github.com/TheMagicMango/mangomail/configs"
	"github.com/TheMagicMango/mangomail/internal/infra/http/controller"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"golang.org/x/time/rate"
)

func NewRouter(emailController *controller.EmailController, cfg *configs.MangomailConfig) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.CorsAllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "X-API-Key"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	// Protected API routes with rate limiting and authentication
	r.Group(func(r chi.Router) {
		limiter := rate.NewLimiter(rate.Limit(cfg.MangomailRateLimit), int(cfg.MangomailRateLimit))
		r.Use(rateLimitMiddleware(limiter))

		r.Use(apiKeyMiddleware(cfg.MangomailApiKey.Value))

		r.Post("/api/v1/send", emailController.SendEmail)
	})

	return r
}

func rateLimitMiddleware(limiter *rate.Limiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"success": false,
					"error":   "Rate limit exceeded. Please try again later.",
				})
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func apiKeyMiddleware(validKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.Header.Get("X-API-Key")
			if key == "" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"success": false,
					"error":   "Missing API key. Please provide X-API-Key header.",
				})
				return
			}

			if key != validKey {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"success": false,
					"error":   "Invalid API key.",
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
