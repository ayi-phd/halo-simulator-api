package domain

import "time"

type AssignmentStatus string

const (
	AssignmentStatusActive    AssignmentStatus = "active"
	AssignmentStatusCompleted AssignmentStatus = "completed"
	AssignmentStatusCancelled AssignmentStatus = "cancelled"
)

type Assignment struct {
	ID          string
	TaskID      string
	CrewID      string
	EquipmentID string
	Status      AssignmentStatus
	ETA         time.Time
	CreatedAt   time.Time
	CompletedAt *time.Time
}
