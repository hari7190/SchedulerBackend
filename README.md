# Scheduling Application Backend

Goal for this repo is to  build a Scheduling application backend in Go for another project I am building. This application will be setup like a sidecar with the main application to handle the scheduling needs for the main app. All communication will be done via gRPC requests.


- Get all events
- Create an event
- Edit/Cancel event
- Check for conflict
- Get all contacts
- Create a contact
- Delete contact
- Add/update contact availability
- record logs

Database is MySQL and Redis for caching.

Peak load is expected to be around 120 req/s.

Non functional requirements are:

1. Low latency event lookup.
2. Data consistency for event and availability.
3. Must be scalable horizontally.