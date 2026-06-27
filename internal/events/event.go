package events

// EventType identifies what happened. Using dot-notation keeps them
// grouped and readable in logs ("flight.arrived", "crew.available").
type EventType string

const (
	FlightArrived       EventType = "flight.arrived"
	FlightDelayed       EventType = "flight.delayed"
	CrewAvailable       EventType = "crew.available"
	CrewUnavailable     EventType = "crew.unavailable"
	EquipmentAvailable  EventType = "equipment.available"
	EquipmentBroken     EventType = "equipment.broken"
	TaskCompleted       EventType = "task.completed"
	AssignmentCreated   EventType = "assignment.created"
	AssignmentCompleted EventType = "assignment.completed"
)

// Event carries a type and an arbitrary payload.
// Subscribers type-assert the payload to the concrete type they expect,
// e.g.: flight := e.Payload.(domain.Flight)
type Event struct {
	Type    EventType
	Payload any
}
