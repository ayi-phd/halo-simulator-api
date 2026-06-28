package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"halo-simulator/internal/api"
	"halo-simulator/internal/dispatcher"
	"halo-simulator/internal/events"
	"halo-simulator/internal/metrics"
	"halo-simulator/internal/planner"
	"halo-simulator/internal/scheduler"
	"halo-simulator/internal/store"
	"halo-simulator/internal/ws"
)

func main() {
	ctx := context.Background()

	// Stores
	flights     := store.NewFlightStore()
	tasks       := store.NewTaskStore()
	crews       := store.NewCrewStore()
	equipment   := store.NewEquipmentStore()
	assignments := store.NewAssignmentStore()

	// Event bus
	bus := events.NewBus()

	// Components
	hub        := ws.NewHub(bus)
	collector  := metrics.NewCollector(flights, tasks, crews, assignments)
	plan       := planner.NewPlanner(tasks, bus)
	sched      := scheduler.NewScheduler(tasks, crews, equipment, assignments, bus)
	disp       := dispatcher.NewDispatcher(tasks, crews, equipment, assignments, bus)

	// Start background goroutines
	go bus.Run(ctx)
	go hub.Run(ctx)
	go plan.Run(ctx)
	go sched.Run(ctx)
	go disp.Run(ctx)

	// HTTP
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Mount("/", api.NewHandler(hub, collector).Routes())

	log.Println("HALO Simulator starting on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
