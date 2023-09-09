package server

import (
	"net/http"
	"time"

	"github.com/flukis/pagination-comparation/utils/config"
	"github.com/go-chi/chi/v5"
)

func Run(c *config.Config, r *chi.Mux) error {
	server := &http.Server{
		Handler:      r,
		Addr:         c.Listen.Address(),
		ReadTimeout:  time.Second * time.Duration(c.Listen.ReadTO),
		WriteTimeout: time.Second * time.Duration(c.Listen.WriteTO),
		IdleTimeout:  time.Second * time.Duration(c.Listen.IdleTO),
	}

	// Start server.
	return server.ListenAndServe()
}
