---
title: Scheduling sidecar technical addendum
status: draft
created: 2026-10-08
updated: 2026-10-08
---

# Addendum: technical how

Product behavior is in `prd.md`. This file holds mechanism, the README's stack, and how that differs from the code that is already in the repo. Downstream architecture should read this with the PRD. It is not a second requirements list.

## Stack the README requires

- Language and shape: a Go backend, run as a sidecar beside the main application.
- Caller-facing transport: gRPC only. The README says all communication is gRPC requests.
- System of record: MySQL.
- Cache: Redis, for caching. Redis is not the system of record. A read must match the last successful write, so a cached event or availability row cannot be older than that write.
- Load: design for a peak of about 120 requests per second.
- Package sketch from the README, mapped onto this module (`hari.foo/schedulerbackend`), not the `workshops/` root drawn in the README:
  - `cmd/api/main.go` — process entry, wires dependencies, starts the server.
  - `internal/config` — environment and app configuration. Not present in the tree yet.
  - `internal/database` — connections, migrations, pool setup.
  - `internal/event` — event domain: transport, service, store.
  - `internal/contact` — contact domain: transport, service, store. Not present in the tree yet.

The README's ASCII tree labels handlers as HTTP JSON. That drawing disagrees with the gRPC sentence. The PRD treats the gRPC sentence as the requirement. The HTTP labels are not the target surface.

## What the repo contains today

Observed 2026-10-08. This is scaffold, not the v1 contract.

- `cmd/api/main.go` listens with `net/http` on `:8080` and registers event routes. It does not start a gRPC server. It does not connect to Redis.
- `internal/event` implements create via `POST /v1/events`. The stored event has id, title, created time, and modified time. There is no time range, contact, edit, cancel, list-all, or conflict check.
- `internal/database` opens MySQL and pings it. Pool limits in code are 25 open and 25 idle connections, 5-minute max lifetime. Those numbers are current code, not a README requirement.
- There is no `internal/contact`, no `internal/config`, no Redis client, and no protobuf or gRPC dependency in `go.mod`.
- Go version in `go.mod` is 1.26.4.

Database credentials are hardcoded in `cmd/api/main.go`. This addendum does not repeat them. Configuration belongs in `internal/config` once that package exists. Do not copy those literals into new code as the supported setup.

## Decisions this addendum does not make

- Protobuf package names, RPC names, and message fields. Those wait on open questions 3, 4, 6, and 7 in the PRD.
- Whether Redis caches event lists, single events, availability, or all three.
- How a second instance learns of a write. NFR-2 and NFR-3 require that it does; the mechanism is architecture work.
- Migration tool and schema. The README names MySQL and does not name a migration tool. `internal/database` does not run migrations today.
