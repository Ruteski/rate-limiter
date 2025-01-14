package main

import (
	"fmt"
	"log"
	"net/http"
	"rate-limiter/config"
	"rate-limiter/limiter"
	mid "rate-limiter/middleware"
	"rate-limiter/store"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	// carrega configurações do arquivo .env
	configs, err := config.LoadConfig(".")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	// seleciona o tipo de store do limiter
	store, err := store.NewRedisStore(configs.RedisAddr)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
		return
	}
	//store := store.NewInMemoryStore() // Alternativa: usar In-Memory que ja está implementado

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.WithValue("configs", configs))
	r.Use(middleware.WithValue("limiter", limiter.NewRateLimiter(store)))
	r.Use(mid.LogRequest)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!\n\n"))
	})

	fmt.Println("Server is running on port 8080")
	fmt.Println("")
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Println("Error: ", err)
	}
}
