package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"

	"github.com/pteich/gosea/api"
	"github.com/pteich/gosea/posts"
	"github.com/pteich/gosea/status"
)

func main() {
	var err error

	// initialize logger
	logfile, err := os.Create("messages.log")
	if err != nil {
		log.Fatalf("error opening log file: %s", err.Error())
	}
	defer func() {
		log.Print("closing log file")
		logfile.Close()
	}()
	logger := log.New(os.Stdout, "gosea ", log.LstdFlags)

	// init signal handling
	sigChan := make(chan os.Signal)
	defer close(sigChan)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// create services
	postsService := posts.NewWithSEA()
	apiService := api.New(postsService)

	chiRouter := chi.NewRouter()
	chiRouter.Get("/health", status.Health)
	chiRouter.Get("/api", apiService.Posts)

	srv := &http.Server{
		Addr:    ":8000",
		Handler: chiRouter,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("error starting server: %s", err.Error())
		}
	}()

	logger.Print("starting service")

	<-sigChan

	srv.Close()

	logger.Print("stopping service")
}
