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
		name          string
		setup         func(f *store.FlightStore, ts *store.TaskStore, cs *store.CrewStore, as *store.AssignmentStore)
		wantActive    int
		wantPending   int
		wantCompleted int
		wantBusy      int
	}{
		{
			name:  "empty stores return all zeros",
			setup: func(f *store.FlightStore, ts *store.TaskStore, cs *store.CrewStore, as *store.AssignmentStore) {},
		},
		{
			name: "counts reflect store state correctly",
			setup: func(f *store.FlightStore, ts *store.TaskStore, cs *store.CrewStore, as *store.AssignmentStore) {
				f.SaveFlight(domain.Flight{Status: domain.FlightStatusArrived})
				f.SaveFlight(domain.Flight{Status: domain.FlightStatusDeparted})
				ts.Save(domain.Task{Status: domain.TaskStatusPending})
				ts.Save(domain.Task{Status: domain.TaskStatusCompleted})
				ts.Save(domain.Task{Status: domain.TaskStatusCompleted})
				ts.Save(domain.Task{Status: domain.TaskStatusAssigned})
				cs.Save(domain.Crew{Status: domain.CrewStatusBusy})
				cs.Save(domain.Crew{Status: domain.CrewStatusBusy})
				cs.Save(domain.Crew{Status: domain.CrewStatusAvailable})
			},
			wantActive:    1,
			wantPending:   1,
			wantCompleted: 2,
			wantBusy:      2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flights     := store.NewFlightStore()
			tasks       := store.NewTaskStore()
			crews       := store.NewCrewStore()
			assignments := store.NewAssignmentStore()
			tt.setup(flights, tasks, crews, assignments)

			snap := metrics.NewCollector(flights, tasks, crews, assignments).Snapshot()

			if snap.ActiveFlights != tt.wantActive {
				t.Errorf("ActiveFlights: want %d, got %d", tt.wantActive, snap.ActiveFlights)
			}
			if snap.PendingTasks != tt.wantPending {
				t.Errorf("PendingTasks: want %d, got %d", tt.wantPending, snap.PendingTasks)
			}
			if snap.CompletedTasks != tt.wantCompleted {
				t.Errorf("CompletedTasks: want %d, got %d", tt.wantCompleted, snap.CompletedTasks)
			}
			if snap.BusyCrews != tt.wantBusy {
				t.Errorf("BusyCrews: want %d, got %d", tt.wantBusy, snap.BusyCrews)
			}
		})
	}
}

func TestCollector_AvgDispatchSeconds(t *testing.T) {
	flights     := store.NewFlightStore()
	tasks       := store.NewTaskStore()
	crews       := store.NewCrewStore()
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

	snap := metrics.NewCollector(flights, tasks, crews, assignments).Snapshot()
	if snap.AvgDispatchSeconds != 4.0 {
		t.Errorf("AvgDispatchSeconds: want 4.0, got %f", snap.AvgDispatchSeconds)
	}
}
