package main

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3"

	"github.com/pteich/gosea/src/seabackend"
)

var Version = "latest"

func main() {
	flamingo.App([]dingo.Module{
		new(seabackend.Module),
	})

	/*
		var err error

		ctx, cancel := context.WithCancel(context.Background())

		// initialize logger
		logfile, err := os.Create("/tmp/messages.log")
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
		go func() {
			sig := <-sigChan
			log.Printf("received signal %s", sig.String())
			cancel()
		}()
		defer close(sigChan)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		// create services
		postsService := seabackend.NewWithSEA()
		apiService := api.New(postsService, logger)

		mux := http.NewServeMux()
		mux.HandleFunc("/health", status.Health)
		mux.HandleFunc("/api", apiService.Posts)

		srv := &http.Server{
			Addr:    ":8000",
			Handler: mux,
		}

		go func() {
			err := srv.ListenAndServe()
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				logger.Fatalf("error starting server: %s", err.Error())
			}
		}()

		logger.Printf("starting gosea %s", Version)

		<-ctx.Done()

		srv.Close()

		logger.Print("stopping service")
	*/
}
