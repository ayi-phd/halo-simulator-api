package scheduler

import (
	"context"
	"log"
	"time"

	"halo-simulator/internal/domain"
	"halo-simulator/internal/events"
	"halo-simulator/internal/store"
)

// taskCrewRole maps each task type to the crew role required to perform it.
var taskCrewRole = map[domain.TaskType]domain.CrewRole{
	domain.TaskTypeFuel:     domain.CrewRoleFueler,
	domain.TaskTypeCleaning: domain.CrewRoleCleaner,
	domain.TaskTypeCatering: domain.CrewRoleCaterer,
	domain.TaskTypeBaggage:  domain.CrewRoleBaggageHandler,
	domain.TaskTypePushback: domain.CrewRolePushback,
}

// taskEquipmentType maps task types that require a specific piece of equipment.
// Cleaning is intentionally absent — it needs crew only.
var taskEquipmentType = map[domain.TaskType]domain.EquipmentType{
	domain.TaskTypeFuel:     domain.EquipmentTypeFuelTruck,
	domain.TaskTypeCatering: domain.EquipmentTypeCateringTruck,
	domain.TaskTypeBaggage:  domain.EquipmentTypeBeltLoader,
	domain.TaskTypePushback: domain.EquipmentTypePushbackTug,
}

// Scheduler subscribes to TaskCreated events, selects the best available
// crew and equipment, and publishes an AssignmentCreated event.
type Scheduler struct {
	tasks       *store.TaskStore
	crews       *store.CrewStore
	equipment   *store.EquipmentStore
	assignments *store.AssignmentStore
	bus         *events.Bus
}

func NewScheduler(
	tasks *store.TaskStore,
	crews *store.CrewStore,
	equipment *store.EquipmentStore,
	assignments *store.AssignmentStore,
	bus *events.Bus,
) *Scheduler {
	return &Scheduler{
		tasks:       tasks,
		crews:       crews,
		equipment:   equipment,
		assignments: assignments,
		bus:         bus,
	}
}

// Run starts the Scheduler's event loop. Call in its own goroutine.
func (s *Scheduler) Run(ctx context.Context) {
	sub := s.bus.Subscribe()
	for {
		select {
		case e := <-sub:
			s.handle(e)
		case <-ctx.Done():
			return
		}
	}
}

func (s *Scheduler) handle(e events.Event) {
	if e.Type != events.TaskCreated {
		return
	}
	task, ok := e.Payload.(domain.Task)
	if !ok {
		log.Printf("scheduler: unexpected payload type for %s", e.Type)
		return
	}
	s.schedule(task)
}

func (s *Scheduler) schedule(task domain.Task) {
	crew, ok := s.findCrew(task.Type)
	if !ok {
		log.Printf("scheduler: no available crew for task %s (%s) — will retry on next crew.available event", task.ID, task.Type)
		return
	}

	// Mark crew busy so subsequent tasks don't double-assign them.
	crew.Status = domain.CrewStatusBusy
	s.crews.Save(crew)

	// Equipment is optional; cleaning tasks need no equipment.
	var equipmentID string
	if equip, hasEquip := s.findEquipment(task.Type); hasEquip {
		equip.Status = domain.EquipmentStatusBusy
		s.equipment.Save(equip)
		equipmentID = equip.ID
	}

	// Advance task state.
	task.Status = domain.TaskStatusAssigned
	s.tasks.Save(task)

	assignment := s.assignments.Save(domain.Assignment{
		TaskID:      task.ID,
		CrewID:      crew.ID,
		EquipmentID: equipmentID,
		Status:      domain.AssignmentStatusActive,
		ETA:         time.Now().Add(5 * time.Minute),
		CreatedAt:   time.Now(),
	})

	s.bus.Publish(events.Event{
		Type:    events.AssignmentCreated,
		Payload: assignment,
	})

	log.Printf("scheduler: assigned task %s (%s) to crew %s", task.ID, task.Type, crew.Name)
}

// findCrew returns the first available crew member whose role matches the task.
func (s *Scheduler) findCrew(taskType domain.TaskType) (domain.Crew, bool) {
	role := taskCrewRole[taskType]
	for _, c := range s.crews.Available() {
		if c.Role == role {
			return c, true
		}
	}
	return domain.Crew{}, false
}

// findEquipment returns the first available equipment matching the task type,
// if equipment is required for that task at all.
func (s *Scheduler) findEquipment(taskType domain.TaskType) (domain.GroundEquipment, bool) {
	eType, required := taskEquipmentType[taskType]
	if !required {
		return domain.GroundEquipment{}, false
	}
	for _, e := range s.equipment.Available() {
		if e.Type == eType {
			return e, true
		}
	}
	return domain.GroundEquipment{}, false
}
