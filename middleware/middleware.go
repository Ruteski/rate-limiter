package middleware

import (
	"net/http"
	"rate-limiter/limiter"
)

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		limiter := req.Context().Value("limiter").(*limiter.RateLimiter)
		// caso não tenha alcançado o limite de requisições
		if limiter.Limiter(req, w) {
			// Chama o próximo handler
			next.ServeHTTP(w, req)
		}
	})
}
