package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"rate-limiter/config"
	"rate-limiter/limiter"
	"rate-limiter/store"
	"sync"
	"testing"
	"time"
)

// TestRateLimiterIP testa o limite de requisições por IP
func TestRateLimiterIP(t *testing.T) {
	// Configurações do rate limiter
	cfg := &config.Conf{
		LimitPerIP:           true,
		LimitPerAPI_KEY:      false,
		MaxRequestsPerSecond: 2, // Limite de 2 requisições por segundo
		BlockTime:            1, // Bloqueia por 1 segundo
		HPPTCodeLimitReached: 429,
		MessageLimitReached:  "Too many requests",
	}

	// Usa o armazenamento em memória para os testes
	store := store.NewInMemoryStore()
	rateLimiter := limiter.NewRateLimiter(store)

	// Simula requisições de um mesmo IP
	req, _ := http.NewRequest("GET", "/", nil)
	req.RemoteAddr = "127.0.0.1:12345"

	for i := 0; i < 3; i++ { // Faz 3 requisições
		w := httptest.NewRecorder()

		// Adiciona as configurações ao contexto da requisição
		ctx := context.WithValue(req.Context(), "configs", cfg)
		req = req.WithContext(ctx)

		allowed := rateLimiter.Limiter(req, w)

		if i < 2 {
			// As primeiras 2 requisições devem ser permitidas
			if !allowed {
				t.Errorf("Expected request to be allowed, but it was blocked")
			}
			if w.Code != http.StatusOK {
				t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
			}
		} else {
			// A terceira requisição deve ser bloqueada
			if allowed {
				t.Errorf("Expected request to be blocked, but it was allowed")
			}
			if w.Code != cfg.HPPTCodeLimitReached {
				t.Errorf("Expected status %d, got %d", cfg.HPPTCodeLimitReached, w.Code)
			}
		}
	}
}

// TestRateLimiterAPIKey testa o limite de requisições por API_KEY
func TestRateLimiterAPIKey(t *testing.T) {
	// Configurações do rate limiter
	cfg := &config.Conf{
		LimitPerIP:           false,
		LimitPerAPI_KEY:      true,
		MaxRequestsPerSecond: 2, // Limite de 2 requisições por segundo
		BlockTime:            1, // Bloqueia por 1 segundo
		HPPTCodeLimitReached: 429,
		MessageLimitReached:  "Too many requests",
	}

	// Usa o armazenamento em memória para os testes
	store := store.NewInMemoryStore()
	rateLimiter := limiter.NewRateLimiter(store)

	// Simula requisições com uma mesma API_KEY
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("API_KEY", "test-key")

	for i := 0; i < 3; i++ { // Faz 3 requisições
		w := httptest.NewRecorder()

		// Adiciona as configurações ao contexto da requisição
		ctx := context.WithValue(req.Context(), "configs", cfg)
		req = req.WithContext(ctx)

		allowed := rateLimiter.Limiter(req, w)

		if i < 2 {
			// As primeiras 2 requisições devem ser permitidas
			if !allowed {
				t.Errorf("Expected request to be allowed, but it was blocked")
			}
			if w.Code != http.StatusOK {
				t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
			}
		} else {
			// A terceira requisição deve ser bloqueada
			if allowed {
				t.Errorf("Expected request to be blocked, but it was allowed")
			}
			if w.Code != cfg.HPPTCodeLimitReached {
				t.Errorf("Expected status %d, got %d", cfg.HPPTCodeLimitReached, w.Code)
			}
		}
	}
}

// TestRateLimiterBlockDuration testa o tempo de bloqueio
func TestRateLimiterBlockDuration(t *testing.T) {
	// Configurações do rate limiter
	cfg := &config.Conf{
		LimitPerIP:           true,
		LimitPerAPI_KEY:      false,
		MaxRequestsPerSecond: 2, // Limite de 2 requisições por segundo
		BlockTime:            1, // Bloqueia por 1 segundo
		HPPTCodeLimitReached: 429,
		MessageLimitReached:  "Too many requests",
	}

	// Usa o armazenamento em memória para os testes
	store := store.NewInMemoryStore()
	rateLimiter := limiter.NewRateLimiter(store)

	// Simula requisições de um mesmo IP
	req, _ := http.NewRequest("GET", "/", nil)
	req.RemoteAddr = "127.0.0.1:12345"

	// Adiciona as configurações ao contexto da requisição
	ctx := context.WithValue(req.Context(), "configs", cfg)
	req = req.WithContext(ctx)

	// Faz 3 requisições para exceder o limite
	for i := 0; i < 3; i++ {
		w := httptest.NewRecorder()
		rateLimiter.Limiter(req, w)
	}

	// Aguarda o tempo de bloqueio expirar
	time.Sleep(time.Duration(cfg.BlockTime) * time.Second)

	// Faz uma nova requisição após o bloqueio expirar
	w := httptest.NewRecorder()
	allowed := rateLimiter.Limiter(req, w)

	// A requisição deve ser permitida após o bloqueio expirar
	if !allowed {
		t.Errorf("Expected request to be allowed after block time expired, but it was blocked")
	}
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

// TestRateLimiterConcurrency testa o rate limiter com requisições concorrentes
func TestRateLimiterConcurrency(t *testing.T) {
	// Configurações do rate limiter
	cfg := &config.Conf{
		LimitPerIP:           true,
		LimitPerAPI_KEY:      false,
		MaxRequestsPerSecond: 10, // Limite de 10 requisições por segundo
		BlockTime:            1,  // Bloqueia por 1 segundo
		HPPTCodeLimitReached: 429,
		MessageLimitReached:  "Too many requests",
	}

	// Usa o armazenamento em memória para os testes
	store := store.NewInMemoryStore()
	rateLimiter := limiter.NewRateLimiter(store)

	// Simula requisições concorrentes de um mesmo IP
	var wg sync.WaitGroup
	req, _ := http.NewRequest("GET", "/", nil)
	req.RemoteAddr = "127.0.0.1:12345"

	allowed := 0
	blocked := 0

	for i := 0; i < 15; i++ { // Faz 15 requisições concorrentes
		wg.Add(1)
		go func() {
			defer wg.Done()

			w := httptest.NewRecorder()

			// Adiciona as configurações ao contexto da requisição
			ctx := context.WithValue(req.Context(), "configs", cfg)
			req := req.WithContext(ctx)

			if rateLimiter.Limiter(req, w) {
				allowed++
			} else {
				blocked++
			}
		}()
	}

	wg.Wait()

	// Verifica se o número de requisições permitidas e bloqueadas está correto
	if allowed != cfg.MaxRequestsPerSecond {
		t.Errorf("Expected %d allowed requests, got %d", cfg.MaxRequestsPerSecond, allowed)
	}
	if blocked != 5 { // 15 requisições - 10 permitidas = 5 bloqueadas
		t.Errorf("Expected 5 blocked requests, got %d", blocked)
	}
}
