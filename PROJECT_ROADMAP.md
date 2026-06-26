# PROJECT_ROADMAP.md

# Project

**HALO Simulator** — A HALO-inspired airport turnaround decision engine.

The application simulates airport ground operations by reacting to real-time events and assigning crews and equipment to aircraft turnaround tasks.

---

# Objectives

- Learn Go by building a realistic backend.
- Demonstrate event-driven architecture.
- Demonstrate Go concurrency.
- Build something directly relevant to Moonware.

---

# Tech Stack

| Layer | Technology |
|--------|------------|
| Router | chi |
| WebSocket | gorilla/websocket |
| UUID | google/uuid |
| Testing | Go testing |
| Storage | map + sync.RWMutex |

---

# Package Structure

```
cmd/
  server/main.go

internal/
  domain/
  store/
  events/
  planner/
  scheduler/
  dispatcher/
  api/
  ws/
  metrics/
```

---

# Domain Model

Airport
- Airline
- Flight
- Crew
- GroundEquipment
- Task
- Assignment
- Event

Relationships

Airline
 -> Flight
 -> Task
 -> Assignment
 -> Crew + Equipment

---

# Core Responsibilities

Planner

Receives operational events and determines required work.

Examples:

- FlightArrived
- FlightDelayed

Creates Tasks:

- Fuel
- Cleaning
- Baggage
- Catering
- Pushback

---

Scheduler

Chooses the best assignment.

Simple scoring:

- Equipment compatibility
- Crew availability
- Distance
- Flight priority
- Time until departure

---

Dispatcher

Executes assignments.

Publishes:

- AssignmentCreated
- AssignmentCompleted

Notifies WebSocket clients.

---

# Event Flow

FlightArrived
    ↓
Planner
    ↓
Tasks Created
    ↓
Scheduler
    ↓
Assignments Created
    ↓
Dispatcher
    ↓
WebSocket Notification
    ↓
Metrics Updated

---

# REST API

POST /flights
POST /crew
POST /equipment
POST /events

GET /assignments
GET /dashboard
GET /metrics

---

# Events

FlightArrived

FlightDelayed

CrewAvailable

CrewUnavailable

EquipmentAvailable

EquipmentBroken

TaskCompleted

AssignmentCompleted

---

# In-Memory Storage

Maps protected by RWMutex.

Stores:

Flights

Crews

Equipment

Tasks

Assignments

---

# Concurrency

Use goroutines for:

- Dispatcher
- WebSocket hub
- Event processing
- Airport simulator

Use channels for communication between components.

---

# Metrics

Expose:

- Active flights
- Active crews
- Pending tasks
- Completed tasks
- Average dispatch time

---

# Airport Simulator

Generate random events every few seconds:

- Flight arrivals
- Flight delays
- Equipment failures
- Crew availability
- Task completion

The system should continuously react to these events.

---

# Development Plan

| # | Step | Status |
|---|------|--------|
| 1 | Initialize Go module | ✅ Done |
| 2 | HTTP server (`/health` endpoint, `net/http`) | ✅ Done |
| 3 | Routing (chi router, structured routes) | ⬜ Next |
| 4 | Domain models (Flight, Crew, Equipment, Task, Assignment) | ⬜ Pending |
| 5 | In-memory stores (maps + RWMutex) | ⬜ Pending |
| 6 | Event bus (channels + pub/sub) | ⬜ Pending |
| 7 | Planner (react to events, create tasks) | ⬜ Pending |
| 8 | Scheduler (score and select best crew/equipment) | ⬜ Pending |
| 9 | Dispatcher (execute assignments, update state) | ⬜ Pending |
| 10 | WebSocket hub (broadcast real-time updates) | ⬜ Pending |
| 11 | Metrics endpoint | ⬜ Pending |
| 12 | Airport simulator (random event generator) | ⬜ Pending |
| 13 | Unit tests | ⬜ Pending |

---

# Out of Scope

- Authentication
- Database
- Docker
- Kubernetes
- Kafka
- Redis
- Cloud deployment
- UI (optional)

Focus on backend architecture, concurrency, and clear Go code.
