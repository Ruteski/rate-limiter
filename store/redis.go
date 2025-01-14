package store

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(addr string) (*RedisStore, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // Senha, se houver
		DB:       0,  // Banco de dados padrão
	})

	// Testa a conexão com o Redis
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return &RedisStore{client: client}, nil
}

func (r *RedisStore) Increment(ctx context.Context, key string) (int, error) {
	count, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	// Define a expiração da chave para 1 segundo (reset do contador)
	if count == 1 {
		r.client.Expire(ctx, key, time.Second).Result()
	}

	return int(count), nil
}

func (r *RedisStore) IsBlocked(ctx context.Context, key string) (bool, error) {
	exists, err := r.client.Exists(ctx, "blocked:"+key).Result()
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}

func (r *RedisStore) Block(ctx context.Context, key string, blockTime time.Duration) error {
	err := r.client.Set(ctx, "blocked:"+key, "1", blockTime).Err()
	return err
}
