// Package simclient is an HTTP client used by the simulator-api binary to
// drive the HALO server without sharing in-process memory.
package simclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Client sends requests to the HALO server REST API.
type Client struct {
	base string
	http *http.Client
}

func New(serverURL string) *Client {
	return &Client{
		base: serverURL,
		http: &http.Client{Timeout: 5 * time.Second},
	}
}

// HealthCheck returns nil when the server is up and responding.
func (c *Client) HealthCheck() error {
	resp, err := c.http.Get(c.base + "/health")
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check: status %d", resp.StatusCode)
	}
	return nil
}

// CreateFlight POSTs to /flights and returns the new flight's ID.
func (c *Client) CreateFlight(flightNumber, gate string) (string, error) {
	var out struct {
		ID string `json:"id"`
	}
	err := c.post("/flights", map[string]any{
		"flight_number": flightNumber,
		"gate":          gate,
	}, &out)
	return out.ID, err
}

// CreateCrew POSTs to /crew and returns the new crew member's ID.
func (c *Client) CreateCrew(name, role string) (string, error) {
	var out struct {
		ID string `json:"id"`
	}
	err := c.post("/crew", map[string]any{
		"name": name,
		"role": role,
	}, &out)
	return out.ID, err
}

// CreateEquipment POSTs to /equipment and returns the new equipment's ID.
func (c *Client) CreateEquipment(equipType string) (string, error) {
	var out struct {
		ID string `json:"id"`
	}
	err := c.post("/equipment", map[string]any{
		"type": equipType,
	}, &out)
	return out.ID, err
}

// DelayFlight publishes a flight.delayed event for the given flight ID.
func (c *Client) DelayFlight(flightID string) error {
	return c.post("/events", map[string]any{
		"type":      "flight.delayed",
		"flight_id": flightID,
	}, nil)
}

// BreakEquipment publishes an equipment.broken event.
func (c *Client) BreakEquipment(equipmentID string) error {
	return c.post("/events", map[string]any{
		"type":         "equipment.broken",
		"equipment_id": equipmentID,
	}, nil)
}

// RepairEquipment publishes an equipment.available event (auto-repair).
func (c *Client) RepairEquipment(equipmentID string) error {
	return c.post("/events", map[string]any{
		"type":         "equipment.available",
		"equipment_id": equipmentID,
	}, nil)
}

func (c *Client) post(path string, body any, out any) error {
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}
	resp, err := c.http.Post(c.base+path, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("POST %s: %w", path, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("POST %s: server returned %d", path, resp.StatusCode)
	}
	if out != nil {
		return json.NewDecoder(resp.Body).Decode(out)
	}
	return nil
}
