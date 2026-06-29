package planner

import (
	"context"
	"log"
	"time"

	"halo-simulator/internal/domain"
	"halo-simulator/internal/events"
	"halo-simulator/internal/store"
)

// Planner subscribes to the event bus and determines what work needs
// to happen. It does not decide who performs the work — that is the
// Scheduler's responsibility.
type Planner struct {
	tasks *store.TaskStore
	bus   *events.Bus
	sub   <-chan events.Event
}

func NewPlanner(tasks *store.TaskStore, bus *events.Bus) *Planner {
	return &Planner{tasks: tasks, bus: bus, sub: bus.Subscribe()}
}

// Run starts the Planner's event loop. Call in its own goroutine.
// Exits cleanly when ctx is cancelled.
func (p *Planner) Run(ctx context.Context) {
	for {
		select {
		case e := <-p.sub:
			p.handle(e)
		case <-ctx.Done():
			return
		}
	}
}

func (p *Planner) handle(e events.Event) {
	switch e.Type {
	case events.FlightArrived:
		flight, ok := e.Payload.(domain.Flight)
		if !ok {
			log.Printf("planner: unexpected payload type for %s", e.Type)
			return
		}
		p.createTasks(flight)
	}
}

// createTasks saves one Task per required ground operation for the
// arriving flight and publishes a TaskCreated event for each.
func (p *Planner) createTasks(flight domain.Flight) {
	required := []domain.TaskType{
		domain.TaskTypeFuel,
		domain.TaskTypeCleaning,
		domain.TaskTypeCatering,
		domain.TaskTypeBaggage,
		domain.TaskTypePushback,
	}

	for _, tt := range required {
		task := p.tasks.Save(domain.Task{
			FlightID:  flight.ID,
			Type:      tt,
			Status:    domain.TaskStatusPending,
			CreatedAt: time.Now(),
		})

		p.bus.Publish(events.Event{
			Type:    events.TaskCreated,
			Payload: task,
		})

		log.Printf("planner: created task %s (%s) for flight %s", task.ID, task.Type, flight.FlightNumber)
	}
}
