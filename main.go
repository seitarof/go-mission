package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TechBowl-japan/go-stations/db"
	"github.com/TechBowl-japan/go-stations/handler/router"
)

func main() {
	err := realMain()
	if err != nil {
		log.Fatalln("main: failed to exit successfully, err =", err)
	}
}

func realMain() error {
	// config values
	const (
		defaultPort   = ":8080"
		defaultDBPath = ".sqlite3/todo.db"
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = defaultDBPath
	}

	// set time zone
	var err error
	time.Local, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return err
	}

	// set up sqlite3
	todoDB, err := db.NewDB(dbPath)
	if err != nil {
		return err
	}
	defer todoDB.Close()

	// set basic-auth info
	os.Setenv("BASIC_AUTH_USER_ID", "test")
	os.Setenv("BASIC_AUTH_PASSWORD", "4c716d4cf211c7b7d2f3233c941771ad0507ea5bacf93b492766aa41ae9f720d")

	// set http handlers
	mux := router.NewRouter(todoDB)

	// TODO: ここから実装を行う
	// コンテキストをセットして
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt, os.Kill)
	defer stop()

	srv := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	stop()

	go srv.ListenAndServe()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = srv.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}
