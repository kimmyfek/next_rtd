Instructions:
```
./next_rtd --parse=true --dbDir /opt/next-rtd --sourceDir /Users/jjob200/Downloads/google_transit
```

## Endpoints
### Stations
#### GET
`HTTP/1.1 GET /stations`
v1 -- No params required
Returns stations and connections
```json
[{
	"name": "13th Ave Station",
	"connections": [{
		"name": "2nd \u0026 Abilene Station"
	}, {
		"name": "Colfax Station"
	}, {
		"name": "County Line Station"
	}, {
		"name": "Dry Creek Station"
	}, {
		"name": "Fitzsimons Station"
	}, {
		"name": "Iliff Station"
	}, {
		"name": "Peoria Station"
	}, {
		"name": "Arapahoe at Village Center Station"
	}, {
		"name": "Lincoln Station"
	}, {
		"name": "Aurora Metro Center Station"
	}, {
		"name": "Florida Station"
	}, {
		"name": "Nine Mile Station"
	}, {
		"name": "Orchard Station"
	}, {
		"name": "Belleview Station"
	}, {
		"name": "Dayton Station"
	}]
}, {
	"name": "Aurora Metro Center Station",
	"connections": [{
		"name": "Belleview Station"
	}, {
		"name": "Iliff Station"
	}, {
		"name": "Orchard Station"
	}, {
		"name": "Peoria Station"
	}, {
		"name": "Arapahoe at Village Center Station"
	}, {
		"name": "Colfax Station"
	}, {
		"name": "Nine Mile Station"
	}, {
		"name": "13th Ave Station"
	}, {
		"name": "2nd \u0026 Abilene Station"
	}, {
		"name": "Fitzsimons Station"
	}, {
		"name": "Florida Station"
	}, {
		"name": "Lincoln Station"
	}, {
		"name": "County Line Station"
	}, {
		"name": "Dayton Station"
	}, {
		"name": "Dry Creek Station"
	}]
}, ...]
```

### Times



# Major TODO items
- [X] Save calendar data
- [X] get times based on now
- [X] Only get next times for today
- [ ] Figure out how to grab tomorrow if rollover :: MAJOR PRIORITY ::
- [X] Add flags to startup
  - [X] Specify DB path
  - [X] Whether to parse
  - [X] What dir to get parse files from
- [X] Have the UI get served from the back-end
- [ ] Logging > printlns
- [ ] Clean up parser
- [ ] Parser writes to temp table and replaces instead
- [ ] Fix union stations
- [ ] Fucking interfaces
- [X] Change "arrival_from" to departure time

## Parsing
- [X] Parse the data
- [ ] Re-Pull data and parse after X duration

## Saving
- [X] Save the data to the DB
- [X] Create queries for saving to DB
- [X] Create generic function for saving to DB

## Retrieving Data
- [X] Create an API to serve the data
- [X] Create endpoint to serve all stop locations: GET /stations
- [X] Endpoint to serve all locations with connections
  - [ ] Connections should only be provided if `connections=true` /stations&connections=true
- [X] Endpoint to serve all locatons, with connections, with next X incoming times
  - [ ] If `to` isn't provided, show all directions
  - [ ] Better error codes and error json responses
  - [ ] Handle actual days
	- [ ] Handle now = 11:59, provide tomorrow
- [X] Get Times, provide to and from station

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


-- If time is >= 24, subtract 24
-- If day is Friday, Sat, Sun
