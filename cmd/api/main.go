package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pruthuvifernando/url-shortner/config"
	"github.com/pruthuvifernando/url-shortner/internal/handler"
	"github.com/pruthuvifernando/url-shortner/internal/usecase"
)

func main() {
	cfg := config.Load()

	// TODO: initialise Postgres connection (database/sql + pgx driver)
	// db, err := sql.Open("pgx", cfg.DBDSN)

	// TODO: initialise Valkey client
	// valkeyClient := ...

	// Wire up layers
	// repo := postgres.NewURLRepository(db)
	// cache := valkey.NewURLCache(valkeyClient)
	uc := usecase.NewURLUsecase(nil, nil) // replace nils with real impls
	h := handler.NewURLHandler(uc)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /urls", h.ShortenURL)
	mux.HandleFunc("GET /{shortCode}", h.Redirect)
	mux.HandleFunc("DELETE /urls/{shortCode}", h.DeleteURL)

	wrapped := handler.LoggingMiddleware(handler.HeaderMiddleware(mux))

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      wrapped,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		BaseContext: func(l net.Listener) context.Context {
			return context.WithValue(context.Background(), ctxKeyAddr, l.Addr().String())
		},
	}

	go func() {
		fmt.Printf("server starting on :%s\n", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("server error: %s\n", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("forced shutdown: %s\n", err)
	}
	fmt.Println("stopped")
}

type ctxKey string

const ctxKeyAddr ctxKey = "serverAddr"
