package main

import (
	"net/http"
	"time"
)

func (app *Application) Serve() error {
	srv := http.Server{
		Addr:         ":4040",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 30 * time.Second,
		Handler:      app.Routes(),
		ErrorLog:     app.errorLog,
	}

	return srv.ListenAndServe()
}
