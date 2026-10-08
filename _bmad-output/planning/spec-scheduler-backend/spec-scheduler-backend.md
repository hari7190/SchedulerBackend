---
id: SPEC-scheduler-backend
companions:
  - ../addendum.md
sources:
  - ../prd.md
---

> **Canonical contract.** This SPEC and the files in `companions:` are the complete, preservation-validated contract for what to build, test, and validate. Source documents listed in frontmatter are for traceability — consult them only if you need narrative rationale or prose color this contract intentionally omits.

# Scheduling sidecar

## Why

The main application needs a separate scheduling sidecar so it can store events and contacts, check conflicts, record availability, and record logs. The force is a vision to realize: that sidecar, called by the main application, at about 120 requests per second, without a UI or authentication of its own.

## Capabilities

- **CAP-1**
  - **intent:** The caller can list events so it can see what is scheduled.
  - **success:** A list returns every stored event, and an empty store returns an empty collection.
- **CAP-2**
  - **intent:** The caller can create an event so the schedule includes it.
  - **success:** A later list includes the event, and a conflicting create is not stored as active.
- **CAP-3**
  - **intent:** The caller can edit an event so the schedule matches the change.
  - **success:** A later list shows the edited values, and a conflicting edit leaves the previous values.
- **CAP-4**
  - **intent:** The caller can cancel an event so it no longer occupies time.
  - **success:** Later conflict checks ignore it, and a list still returns it.
- **CAP-5**
  - **intent:** The caller can check a proposed event for conflict so it does not double-book.
  - **success:** A time overlap is reported when the organizer or any participant is the organizer or a participant of another active event. A start time equal to another event's end time is not a conflict. The check does not create an event.
- **CAP-6**
  - **intent:** The caller can list contacts so it can see who can be scheduled.
  - **success:** A list returns every stored contact. Whether inactive contacts are included is open.
- **CAP-7**
  - **intent:** The caller can create a contact so that party can be scheduled.
  - **success:** A later list includes the contact.
- **CAP-8**
  - **intent:** The caller can delete a contact so that person becomes inactive.
  - **success:** The contact's status becomes inactive, events are unchanged, and the delete is rejected while that contact is organizer or participant on an event whose start time is still in the future.
- **CAP-9**
  - **intent:** The caller can add or update a contact's availability so conflict checks use current time.
  - **success:** A later conflict check uses the availability from the last successful update.
- **CAP-10**
  - **intent:** The caller can record a log so a scheduling action is stored.
  - **success:** The submitted log is stored. Whether the caller can read it back is an open question.

## Constraints

- The only caller-facing transport is gRPC. HTTP is not the product surface.
- MySQL is the system of record for events, contacts, availability, and logs.
- Redis may cache. It must not be the only copy of an event or of availability.
- A read matches the last successful write. A cached event or availability must not be older than that write. Two successful bookings of the same contact time are not allowed.
- The service has no UI and does not authenticate callers.
- More than one instance must be able to serve the same data.
- The design must hold a peak of about 120 requests per second.
- Event lookup p99 is under 50 ms.

## Non-goals

- A UI, or authentication and accounts inside this service.
- Replacing the main application.
- Scheduling operations other than CAP-1 through CAP-10.
- A log read API, until the open question on logs is answered.

## Success signal

The main application, over gRPC, can list, create, edit, and cancel events, check conflicts, list, create, and delete contacts, update availability, and record logs. Conflict checks agree with the last successful event and availability writes. Two instances can serve that at a peak of about 120 requests per second, with no UI and no auth in this service.

## Assumptions

- An event's fields are title, capacity, date, start time, and end time.
- A contact is a person with name, email, and phone.
- Every event has one organizer contact and may have many participant contacts.
- Cancel keeps the event readable and removes it from conflict consideration.
- Availability is a dated time block on a 15-minute grid, for example 3:00 PM–10:00 PM on a date. Whether update replaces or merges is open.
- No gap is required between events. A start time equal to another event's end time is not a conflict.
- A conflict counts both the organizer and the participants.
- The caller submits logs; the sidecar does not invent them.
- "Get all" means every stored row, not a tenant slice.

## Open Questions

- Does an availability update replace the contact's blocks or merge with them?
- Does listing contacts include inactive contacts?
- Must recorded logs be readable by the caller?
