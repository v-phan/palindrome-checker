package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/viptra/palindrom-ee/db"
	"github.com/viptra/palindrom-ee/handlers"
)

func main() {
	dbAdr := os.Getenv("PALINDROME_DB")
	ctx, cancel := context.WithCancel(context.Background())

	pool, err := pgxpool.Connect(ctx, dbAdr)
	if err != nil {
		log.Fatalf(err.Error())
	}

	//channel for error
	errChan := make(chan error, 1)

	dbConn := &db.DbConnection{Db: pool, Ctx: ctx}
	userHandler := handlers.NewUserHandler(dbConn)
	palindromeHandler := handlers.NewPalindromeHandler(dbConn)
	sm := http.NewServeMux()
	sm.Handle("/user", userHandler)
	sm.Handle("/palindrome", palindromeHandler)

	s := http.Server{
		Addr:         ":8080",
		Handler:      sm,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			errChan <- err
			return
		}
	}()

	log.Println("App is up and running")

	err = <-errChan
	if err != nil {
		log.Println("Recieved error... Shutting down")
		log.Printf("Error: %v", err)
		cancel()
	}
}
