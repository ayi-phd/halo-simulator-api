package store_test

import (
	"testing"

	"halo-simulator/internal/domain"
	"halo-simulator/internal/store"
)

func TestCrewStore_SaveAndGet(t *testing.T) {
	s := store.NewCrewStore()
	saved := s.Save(domain.Crew{Name: "Alice", Role: domain.CrewRoleFueler, Status: domain.CrewStatusAvailable})

	if saved.ID == "" {
		t.Fatal("expected ID to be generated")
	}
	got, ok := s.Get(saved.ID)
	if !ok {
		t.Fatal("expected crew to be found")
	}
	if got.Name != "Alice" {
		t.Errorf("want Alice, got %s", got.Name)
	}
}

func TestCrewStore_Available_FiltersCorrectly(t *testing.T) {
	s := store.NewCrewStore()
	s.Save(domain.Crew{Status: domain.CrewStatusAvailable})
	s.Save(domain.Crew{Status: domain.CrewStatusBusy})
	s.Save(domain.Crew{Status: domain.CrewStatusAvailable})
	s.Save(domain.Crew{Status: domain.CrewStatusOffDuty})

	got := s.Available()
	if len(got) != 2 {
		t.Errorf("want 2 available crew, got %d", len(got))
	}
}

func TestEquipmentStore_Available_FiltersCorrectly(t *testing.T) {
	s := store.NewEquipmentStore()
	s.Save(domain.GroundEquipment{Type: domain.EquipmentTypeFuelTruck, Status: domain.EquipmentStatusAvailable})
	s.Save(domain.GroundEquipment{Type: domain.EquipmentTypeFuelTruck, Status: domain.EquipmentStatusBusy})
	s.Save(domain.GroundEquipment{Type: domain.EquipmentTypeBeltLoader, Status: domain.EquipmentStatusBroken})

	got := s.Available()
	if len(got) != 1 {
		t.Errorf("want 1 available equipment, got %d", len(got))
	}
}

func TestTaskStore_ByFlight(t *testing.T) {
	s := store.NewTaskStore()
	s.Save(domain.Task{FlightID: "f1", Type: domain.TaskTypeFuel, Status: domain.TaskStatusPending})
	s.Save(domain.Task{FlightID: "f1", Type: domain.TaskTypeCleaning, Status: domain.TaskStatusPending})
	s.Save(domain.Task{FlightID: "f2", Type: domain.TaskTypeFuel, Status: domain.TaskStatusPending})

	got := s.ByFlight("f1")
	if len(got) != 2 {
		t.Errorf("want 2 tasks for f1, got %d", len(got))
	}
}

func TestFlightStore_SaveUpsert(t *testing.T) {
	s := store.NewFlightStore()
	f := s.SaveFlight(domain.Flight{FlightNumber: "DL100", Status: domain.FlightStatusArrived})

	f.Status = domain.FlightStatusDelayed
	s.SaveFlight(f)

	got, ok := s.GetFlight(f.ID)
	if !ok {
		t.Fatal("expected flight to be found")
	}
	if got.Status != domain.FlightStatusDelayed {
		t.Errorf("want delayed, got %s", got.Status)
	}
}
