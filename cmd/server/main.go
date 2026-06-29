package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"halo-simulator/internal/api"
	"halo-simulator/internal/dispatcher"
	"halo-simulator/internal/domain"
	"halo-simulator/internal/events"
	"halo-simulator/internal/metrics"
	"halo-simulator/internal/planner"
	"halo-simulator/internal/scheduler"
	"halo-simulator/internal/store"
	"halo-simulator/internal/ws"
)

func main() {
	// ctx is cancelled when SIGINT (Ctrl+C) or SIGTERM is received.
	// All goroutines block on ctx.Done() and exit cleanly on cancellation.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Stores
	flights     := store.NewFlightStore()
	tasks       := store.NewTaskStore()
	crews       := store.NewCrewStore()
	equipment   := store.NewEquipmentStore()
	assignments := store.NewAssignmentStore()

	// Event bus
	bus := events.NewBus()

	// Seed static airline reference data.
	for _, a := range []domain.Airline{
		{Name: "Delta", IATA: "DL"},
		{Name: "United", IATA: "UA"},
		{Name: "Southwest", IATA: "WN"},
	} {
		flights.SaveAirline(a)
	}

	// Components
	hub       := ws.NewHub(bus)
	collector := metrics.NewCollector(flights, tasks, crews, assignments)
	plan      := planner.NewPlanner(tasks, bus)
	sched     := scheduler.NewScheduler(tasks, crews, equipment, assignments, bus)
	disp      := dispatcher.NewDispatcher(tasks, crews, equipment, assignments, bus)

	// Start background goroutines — all exit when ctx is cancelled.
	go bus.Run(ctx)
	go hub.Run(ctx)
	go plan.Run(ctx)
	go sched.Run(ctx)
	go disp.Run(ctx)

	// HTTP server — runs in its own goroutine so main can wait on ctx.
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Mount("/", api.NewHandler(api.Config{
		Hub:         hub,
		Metrics:     collector,
		Flights:     flights,
		Tasks:       tasks,
		Crews:       crews,
		Equipment:   equipment,
		Assignments: assignments,
		Bus:         bus,
	}).Routes())

	srv := &http.Server{Addr: ":8080", Handler: r}

	go func() {
		log.Println("HALO server starting on :8080 — run simulator-api to generate traffic")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Block until a signal arrives.
	<-ctx.Done()
	log.Println("shutdown signal received — draining...")

	// Give in-flight HTTP requests up to 5 seconds to complete.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}

	log.Println("HALO Simulator stopped")
}
