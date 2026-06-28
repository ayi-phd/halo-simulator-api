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
	ID   string `json:"id"`
	Name string `json:"name"`
	IATA string `json:"iata"`
}

type Flight struct {
	ID                 string       `json:"id"`
	FlightNumber       string       `json:"flight_number"`
	AirlineID          string       `json:"airline_id"`
	Gate               string       `json:"gate"`
	Status             FlightStatus `json:"status"`
	ScheduledDeparture time.Time    `json:"scheduled_departure"`
	ActualArrival      time.Time    `json:"actual_arrival"`
}
