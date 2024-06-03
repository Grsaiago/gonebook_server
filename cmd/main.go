package main

import (
	"context"
	"github.com/Grsaiago/gonebook_server/internal/application"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
)

func main() {

	// setup custom logger as default slogger (structured logging)
	jsonLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(jsonLogger)

	// create app 'controller'
	app, err := application.New()
	if err != nil {
		slog.Error("db: Error on connect to database", "message", err.Error())
		os.Exit(1)
	}
	defer app.DbConn.Close(app.DbCtx) // isso deveria ficar aqui, invés de no New, não?

	//sync channel for graceful shutdown
	shutdownChannel := make(chan struct{})

	// create and setup router
	mux := http.NewServeMux()
	mux.HandleFunc("GET /contact", app.ContactService.GetAllContacts)
	mux.HandleFunc("POST /contact", app.ContactService.CreateContact)

	// Create the server
	server := http.Server{
		Addr:    ":8034",
		Handler: mux,
	}

	// setup the sighandler in another goroutine
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		slog.Info("http: Starting graceful shutdown")
		// shutdown signal received, kill server
		if err := server.Shutdown(context.Background()); err != nil {
			slog.Error("http: Shutdown error: %v\n", err)
		}
		close(shutdownChannel)
	}()

	// start server
	slog.Info("http: Starting server")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

	// wait for graceful shutdown
	<-shutdownChannel
	return
}
