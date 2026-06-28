package scheduler_test

import (
	"context"
	"testing"
	"time"

	"halo-simulator/internal/domain"
	"halo-simulator/internal/events"
	"halo-simulator/internal/scheduler"
	"halo-simulator/internal/store"
)

func TestScheduler_TaskCreated_ProducesAssignment(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	taskStore  := store.NewTaskStore()
	crewStore  := store.NewCrewStore()
	equipStore := store.NewEquipmentStore()
	assignStore := store.NewAssignmentStore()
	bus := events.NewBus()
	go bus.Run(ctx)

	crew  := crewStore.Save(domain.Crew{Name: "Alice", Role: domain.CrewRoleFueler, Status: domain.CrewStatusAvailable})
	equip := equipStore.Save(domain.GroundEquipment{Type: domain.EquipmentTypeFuelTruck, Status: domain.EquipmentStatusAvailable})

	sub := bus.Subscribe()
	go scheduler.NewScheduler(taskStore, crewStore, equipStore, assignStore, bus).Run(ctx)

	task := taskStore.Save(domain.Task{FlightID: "f1", Type: domain.TaskTypeFuel, Status: domain.TaskStatusPending, CreatedAt: time.Now()})
	bus.Publish(events.Event{Type: events.TaskCreated, Payload: task})

	select {
	case e := <-sub:
		if e.Type != events.AssignmentCreated {
			// skip other events (e.g. echoes) and wait for the one we want
			select {
			case e = <-sub:
				if e.Type != events.AssignmentCreated {
					t.Fatalf("want AssignmentCreated, got %s", e.Type)
				}
			case <-time.After(2 * time.Second):
				t.Fatal("timed out waiting for AssignmentCreated")
			}
		}
		a := e.Payload.(domain.Assignment)
		if a.CrewID != crew.ID {
			t.Errorf("want crew %s, got %s", crew.ID, a.CrewID)
		}
		if a.EquipmentID != equip.ID {
			t.Errorf("want equipment %s, got %s", equip.ID, a.EquipmentID)
		}
		if a.Status != domain.AssignmentStatusActive {
			t.Errorf("want active, got %s", a.Status)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for AssignmentCreated")
	}

	// Crew and equipment should now be marked busy.
	c, _ := crewStore.Get(crew.ID)
	if c.Status != domain.CrewStatusBusy {
		t.Errorf("want crew busy, got %s", c.Status)
	}
	eq, _ := equipStore.Get(equip.ID)
	if eq.Status != domain.EquipmentStatusBusy {
		t.Errorf("want equipment busy, got %s", eq.Status)
	}
}

func TestScheduler_NoAvailableCrew_SkipsTask(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	taskStore   := store.NewTaskStore()
	crewStore   := store.NewCrewStore()
	equipStore  := store.NewEquipmentStore()
	assignStore := store.NewAssignmentStore()
	bus := events.NewBus()
	go bus.Run(ctx)

	// No crew seeded — scheduler should produce no assignment.
	sub := bus.Subscribe()
	go scheduler.NewScheduler(taskStore, crewStore, equipStore, assignStore, bus).Run(ctx)

	task := taskStore.Save(domain.Task{FlightID: "f1", Type: domain.TaskTypeFuel, Status: domain.TaskStatusPending, CreatedAt: time.Now()})
	bus.Publish(events.Event{Type: events.TaskCreated, Payload: task})

	select {
	case e := <-sub:
		if e.Type == events.AssignmentCreated {
			t.Error("expected no assignment when no crew available")
		}
	case <-time.After(300 * time.Millisecond):
		// Correct: no assignment produced.
	}

	if n := len(assignStore.All()); n != 0 {
		t.Errorf("want 0 assignments, got %d", n)
	}
}
