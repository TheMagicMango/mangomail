// (c) Magic Mango and individual authors
// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"encoding/json"
	"net/http"
)

func RequireApiKey(validKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.Header.Get("X-API-Key")
			if key != validKey {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(UnauthorizedError{
					Message: "Invalid API key",
				})
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
