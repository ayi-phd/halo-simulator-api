package metrics

import (
	"time"

	"halo-simulator/internal/domain"
	"halo-simulator/internal/store"
)

// Snapshot is a point-in-time read of operational KPIs.
type Snapshot struct {
	TotalFlights       int     `json:"total_flights"`
	TotalTasks         int     `json:"total_tasks"`
	PendingTasks       int     `json:"pending_tasks"`
	CompletedTasks     int     `json:"completed_tasks"`
	BusyCrew           int     `json:"busy_crew"`
	BusyEquipment      int     `json:"busy_equipment"`
	AvgDispatchSeconds float64 `json:"avg_dispatch_seconds"`
}

// Collector computes metrics by reading directly from the stores.
// No redundant counters to keep in sync — the stores are the source of truth.
type Collector struct {
	flights     *store.FlightStore
	tasks       *store.TaskStore
	crews       *store.CrewStore
	equipment   *store.EquipmentStore
	assignments *store.AssignmentStore
}

func NewCollector(
	flights *store.FlightStore,
	tasks *store.TaskStore,
	crews *store.CrewStore,
	equipment *store.EquipmentStore,
	assignments *store.AssignmentStore,
) *Collector {
	return &Collector{
		flights:     flights,
		tasks:       tasks,
		crews:       crews,
		equipment:   equipment,
		assignments: assignments,
	}
}

// Snapshot returns current KPIs computed from live store state.
func (c *Collector) Snapshot() Snapshot {
	var pending, completed, total int
	for _, t := range c.tasks.All() {
		total++
		switch t.Status {
		case domain.TaskStatusPending:
			pending++
		case domain.TaskStatusCompleted:
			completed++
		}
	}

	var busyCrew int
	for _, crew := range c.crews.All() {
		if crew.Status == domain.CrewStatusBusy {
			busyCrew++
		}
	}

	var busyEquip int
	for _, e := range c.equipment.All() {
		if e.Status == domain.EquipmentStatusBusy {
			busyEquip++
		}
	}

	return Snapshot{
		TotalFlights:       len(c.flights.AllFlights()),
		TotalTasks:         total,
		PendingTasks:       pending,
		CompletedTasks:     completed,
		BusyCrew:           busyCrew,
		BusyEquipment:      busyEquip,
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
