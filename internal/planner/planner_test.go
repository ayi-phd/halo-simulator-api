package planner_test

import (
	"context"
	"testing"
	"time"

	"halo-simulator/internal/domain"
	"halo-simulator/internal/events"
	"halo-simulator/internal/planner"
	"halo-simulator/internal/store"
)

func TestPlanner_FlightArrived_CreatesAllTaskTypes(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tasks := store.NewTaskStore()
	bus := events.NewBus()
	go bus.Run(ctx)

	sub := bus.Subscribe()

	go planner.NewPlanner(tasks, bus).Run(ctx)

	flight := domain.Flight{ID: "f1", FlightNumber: "DL100", Status: domain.FlightStatusArrived}
	bus.Publish(events.Event{Type: events.FlightArrived, Payload: flight})

	// Collect all TaskCreated events published by the planner.
	got := collectEvents(t, sub, events.TaskCreated, 5, 2*time.Second)

	wantTypes := map[domain.TaskType]bool{
		domain.TaskTypeFuel:     false,
		domain.TaskTypeCleaning: false,
		domain.TaskTypeCatering: false,
		domain.TaskTypeBaggage:  false,
		domain.TaskTypePushback: false,
	}
	for _, e := range got {
		task := e.Payload.(domain.Task)
		if task.FlightID != "f1" {
			t.Errorf("task %s: want flightID f1, got %s", task.ID, task.FlightID)
		}
		if task.Status != domain.TaskStatusPending {
			t.Errorf("task %s: want pending, got %s", task.ID, task.Status)
		}
		wantTypes[task.Type] = true
	}
	for tt, seen := range wantTypes {
		if !seen {
			t.Errorf("missing task type: %s", tt)
		}
	}

	if n := len(tasks.ByFlight("f1")); n != 5 {
		t.Errorf("want 5 tasks in store, got %d", n)
	}
}

func TestPlanner_OtherEvents_Ignored(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tasks := store.NewTaskStore()
	bus := events.NewBus()
	go bus.Run(ctx)
	go planner.NewPlanner(tasks, bus).Run(ctx)

	bus.Publish(events.Event{Type: events.EquipmentBroken})
	bus.Publish(events.Event{Type: events.CrewAvailable})

	time.Sleep(100 * time.Millisecond)
	if n := len(tasks.All()); n != 0 {
		t.Errorf("want 0 tasks for non-flight events, got %d", n)
	}
}

// collectEvents waits for n events of the given type, failing if timeout elapses.
func collectEvents(t *testing.T, sub <-chan events.Event, kind events.EventType, n int, timeout time.Duration) []events.Event {
	t.Helper()
	var result []events.Event
	deadline := time.After(timeout)
	for len(result) < n {
		select {
		case e := <-sub:
			if e.Type == kind {
				result = append(result, e)
			}
		case <-deadline:
			t.Fatalf("timeout: want %d %s events, got %d", n, kind, len(result))
		}
	}
	return result
}
