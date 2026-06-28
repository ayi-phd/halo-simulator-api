package simulator

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"halo-simulator/internal/domain"
	"halo-simulator/internal/events"
	"halo-simulator/internal/store"
)

// Simulator seeds the world with crew and equipment, then continuously
// publishes random operational events to drive the pipeline.
type Simulator struct {
	flights   *store.FlightStore
	crews     *store.CrewStore
	equipment *store.EquipmentStore
	bus       *events.Bus
	airlines  []domain.Airline
}

func NewSimulator(
	flights *store.FlightStore,
	crews *store.CrewStore,
	equipment *store.EquipmentStore,
	bus *events.Bus,
) *Simulator {
	return &Simulator{
		flights:   flights,
		crews:     crews,
		equipment: equipment,
		bus:       bus,
	}
}

// Run seeds initial data then starts two tickers:
//   - every 10s: a new flight arrives
//   - every 25s: a random operational event (equipment breakdown, flight delay)
func (s *Simulator) Run(ctx context.Context) {
	s.seed()

	flightTicker := time.NewTicker(10 * time.Second)
	eventTicker  := time.NewTicker(25 * time.Second)
	defer flightTicker.Stop()
	defer eventTicker.Stop()

	for {
		select {
		case <-flightTicker.C:
			s.arriveRandomFlight()
		case <-eventTicker.C:
			s.randomOperationalEvent(ctx)
		case <-ctx.Done():
			return
		}
	}
}

// seed registers airlines, crew, and equipment into their stores.
func (s *Simulator) seed() {
	for _, a := range []domain.Airline{
		{Name: "Delta", IATA: "DL"},
		{Name: "United", IATA: "UA"},
		{Name: "Southwest", IATA: "WN"},
	} {
		s.airlines = append(s.airlines, s.flights.SaveAirline(a))
	}

	crewSeeds := []struct {
		name string
		role domain.CrewRole
	}{
		{"Alice", domain.CrewRoleFueler},
		{"Bob", domain.CrewRoleFueler},
		{"Carol", domain.CrewRoleCleaner},
		{"Dave", domain.CrewRoleCleaner},
		{"Eve", domain.CrewRoleCaterer},
		{"Frank", domain.CrewRoleCaterer},
		{"Grace", domain.CrewRoleBaggageHandler},
		{"Henry", domain.CrewRoleBaggageHandler},
		{"Iris", domain.CrewRolePushback},
		{"Jack", domain.CrewRolePushback},
	}
	for _, cs := range crewSeeds {
		s.crews.Save(domain.Crew{
			Name:     cs.name,
			Role:     cs.role,
			Status:   domain.CrewStatusAvailable,
			Location: "Base",
		})
	}

	for _, t := range []domain.EquipmentType{
		domain.EquipmentTypeFuelTruck,
		domain.EquipmentTypeFuelTruck,
		domain.EquipmentTypePushbackTug,
		domain.EquipmentTypePushbackTug,
		domain.EquipmentTypeBeltLoader,
		domain.EquipmentTypeBeltLoader,
		domain.EquipmentTypeCateringTruck,
		domain.EquipmentTypeCateringTruck,
	} {
		s.equipment.Save(domain.GroundEquipment{
			Type:     t,
			Status:   domain.EquipmentStatusAvailable,
			Location: "Depot",
		})
	}

	log.Printf("simulator: seeded %d airlines, %d crew, 8 equipment", len(s.airlines), len(crewSeeds))
}

func (s *Simulator) arriveRandomFlight() {
	airline := s.airlines[rand.Intn(len(s.airlines))]
	flight := s.flights.SaveFlight(domain.Flight{
		FlightNumber:       fmt.Sprintf("%s%d", airline.IATA, 100+rand.Intn(900)),
		AirlineID:          airline.ID,
		Gate:               fmt.Sprintf("%c%d", 'A'+rune(rand.Intn(4)), 1+rand.Intn(10)),
		Status:             domain.FlightStatusArrived,
		ScheduledDeparture: time.Now().Add(2 * time.Hour),
		ActualArrival:      time.Now(),
	})
	s.bus.Publish(events.Event{Type: events.FlightArrived, Payload: flight})
	log.Printf("simulator: flight %s arrived at gate %s", flight.FlightNumber, flight.Gate)
}

// randomOperationalEvent picks one of: equipment breakdown or flight delay.
func (s *Simulator) randomOperationalEvent(ctx context.Context) {
	if rand.Intn(2) == 0 {
		s.breakRandomEquipment(ctx)
	} else {
		s.delayRandomFlight()
	}
}

func (s *Simulator) breakRandomEquipment(ctx context.Context) {
	available := s.equipment.Available()
	if len(available) == 0 {
		return
	}
	equip := available[rand.Intn(len(available))]
	equip.Status = domain.EquipmentStatusBroken
	s.equipment.Save(equip)
	s.bus.Publish(events.Event{Type: events.EquipmentBroken, Payload: equip})
	log.Printf("simulator: equipment %s (%s) broke down", equip.ID, equip.Type)

	// Auto-repair after 30 seconds using the same cancellable-sleep pattern
	// as the Dispatcher — the goroutine exits cleanly on shutdown.
	go func() {
		select {
		case <-time.After(30 * time.Second):
			equip.Status = domain.EquipmentStatusAvailable
			s.equipment.Save(equip)
			s.bus.Publish(events.Event{Type: events.EquipmentAvailable, Payload: equip})
			log.Printf("simulator: equipment %s (%s) repaired", equip.ID, equip.Type)
		case <-ctx.Done():
		}
	}()
}

func (s *Simulator) delayRandomFlight() {
	flights := s.flights.AllFlights()
	arrived := flights[:0]
	for _, f := range flights {
		if f.Status == domain.FlightStatusArrived {
			arrived = append(arrived, f)
		}
	}
	if len(arrived) == 0 {
		return
	}
	flight := arrived[rand.Intn(len(arrived))]
	flight.Status = domain.FlightStatusDelayed
	flight.ScheduledDeparture = flight.ScheduledDeparture.Add(30 * time.Minute)
	s.flights.SaveFlight(flight)
	s.bus.Publish(events.Event{Type: events.FlightDelayed, Payload: flight})
	log.Printf("simulator: flight %s delayed, new departure %s", flight.FlightNumber, flight.ScheduledDeparture.Format(time.Kitchen))
}
