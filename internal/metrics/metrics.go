package metrics

import (
	"time"

	"halo-simulator/internal/domain"
	"halo-simulator/internal/store"
)

// Snapshot is a point-in-time read of operational KPIs.
type Snapshot struct {
	ActiveFlights      int     `json:"active_flights"`
	PendingTasks       int     `json:"pending_tasks"`
	BusyCrews          int     `json:"busy_crews"`
	CompletedTasks     int     `json:"completed_tasks"`
	AvgDispatchSeconds float64 `json:"avg_dispatch_seconds"`
}

// Collector computes metrics by reading directly from the stores.
// No redundant counters to keep in sync — the stores are the source of truth.
type Collector struct {
	flights     *store.FlightStore
	tasks       *store.TaskStore
	crews       *store.CrewStore
	assignments *store.AssignmentStore
}

func NewCollector(
	flights *store.FlightStore,
	tasks *store.TaskStore,
	crews *store.CrewStore,
	assignments *store.AssignmentStore,
) *Collector {
	return &Collector{
		flights:     flights,
		tasks:       tasks,
		crews:       crews,
		assignments: assignments,
	}
}

// Snapshot returns current KPIs computed from live store state.
func (c *Collector) Snapshot() Snapshot {
	var pending, completed int
	for _, t := range c.tasks.All() {
		switch t.Status {
		case domain.TaskStatusPending:
			pending++
		case domain.TaskStatusCompleted:
			completed++
		}
	}

	var busy int
	for _, crew := range c.crews.All() {
		if crew.Status == domain.CrewStatusBusy {
			busy++
		}
	}

	var active int
	for _, f := range c.flights.AllFlights() {
		if f.Status == domain.FlightStatusArrived {
			active++
		}
	}

	return Snapshot{
		ActiveFlights:      active,
		PendingTasks:       pending,
		BusyCrews:          busy,
		CompletedTasks:     completed,
		AvgDispatchSeconds: c.avgDispatchSeconds(),
	}
}

// avgDispatchSeconds computes the mean time between task creation and
// assignment creation across all assignments seen so far.
func (c *Collector) avgDispatchSeconds() float64 {
	taskCreatedAt := make(map[string]time.Time)
	for _, t := range c.tasks.All() {
		taskCreatedAt[t.ID] = t.CreatedAt
	}

	var total float64
	var count int
	for _, a := range c.assignments.All() {
		if created, ok := taskCreatedAt[a.TaskID]; ok {
			total += a.CreatedAt.Sub(created).Seconds()
			count++
		}
	}
	if count == 0 {
		return 0
	}
	return total / float64(count)
}
