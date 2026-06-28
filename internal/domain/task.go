package domain

import "time"

type TaskType string

const (
	TaskTypeFuel     TaskType = "fuel"
	TaskTypeCleaning TaskType = "cleaning"
	TaskTypeCatering TaskType = "catering"
	TaskTypeBaggage  TaskType = "baggage"
	TaskTypePushback TaskType = "pushback"
)

type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusAssigned   TaskStatus = "assigned"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusCompleted  TaskStatus = "completed"
)

type Task struct {
	ID        string     `json:"id"`
	FlightID  string     `json:"flight_id"`
	Type      TaskType   `json:"type"`
	Status    TaskStatus `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
}
