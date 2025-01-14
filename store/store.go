package store

import (
	"context"
	"time"
)

// RateLimitStore define as operações necessárias para o rate limiting
type RateLimitStore interface {
	Increment(ctx context.Context, key string) (int, error)               // Incrementa o contador de requisições
	IsBlocked(ctx context.Context, key string) (bool, error)              // Verifica se uma chave está bloqueada
	Block(ctx context.Context, key string, blockTime time.Duration) error // Bloqueia uma chave por um tempo
}
