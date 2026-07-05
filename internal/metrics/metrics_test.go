package metrics_test

import (
	"testing"
	"time"

	"halo-simulator/internal/domain"
	"halo-simulator/internal/metrics"
	"halo-simulator/internal/store"
)

func TestCollector_Snapshot(t *testing.T) {
	tests := []struct {
		name              string
		setup             func(f *store.FlightStore, ts *store.TaskStore, cs *store.CrewStore, es *store.EquipmentStore, as *store.AssignmentStore)
		wantFlights       int
		wantTotalTasks    int
		wantPending       int
		wantCompleted     int
		wantBusyCrew      int
		wantBusyEquipment int
	}{
		{
			name:  "empty stores return all zeros",
			setup: func(f *store.FlightStore, ts *store.TaskStore, cs *store.CrewStore, es *store.EquipmentStore, as *store.AssignmentStore) {},
		},
		{
			name: "counts reflect store state correctly",
			setup: func(f *store.FlightStore, ts *store.TaskStore, cs *store.CrewStore, es *store.EquipmentStore, as *store.AssignmentStore) {
				f.SaveFlight(domain.Flight{Status: domain.FlightStatusArrived})
				f.SaveFlight(domain.Flight{Status: domain.FlightStatusDeparted})
				ts.Save(domain.Task{Status: domain.TaskStatusPending})
				ts.Save(domain.Task{Status: domain.TaskStatusCompleted})
				ts.Save(domain.Task{Status: domain.TaskStatusCompleted})
				ts.Save(domain.Task{Status: domain.TaskStatusAssigned})
				cs.Save(domain.Crew{Status: domain.CrewStatusBusy})
				cs.Save(domain.Crew{Status: domain.CrewStatusBusy})
				cs.Save(domain.Crew{Status: domain.CrewStatusAvailable})
				es.Save(domain.GroundEquipment{Status: domain.EquipmentStatusBusy})
				es.Save(domain.GroundEquipment{Status: domain.EquipmentStatusAvailable})
			},
			wantFlights:       2,
			wantTotalTasks:    4,
			wantPending:       1,
			wantCompleted:     2,
			wantBusyCrew:      2,
			wantBusyEquipment: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flights     := store.NewFlightStore()
			tasks       := store.NewTaskStore()
			crews       := store.NewCrewStore()
			equipment   := store.NewEquipmentStore()
			assignments := store.NewAssignmentStore()
			tt.setup(flights, tasks, crews, equipment, assignments)

			snap := metrics.NewCollector(flights, tasks, crews, equipment, assignments).Snapshot()

			if snap.TotalFlights != tt.wantFlights {
				t.Errorf("TotalFlights: want %d, got %d", tt.wantFlights, snap.TotalFlights)
			}
			if snap.TotalTasks != tt.wantTotalTasks {
				t.Errorf("TotalTasks: want %d, got %d", tt.wantTotalTasks, snap.TotalTasks)
			}
			if snap.PendingTasks != tt.wantPending {
				t.Errorf("PendingTasks: want %d, got %d", tt.wantPending, snap.PendingTasks)
			}
			if snap.CompletedTasks != tt.wantCompleted {
				t.Errorf("CompletedTasks: want %d, got %d", tt.wantCompleted, snap.CompletedTasks)
			}
			if snap.BusyCrew != tt.wantBusyCrew {
				t.Errorf("BusyCrew: want %d, got %d", tt.wantBusyCrew, snap.BusyCrew)
			}
			if snap.BusyEquipment != tt.wantBusyEquipment {
				t.Errorf("BusyEquipment: want %d, got %d", tt.wantBusyEquipment, snap.BusyEquipment)
			}
		})
	}
}

func TestCollector_AvgDispatchSeconds(t *testing.T) {
	flights     := store.NewFlightStore()
	tasks       := store.NewTaskStore()
	crews       := store.NewCrewStore()
	equipment   := store.NewEquipmentStore()
	assignments := store.NewAssignmentStore()

	base := time.Now()
	task := tasks.Save(domain.Task{
		FlightID:  "f1",
		Type:      domain.TaskTypeFuel,
		Status:    domain.TaskStatusCompleted,
		CreatedAt: base,
	})
	assignments.Save(domain.Assignment{
		TaskID:    task.ID,
		Status:    domain.AssignmentStatusCompleted,
		CreatedAt: base.Add(4 * time.Second),
	})

	snap := metrics.NewCollector(flights, tasks, crews, equipment, assignments).Snapshot()
	if snap.AvgDispatchSeconds != 4.0 {
		t.Errorf("AvgDispatchSeconds: want 4.0, got %f", snap.AvgDispatchSeconds)
	}
}
