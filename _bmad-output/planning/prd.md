---
title: Scheduling sidecar
status: draft
created: 2026-10-08
updated: 2026-10-08
---

# PRD: Scheduling sidecar

*Working title — confirm.*

## 0. Document Purpose

This PRD is for the builder of the main application and for downstream BMad architecture, spec, and ticket work. It states what the scheduling sidecar must do. Vocabulary is fixed in the Glossary. Features group functional requirements with stable IDs. Inferences are tagged `[ASSUMPTION]` and indexed in §9. Transport, storage, and the current code scaffold live in `addendum.md`, not here. The source is `README.md`. There is no UX document.

## 1. Vision

The main application needs scheduling without owning it. This product is a sidecar process, written in Go, that the main application calls for events, contacts, availability, conflict checks, and logs.

The sidecar exists so the main application can ask "can this happen, and record that it did" without embedding a scheduler. v1 is that service and nothing else: no screens, no accounts, no product of its own.

## 2. Target User

### 2.1 Jobs To Be Done

- When the main application must put something on a schedule, it creates an event and learns whether that time is free.
- When a person's time changes, the main application updates availability and trusts later conflict checks to use it.
- When something on the schedule changes or is called off, the main application edits or cancels the event and the conflict picture updates.
- When the main application needs a record of what the sidecar did, it records a log.

### 2.2 Non-Users (v1)

- End users of the main application. They never call this sidecar.
- Operators who would need a UI, login, or admin console.

### 2.3 Key User Journeys

Pure service. Journeys are one sentence each. The actor is the main application.

- **UJ-1.** The main application creates an event for a contact and is told it was stored or that it conflicts.
- **UJ-2.** The main application lists events and contacts it previously stored.
- **UJ-3.** The main application edits an event, or cancels it, and later conflict checks ignore the cancelled event.
- **UJ-4.** The main application sets a contact's availability, then a conflict check uses that availability.
- **UJ-5.** The main application records a log of a scheduling action.

## 3. Glossary

- **Main application** — the other project this sidecar serves. The only caller in v1.
- **Sidecar** — this scheduling backend. It does not replace the main application.
- **Event** — one scheduled item. Fields: title, capacity, date, start time, end time. Every event has one organizer and may have many participants.
- **Contact** — a person. Fields: name, email, phone. [ASSUMPTION: availability belongs to a contact.]
- **Inactive** — the status a successful contact delete sets. The contact remains stored. A contact who is organizer or participant on a future event cannot be set inactive.
- **Organizer** — the one contact required on an event, in the organizer role.
- **Participant** — a contact associated with an event besides the organizer. An event may have many.
- **Availability** — a block of time on a date, belonging to a contact, with a start and an end. The smallest unit is 15 minutes, so 3:00 PM–10:00 PM on a date is one availability and 3:10 PM is not a valid boundary. [ASSUMPTION: update replaces rather than merges.]
- **Conflict** — a proposed event that must not be stored as a success when its organizer or any participant is the organizer or a participant of an overlapping active event, or is outside their availability. No gap is required: a start time equal to another event's end time is not a conflict. [ASSUMPTION: availability of those same contacts counts.]
- **Active event** — an event that still occupies time. A cancelled event is not active.
- **Cancel** — mark an event so it no longer occupies time. [ASSUMPTION: cancel is not a delete; the event remains readable.]
- **Log** — a stored record of a scheduling action. [ASSUMPTION: the main application submits the log; the README does not say the sidecar writes logs on its own.]
- **Caller** — the main application, for the length of one request.

## 4. Features

### 4.1 Events

**Description:** The caller lists, creates, edits, and cancels events, and checks a proposal for conflict before treating it as booked. Realizes UJ-1, UJ-2, UJ-3.

**Functional Requirements:**

#### FR-1: List events

The caller can get all events.

**Consequences (testable):**
- After one or more successful creates, a list returns each of those events.
- A list of zero events succeeds and returns an empty collection.

**Out of Scope:**
- Filters, paging, and sort order. [ASSUMPTION: v1 returns the full set. The README says "get all events."]

#### FR-2: Create an event

The caller can create an event. Realizes UJ-1.

**Consequences (testable):**
- A successful create is visible to a later list (FR-1).
- A create that conflicts is not stored as an active event (FR-5).
- A create missing title, capacity, date, start time, end time, or the organizer is rejected and not stored.

#### FR-3: Edit an event

The caller can edit an event. Realizes UJ-3.

**Consequences (testable):**
- After a successful edit, a later list returns the edited values.
- An edit that would conflict is rejected and leaves the previous values in place.
- Editing an unknown event is rejected.

#### FR-4: Cancel an event

The caller can cancel an event. Realizes UJ-3.

**Consequences (testable):**
- After a successful cancel, conflict checks do not treat that event as occupying time.
- The cancelled event remains available to list. [ASSUMPTION: cancel does not remove the row.]
- Cancelling an unknown event, or an event that is already cancelled, is rejected and does not change other events.

#### FR-5: Check for conflict

The caller can check a proposed event for conflict. Realizes UJ-1.

**Consequences (testable):**
- A proposal whose organizer or any participant is the organizer or a participant of an overlapping active event is reported as a conflict.
- A proposal that starts when another active event ends is not a conflict. No gap is required between them.
- A proposal is a conflict when its organizer or any participant is outside their availability. [ASSUMPTION: availability of those same contacts counts.]
- A proposal that overlaps only a cancelled event is not a conflict for that reason.
- A check does not itself create an event.
- Create (FR-2) and edit (FR-3) use the same conflict rule and refuse to store a conflicting result.

**Notes:**
- [NOTE FOR PM] No gap is required, and an exact endpoint touch is not a conflict. Both the organizer and the participants count.

### 4.2 Contacts and availability

**Description:** The caller lists, creates, and deletes contacts, and adds or updates a contact's availability. Realizes UJ-2, UJ-4.

**Functional Requirements:**

#### FR-6: List contacts

The caller can get all contacts.

**Consequences (testable):**
- After one or more successful creates, a list returns each of those contacts.
- A list of zero contacts succeeds and returns an empty collection.
- Whether an inactive contact appears in the list is still open.

#### FR-7: Create a contact

The caller can create a contact.

**Consequences (testable):**
- A successful create is visible to a later list (FR-6).
- A create missing name, email, or phone is rejected and not stored.

#### FR-8: Delete a contact

The caller can delete a contact.

**Consequences (testable):**
- A delete sets the contact inactive. It does not remove the contact.
- A delete is rejected, and the status stays unchanged, while the contact is organizer or participant on an event whose start time is still in the future.
- Events are not rewritten when a contact becomes inactive.
- Deleting an unknown contact is rejected.

#### FR-9: Add or update availability

The caller can add or update a contact's availability. Realizes UJ-4.

**Consequences (testable):**
- After a successful update, a conflict check (FR-5) uses the new availability.
- A block whose start or end is not on a 15-minute boundary is rejected and not stored.
- Updating availability for an unknown contact is rejected.
- A second update for the same contact replaces the availability conflict checks use. [ASSUMPTION: update is a replacement, not a merge.]

### 4.3 Logs

**Description:** The caller records logs. Realizes UJ-5.

**Functional Requirements:**

#### FR-10: Record a log

The caller can record a log.

**Consequences (testable):**
- A successful record stores the log the caller submitted.
- A record missing whatever the sidecar requires in a log is rejected and stores nothing.

**Out of Scope:**
- Listing or searching logs. The README says "record logs" and does not say the caller can read them back. [NOTE FOR PM] If the caller cannot read logs, FR-10 cannot be demonstrated except by inspecting storage.

### 4.4 Feature-specific quality

These bind a feature tighter than the cross-cutting NFRs.

- Event list and single-event lookup used by conflict checks stay correct after a create, edit, or cancel on any instance (see NFR-2).
- Availability used by FR-5 is the availability last successfully written by FR-9 (see NFR-2).

## 5. Non-Goals (Explicit)

- This sidecar does not provide a UI.
- This sidecar does not authenticate callers and does not own user accounts.
- This sidecar is not the main application and does not take on that application's other jobs.
- v1 does not add scheduling operations beyond FR-1 through FR-10.

## 6. MVP Scope

### 6.1 In Scope

- FR-1 through FR-10.
- Conflict checks that agree with stored events and availability.
- Serving the caller while more than one sidecar instance is running.
- Holding a peak of about 120 requests per second.

### 6.2 Out of Scope for MVP

- UI and auth, because the README forbids both.
- A read API for logs, until the open question in FR-10 is answered.
- Filters, search, and notifications. The README does not ask for them.

## 7. Success Metrics

The README states a load expectation, not a numeric latency target. Targets below that are not in the source stay open.

**Primary**

- **SM-1**: The caller can complete UJ-1 through UJ-5 against one shared store. Validates FR-1 through FR-10.
- **SM-2**: At a peak of about 120 requests per second, creates and conflict checks still follow FR-5. Validates FR-2, FR-5, NFR-3.

**Secondary**

- **SM-3**: Event lookup p99 is under 50 ms. Validates NFR-1.

**Counter-metrics (do not optimize)**

- **SM-C1**: Do not count a faster lookup as success if it serves an event or availability that disagrees with the last successful write. Counterbalances SM-3. Validates NFR-2.

## 8. Cross-cutting NFRs

- **NFR-1 Low-latency event lookup.** Event lookup completes with p99 under 50 ms.
- **NFR-2 Consistency.** A read matches the last successful write. A cache must not return an event or availability older than that write. Two successful bookings of the same contact time are not allowed.
- **NFR-3 Horizontal scale.** More than one instance can serve the caller at the same time against the same events and availability.
- **NFR-4 No UI and no auth.** The sidecar exposes no end-user interface and performs no authentication.
- **NFR-5 Load.** The design holds a peak of about 120 requests per second. The README does not split that number by operation.

## 9. Open Questions

1. Does a later availability update replace the contact's blocks or merge with them?
2. What is in a log, and must the caller be able to read logs back?
3. Does "get all" mean every row in the store, or every row for one tenant or one owner? The README does not mention tenancy.
4. Does listing contacts include inactive contacts?

## 10. Assumptions Index

- §3 Event — title, capacity, date, start time, end time, one organizer, and many participants are confirmed. A minimum number of participants is not.
- §3 Contact — a contact is a person with name, email, and phone. Confirmed.
- §3 Availability — a dated block of time on a 15-minute grid is confirmed. Whether update replaces or merges is still open.
- §3 Conflict — the organizer and the participants both count. No gap is required. Confirmed. Availability of those same contacts counting is still an assumption.
- §3 Cancel — cancel keeps the event readable and drops it from conflict consideration.
- §3 Log — the caller submits the log.
- §4 FR-1 — list returns the full set, with no paging.
- §4 FR-8 — a delete sets the contact inactive and does not rewrite events. Confirmed. A future event, judged by start time, blocks the delete. Confirmed. Whether a list includes inactive contacts is still open.
- §1 and addendum — the caller-facing transport is gRPC, because the README prose says all communication is gRPC. The ASCII tree's HTTP handlers are not the target.
