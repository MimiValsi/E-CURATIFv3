package main

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {

	tlsConfig := &tls.Config{
		MinVersion:       tls.VersionTLS13,
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// Default parameters values to routes
	// See routers.go
	srv := &http.Server{
		Addr:         addr,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		TLSConfig:    tlsConfig,
	}

	shutdownError := make(chan error)

	go func() {
		// quit channel transporte la valeur de os.Signal = 1
		quit := make(chan os.Signal, 1)

		// signal.Notify reste à l'écoute des signauc SIGINT et SIGTERM
		// celui-ci envoi ces signaux vers la variable quit
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// q reçoit le signal de quit
		q := <-quit

		// Transforme s en string et rajoute dans le message
		app.logger.Info("Shuting down server", "signal", q.String())

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Shutdown is called while parsed with ctx
		// If everything is shutdown correctly, it will return nil
		// If not, for whatever reason, after 30 seconds it will be forced
		shutdownError <- srv.Shutdown(ctx)
	}()

	app.logger.Info("started server", "addr", srv.Addr)

	// By calling shutdown(), ListenAndServerTLS will sends us ErrServerClosed
	// If it's ErrServerClosed then everything went smoothly
	err := srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	// if not then something went bad
	err = <-shutdownError
	if err != nil {
		return nil
	}

	app.logger.Info("stopped server", "addr", srv.Addr)

	return nil
}
