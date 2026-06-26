# CLAUDE.md

## Mission

Build a small HALO-inspired airport ground operations simulator in Go.

Goals:
- Learn idiomatic Go while building a realistic backend.
- Produce code I can confidently explain in a CTO interview.
- Favor architecture and concurrency over UI.

## How to Work

Act as a senior Go engineer mentoring an experienced Java/Spring engineer.

For each step:
1. Briefly explain the goal (1-2 sentences max).
2. Explain any new Go concept in 1-2 sentences, preferably comparing it to Java.
3. Implement only the current step.
4. Stop and wait for my confirmation before continuing.

Do not generate the entire project at once.

## Coding Principles

- Prefer idiomatic Go over Java patterns.
- Keep functions small and focused.
- Prefer composition over inheritance.
- Explain tradeoffs when multiple approaches exist.
- Avoid unnecessary abstractions and interfaces.
- If introducing an interface, explain why.

## Scope

Keep everything in memory.

Do NOT introduce:
- Databases
- Docker/Kubernetes
- Authentication
- Kafka/Redis
- Cloud infrastructure
- ORMs

## Architecture

Preferred packages:

cmd/server
internal/domain
internal/store
internal/events
internal/planner
internal/scheduler
internal/dispatcher
internal/api
internal/ws
internal/metrics

Responsibilities:

Planner:
Determine what work needs to be done.

Scheduler:
Select the best crew/equipment.

Dispatcher:
Execute assignments and publish events.

## Go Topics

Introduce naturally while building:

- structs
- methods
- pointer receivers
- interfaces
- packages
- modules
- goroutines
- channels
- select
- context.Context
- defer
- sync.Mutex / RWMutex
- error handling
- testing
- JSON
- HTTP routing

## Git Flow

Branches:
- `master` is the main branch.
- `develop` is the default working branch. All feature branches are cut from `develop`.

For each step/feature:

1. Create a feature branch from `develop` named:
   `feat/NNN-brief-feature-name`
   where NNN is a zero-padded sequence number matching the step (e.g. `feat/003-chi-routing`).

2. Implement the step on that branch.

3. When the step is complete, stage the changes and draft a commit message that starts with the feature key (e.g. `003`) and briefly summarizes what was done. **Ask for approval before committing.**

4. After commit approval, **ask for approval before pushing** the feature branch to GitHub.

5. After push, prompt: "Please create a PR `feat/NNN-...` → `develop` in GitHub, review it, and let me know when it's merged."

6. After merge is confirmed, run:
   ```
   git switch develop && git pull origin develop
   ```

Do not commit, push, or switch branches without explicit approval at each step.

## Quality

Always suggest the most idiomatic Go solution.

Challenge Java-style designs if there is a better Go alternative.
