package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"halo-simulator/internal/api"
	"halo-simulator/internal/events"
	"halo-simulator/internal/ws"
)

func main() {
	ctx := context.Background()

	bus := events.NewBus()
	hub := ws.NewHub(bus)

	go bus.Run(ctx)
	go hub.Run(ctx)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	h := api.NewHandler(hub)
	r.Mount("/", h.Routes())

	log.Println("HALO Simulator starting on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
