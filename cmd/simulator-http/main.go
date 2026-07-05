package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"halo-simulator/internal/simclient"
)

var airlinePrefixes = []string{"DL", "UA", "WN"}

var crewSeeds = []struct{ name, role string }{
	{"Alice", "fueler"},
	{"Bob", "fueler"},
	{"Carol", "cleaner"},
	{"Dave", "cleaner"},
	{"Eve", "caterer"},
	{"Frank", "caterer"},
	{"Grace", "baggage_handler"},
	{"Henry", "baggage_handler"},
	{"Iris", "pushback"},
	{"Jack", "pushback"},
}

var equipmentSeeds = []string{
	"fuel_truck", "fuel_truck",
	"pushback_tug", "pushback_tug",
	"belt_loader", "belt_loader",
	"catering_truck", "catering_truck",
}

func main() {
	serverURL := flag.String("server", "http://localhost:8080", "HALO server base URL")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	client := simclient.New(*serverURL)
	waitForServer(ctx, client, *serverURL)

	equipmentIDs := seed(client)
	var flightIDs []string

	flightTicker := time.NewTicker(10 * time.Second)
	eventTicker  := time.NewTicker(25 * time.Second)
	burstTicker  := time.NewTicker(45 * time.Second)
	defer flightTicker.Stop()
	defer eventTicker.Stop()
	defer burstTicker.Stop()

	log.Printf("simulator-http: running — flights every 10s, events every 25s, burst every 45s")

	for {
		select {
		case <-flightTicker.C:
			if id := arriveRandomFlight(client); id != "" {
				flightIDs = append(flightIDs, id)
			}
		case <-eventTicker.C:
			randomOperationalEvent(ctx, client, equipmentIDs, flightIDs)
		case <-burstTicker.C:
			ids := arriveFlightBurst(client, 3)
			flightIDs = append(flightIDs, ids...)
		case <-ctx.Done():
			log.Println("simulator-http: stopped")
			return
		}
	}
}

// waitForServer polls /health until the server responds, then returns.
func waitForServer(ctx context.Context, client *simclient.Client, serverURL string) {
	for {
		if err := client.HealthCheck(); err == nil {
			log.Printf("simulator-http: server ready at %s", serverURL)
			return
		}
		log.Printf("simulator-http: waiting for server at %s...", serverURL)
		select {
		case <-time.After(2 * time.Second):
		case <-ctx.Done():
			return
		}
	}
}

// seed registers crew and equipment with the server. Returns the equipment IDs
// so the simulator can reference them in breakdown events.
func seed(client *simclient.Client) []string {
	for _, cs := range crewSeeds {
		if _, err := client.CreateCrew(cs.name, cs.role); err != nil {
			log.Printf("simulator-http: seed crew error: %v", err)
		}
	}

	var equipmentIDs []string
	for _, t := range equipmentSeeds {
		id, err := client.CreateEquipment(t)
		if err != nil {
			log.Printf("simulator-http: seed equipment error: %v", err)
			continue
		}
		equipmentIDs = append(equipmentIDs, id)
	}

	log.Printf("simulator-http: seeded %d crew, %d equipment", len(crewSeeds), len(equipmentIDs))
	return equipmentIDs
}

func arriveRandomFlight(client *simclient.Client) string {
	prefix := airlinePrefixes[rand.Intn(len(airlinePrefixes))]
	number := fmt.Sprintf("%s%d", prefix, 100+rand.Intn(900))
	gate   := fmt.Sprintf("%c%d", 'A'+rune(rand.Intn(4)), 1+rand.Intn(10))

	id, err := client.CreateFlight(number, gate)
	if err != nil {
		log.Printf("simulator-http: create flight error: %v", err)
		return ""
	}
	log.Printf("simulator-http: flight %s arrived at gate %s", number, gate)
	return id
}

// arriveFlightBurst arrives n flights simultaneously, saturating crew and
// equipment so that some tasks are forced into Pending state.
func arriveFlightBurst(client *simclient.Client, n int) []string {
	log.Printf("simulator-http: burst — arriving %d flights simultaneously", n)
	var ids []string
	for range n {
		if id := arriveRandomFlight(client); id != "" {
			ids = append(ids, id)
		}
	}
	return ids
}

func randomOperationalEvent(ctx context.Context, client *simclient.Client, equipmentIDs, flightIDs []string) {
	if rand.Intn(2) == 0 {
		breakRandomEquipment(ctx, client, equipmentIDs)
	} else {
		delayRandomFlight(client, flightIDs)
	}
}

func breakRandomEquipment(ctx context.Context, client *simclient.Client, ids []string) {
	if len(ids) == 0 {
		return
	}
	id := ids[rand.Intn(len(ids))]
	if err := client.BreakEquipment(id); err != nil {
		log.Printf("simulator-http: break equipment error: %v", err)
		return
	}
	log.Printf("simulator-http: equipment %s broken — repair in 30s", id)

	go func() {
		select {
		case <-time.After(30 * time.Second):
			if err := client.RepairEquipment(id); err != nil {
				log.Printf("simulator-http: repair equipment error: %v", err)
				return
			}
			log.Printf("simulator-http: equipment %s repaired", id)
		case <-ctx.Done():
		}
	}()
}

func delayRandomFlight(client *simclient.Client, ids []string) {
	if len(ids) == 0 {
		log.Println("simulator-http: no flights yet to delay")
		return
	}
	id := ids[rand.Intn(len(ids))]
	if err := client.DelayFlight(id); err != nil {
		log.Printf("simulator-http: delay flight error: %v", err)
		return
	}
	log.Printf("simulator-http: flight %s delayed", id)
}
