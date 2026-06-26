package domain

import "time"

type FlightStatus string

const (
	FlightStatusScheduled FlightStatus = "scheduled"
	FlightStatusArrived   FlightStatus = "arrived"
	FlightStatusDelayed   FlightStatus = "delayed"
	FlightStatusDeparted  FlightStatus = "departed"
)

type Airline struct {
	ID   string
	Name string
	IATA string
}

type Flight struct {
	ID                 string
	FlightNumber       string
	AirlineID          string
	Gate               string
	Status             FlightStatus
	ScheduledDeparture time.Time
	ActualArrival      time.Time
}
