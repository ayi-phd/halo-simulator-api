# HALO Emulator

HALO is an Airport Operations Platform that coordinates real-time ground handling across flights, crew, and equipment. This repository is a backend emulator of that platform — an event-driven, concurrent Go service that models the full ground operations lifecycle from flight arrival through task assignment, dispatch, and resource release.

Built to demonstrate idiomatic Go backend architecture: goroutines, channels, typed pub/sub, and real-time WebSocket streaming — without any external infrastructure dependencies.

## What It Does

When a flight arrives, five ground tasks are automatically created (fueling, cleaning, catering, baggage, pushback). Crew and equipment are scheduled, dispatched, and freed when work completes. All state transitions flow through an in-memory event bus and stream live to a browser dashboard over WebSocket.

```
Flight arrives → Planner creates tasks → Scheduler assigns crew+equipment
             → Dispatcher executes work → Resources freed → next task retried
```

Random operational events (equipment breakdowns, flight delays) are injected by a separate simulator binary, keeping the pipeline under continuous pressure.

## Architecture

The emulator is composed of small, focused components that communicate exclusively through a typed pub/sub event bus — no shared mutable state between components.

| Component    | Responsibility                                               |
|--------------|--------------------------------------------------------------|
| `Planner`    | Reacts to `flight.arrived`, creates one task per operation   |
| `Scheduler`  | Assigns available crew and equipment, retries on availability |
| `Dispatcher` | Executes work concurrently, frees resources on completion    |
| `Bus`        | Buffered fan-out channel — decouples all components          |
| `Hub`        | Broadcasts events to WebSocket clients in real time          |

Two binaries work together:

- **`cmd/halo-server`** — REST API + WebSocket hub + reactive pipeline (the emulator)
- **`cmd/simulator-http`** — external traffic generator; seeds crew/equipment and drives flight arrivals and operational events over HTTP

## Running Locally

**Requirements:** Go 1.21+

```bash
# Clone
git clone https://github.com/ayi-phd/halo-simulator-api
cd halo-simulator-api

# Terminal 1 — start the emulator
go run ./cmd/halo-server

# Terminal 2 — start the simulator
go run ./cmd/simulator-http

# Open the live dashboard
open http://localhost:8080
```

The dashboard connects over WebSocket and shows a color-coded event stream with live metrics. No build step, no Docker, no config files.

## REST API

| Method | Path           | Description                          |
|--------|----------------|--------------------------------------|
| GET    | `/`            | Live dashboard (HTML)                |
| GET    | `/metrics`     | Snapshot: task counts, dispatch time |
| GET    | `/flights`     | All flights                          |
| GET    | `/equipment`   | All ground equipment                 |
| GET    | `/assignments` | All crew assignments                 |
| GET    | `/dashboard`   | Flights paired with tasks + metrics  |
| POST   | `/flights`     | Manually arrive a flight             |
| POST   | `/crew`        | Register a crew member               |
| POST   | `/equipment`   | Register a piece of equipment        |
| POST   | `/events`      | Inject an operational event          |
| GET    | `/ws`          | WebSocket event stream               |

## Key Go Concepts

- **Goroutines and channels** — each pipeline component runs in its own goroutine and communicates via typed channels
- **`select` with `ctx.Done()`** — every goroutine exits cleanly on shutdown
- **`sync.RWMutex`** — all in-memory stores are safe for concurrent read/write
- **`signal.NotifyContext`** — graceful shutdown on `SIGINT`/`SIGTERM` with HTTP drain
- **`//go:embed`** — dashboard HTML is bundled into the server binary at compile time
- **Table-driven tests** with channel-based timeouts — no `time.Sleep` in tests

## Project Structure

```
cmd/
  halo-server/      Emulator binary (API + reactive pipeline)
  simulator-http/   Simulator binary (external HTTP traffic generator)
internal/
  domain/           Core types (Flight, Task, Crew, Equipment, Assignment)
  events/           Typed pub/sub bus
  store/            In-memory stores (map + RWMutex)
  planner/          Task creation from flight events
  scheduler/        Crew and equipment assignment
  dispatcher/       Concurrent work execution
  simclient/        HTTP client used by simulator-http
  api/              Chi router, REST handlers, embedded dashboard
  ws/               WebSocket hub and client pump
  metrics/          On-demand snapshot collector
```
