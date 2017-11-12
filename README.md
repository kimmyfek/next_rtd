- [ ] Save calendar data
- [ ] get times based on now
- [ ] Only get next times for today
- [ ] Figure out how to grab tomorrow if rollover
- [X] Add flags to startup
  - [X] Specify DB path
  - [X] Whether to reload DB
- [ ] Have the UI get served from the back-end
- [ ] DB data replace
- [ ] Logging > printlns

# TODO
## Parsing
- [X] Parse the data
- [ ] Re-Pull data and parse after X duration

## Saving
- [X] Save the data to the DB
- [X] Create queries for saving to DB
- [X] Create generic function for saving to DB

## Retrieving Data
- [X] Create an API to serve the data
- [ ] Serve the UI?
- [X] Create endpoint to serve all stop locations: GET /stations
- [X] Endpoint to serve all locations with connections
  - [ ] Connections should only be provided if `connections=true` /stations&connections=true
- [X] Endpoint to serve all locatons, with connections, with next X incoming times
  - [ ] If `to` isn't provided, show all directions
  - [ ] Better error codes and error json responses
  - [ ] Handle actual days
	- [ ] Handle now = 11:59, provide tomorrow
- [ ] Get Times, provide to and from station

## Abstraction
- [ ] Create abstraction layer _AFTER_ completing all of the above
  - [ ] Abstract save routes
  - [ ] Abstract save times
  - [ ] Abstract save 3
  - [ ] Abstract save 4
  - [ ] Abstract get stations
  - [ ] Abstract get stations with connections
  - [ ] Abstract get all times
  - [ ] Abstract get 1 station times
  - [ ] Abstract get 2 station times

## UI
- [ ] Show arrival time to the _TO_ station
