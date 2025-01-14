package limiter

import (
	"log"
	"net"
	"net/http"
	"rate-limiter/config"
	"rate-limiter/store"
	"time"
)

type RateLimiter struct {
	store store.RateLimitStore
}

func NewRateLimiter(store store.RateLimitStore) *RateLimiter {
	return &RateLimiter{store: store}
}

func (r *RateLimiter) Limiter(req *http.Request, w http.ResponseWriter) bool {
	config := req.Context().Value("configs").(*config.Conf)
	ip := getIPv4(req.RemoteAddr)
	apiKey := req.Header.Get("API_KEY")

	// Verifica se o limite por API_KEY está habilitado e se a API_KEY foi fornecida
	if config.LimitPerAPI_KEY && apiKey != "" {
		blocked, err := r.store.IsBlocked(req.Context(), apiKey)
		if err != nil {
			log.Printf("Error checking if API_KEY is blocked: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return false
		}
		if blocked {
			log.Printf("API_KEY blocked: %s", apiKey)
			w.WriteHeader(config.HPPTCodeLimitReached)
			w.Write([]byte(config.MessageLimitReached))
			return false
		}

		count, err := r.store.Increment(req.Context(), apiKey)
		if err != nil {
			log.Printf("Error incrementing API_KEY counter: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return false
		}

		if count > config.MaxRequestsPerSecond {
			err := r.store.Block(req.Context(), apiKey, time.Duration(config.BlockTime)*time.Second)
			if err != nil {
				log.Printf("Error blocking API_KEY: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return false
			}

			log.Printf("API_KEY rate limit exceeded: %s", apiKey)
			w.WriteHeader(config.HPPTCodeLimitReached)
			w.Write([]byte(config.MessageLimitReached))
			return false
		}
	}

	// Verifica se o limite por IP está habilitado
	if config.LimitPerIP {
		blocked, err := r.store.IsBlocked(req.Context(), ip)
		if err != nil {
			log.Printf("Error checking if IP is blocked: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return false
		}
		if blocked {
			log.Printf("IP blocked: %s", ip)
			w.WriteHeader(config.HPPTCodeLimitReached)
			w.Write([]byte(config.MessageLimitReached))
			return false
		}

		count, err := r.store.Increment(req.Context(), ip)
		if err != nil {
			log.Printf("Error incrementing IP counter: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return false
		}

		if count > config.MaxRequestsPerSecond {
			err := r.store.Block(req.Context(), ip, time.Duration(config.BlockTime)*time.Second)
			if err != nil {
				log.Printf("Error blocking IP: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return false
			}

			log.Printf("IP rate limit exceeded: %s", ip)
			w.WriteHeader(config.HPPTCodeLimitReached)
			w.Write([]byte(config.MessageLimitReached))
			return false
		}
	}

	return true
}

func getIPv4(remoteAddr string) string {
	ip, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return remoteAddr
	}

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return remoteAddr
	}

	if parsedIP.To4() == nil && parsedIP.IsLoopback() {
		return "127.0.0.1"
	}

	return parsedIP.String()
}
