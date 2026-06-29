package domain

import "time"

type AssignmentStatus string

const (
	AssignmentStatusActive    AssignmentStatus = "active"
	AssignmentStatusCompleted AssignmentStatus = "completed"
	AssignmentStatusCancelled AssignmentStatus = "cancelled"
)

type Assignment struct {
	ID          string           `json:"id"`
	TaskID      string           `json:"task_id"`
	CrewID      string           `json:"crew_id"`
	EquipmentID string           `json:"equipment_id,omitempty"`
	Status      AssignmentStatus `json:"status"`
	ETA         time.Time        `json:"eta"`
	CreatedAt   time.Time        `json:"created_at"`
	CompletedAt *time.Time       `json:"completed_at,omitempty"`
}
