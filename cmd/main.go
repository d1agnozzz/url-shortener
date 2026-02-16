package main

import (
	// "context"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/d1agnozzz/url-shortener/internal/aliaser"
	"github.com/d1agnozzz/url-shortener/internal/api"
	"github.com/d1agnozzz/url-shortener/internal/storage"
	"github.com/d1agnozzz/url-shortener/internal/urlsanitizer"
)

func main() {
	var (
		ctx         = context.Background()
		dbUrl       = os.Getenv("DATABASE_URL")
		port        = os.Getenv("PORT")
		storageType = os.Getenv("STORAGE_TYPE")
		store       storage.Storage
	)

	if dbUrl == "" {
		log.Fatalf("DB url env is not set")
		dbUrl = "postgres://postgresuser:password@localhost:5432/urlshortener_db?sslmode=disable"
	}

	if port == "" {
		port = "8080"
	}

	switch storageType {
	case "inmemory":
		store = storage.NewInMemoryStorage()
		log.Println("in-memory storage is set up")

	default:
		s, err := storage.NewPostgresStorage(ctx, dbUrl)
		if err != nil {
			log.Println("DB connection failed, falling back to in-memory storage")
			store = storage.NewInMemoryStorage()
		} else {
			store = s
			log.Println("DB connection succeeded")
		}

	}

	log.Println("URL Shortener starting...")

	addr := ":" + port
	fmt.Printf("Server listening on %s\n", addr)

	apiServer := api.NewApiServer(addr, api.APIConfig{
		Storage:       store,
		Aliaser:       aliaser.NewMd5Aliaser(),
		UrlSanitizer:  *urlsanitizer.NewUrlSanitizer(),
		MaxCollisions: 5,
	})
	http.HandleFunc("/", apiServer.UrlHandler)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
