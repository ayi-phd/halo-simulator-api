package domain

type CrewRole string

const (
	CrewRoleFueler         CrewRole = "fueler"
	CrewRoleCleaner        CrewRole = "cleaner"
	CrewRoleCaterer        CrewRole = "caterer"
	CrewRoleBaggageHandler CrewRole = "baggage_handler"
	CrewRolePushback       CrewRole = "pushback"
)

type CrewStatus string

const (
	CrewStatusAvailable CrewStatus = "available"
	CrewStatusBusy      CrewStatus = "busy"
	CrewStatusOffDuty   CrewStatus = "off_duty"
)

type Crew struct {
	ID       string     `json:"id"`
	Name     string     `json:"name"`
	Role     CrewRole   `json:"role"`
	Status   CrewStatus `json:"status"`
	Location string     `json:"location"`
}
