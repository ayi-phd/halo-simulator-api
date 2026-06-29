package dispatcher

import (
	"context"
	"log"
	"math/rand"
	"time"

	"halo-simulator/internal/domain"
	"halo-simulator/internal/events"
	"halo-simulator/internal/store"
)

// workDuration returns a simulated task duration between 5 and 15 seconds.
func workDuration() time.Duration {
	return time.Duration(5+rand.Intn(11)) * time.Second
}

// Dispatcher reacts to AssignmentCreated events, simulates work execution
// in a goroutine per assignment, then frees resources and publishes results.
type Dispatcher struct {
	tasks       *store.TaskStore
	crews       *store.CrewStore
	equipment   *store.EquipmentStore
	assignments *store.AssignmentStore
	bus         *events.Bus
	sub         <-chan events.Event
}

func NewDispatcher(
	tasks *store.TaskStore,
	crews *store.CrewStore,
	equipment *store.EquipmentStore,
	assignments *store.AssignmentStore,
	bus *events.Bus,
) *Dispatcher {
	return &Dispatcher{
		tasks:       tasks,
		crews:       crews,
		equipment:   equipment,
		assignments: assignments,
		bus:         bus,
		sub:         bus.Subscribe(),
	}
}

// Run starts the Dispatcher's event loop. Call in its own goroutine.
func (d *Dispatcher) Run(ctx context.Context) {
	for {
		select {
		case e := <-d.sub:
			d.handle(ctx, e)
		case <-ctx.Done():
			return
		}
	}
}

func (d *Dispatcher) handle(ctx context.Context, e events.Event) {
	if e.Type != events.AssignmentCreated {
		return
	}
	assignment, ok := e.Payload.(domain.Assignment)
	if !ok {
		log.Printf("dispatcher: unexpected payload type for %s", e.Type)
		return
	}
	// Each assignment executes concurrently in its own goroutine.
	go d.execute(ctx, assignment)
}

// execute simulates the ground crew performing the work, then updates all
// related state and publishes the outcome events.
func (d *Dispatcher) execute(ctx context.Context, assignment domain.Assignment) {
	log.Printf("dispatcher: starting assignment %s (crew %s)", assignment.ID, assignment.CrewID)

	// Cancellable sleep — wakes on completion or server shutdown.
	select {
	case <-time.After(workDuration()):
	case <-ctx.Done():
		return
	}

	now := time.Now()

	// Complete the assignment.
	assignment.Status = domain.AssignmentStatusCompleted
	assignment.CompletedAt = &now
	d.assignments.Save(assignment)

	// Free the crew member and announce availability for pending tasks.
	if crew, ok := d.crews.Get(assignment.CrewID); ok {
		crew.Status = domain.CrewStatusAvailable
		d.crews.Save(crew)
		d.bus.Publish(events.Event{Type: events.CrewAvailable, Payload: crew})
	}

	// Free the equipment (if any was used).
	if assignment.EquipmentID != "" {
		if equip, ok := d.equipment.Get(assignment.EquipmentID); ok {
			equip.Status = domain.EquipmentStatusAvailable
			d.equipment.Save(equip)
			d.bus.Publish(events.Event{Type: events.EquipmentAvailable, Payload: equip})
		}
	}

	// Mark the task completed.
	if task, ok := d.tasks.Get(assignment.TaskID); ok {
		task.Status = domain.TaskStatusCompleted
		d.tasks.Save(task)
	}

	d.bus.Publish(events.Event{
		Type:    events.AssignmentCompleted,
		Payload: assignment,
	})

	log.Printf("dispatcher: assignment %s completed", assignment.ID)
}
