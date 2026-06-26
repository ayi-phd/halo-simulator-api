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
	ID        string
	FlightID  string
	Type      TaskType
	Status    TaskStatus
	CreatedAt time.Time
}
