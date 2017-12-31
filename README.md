# Major TODO items
- [X] Save calendar data
- [X] get times based on now
- [X] Only get next times for today
- [X] Figure out how to grab tomorrow if rollover
- [X] Add flags to startup
  - [X] Specify DB path
  - [X] Whether to parse
  - [X] What dir to get parse files from
- [X] Have the UI get served from the back-end
- [X] Combine Stations
  - [X] Union Station
  - [X] 38th & Blake Station
  - [X] 40th & Colorado Station
  - [X] Central Park Station
  - [X] Peoria Station
  - [X] 40th Ave & Airport Blvd-Gateway Park Station
  - [X] 61st & Pena Station
  - [X] Westminster?
- [X] Change "arrival_from" to departure time
- [ ] UI Should Sort stations as |1 2| |3 4| instead of |1 3| |2 4| (K)
- [ ] Fuckin' interfaces -- Need 'em (J)
- [ ] Now is a really bad name for time in db.go (J)
- [ ] Move rtdtime struct (J) 
- [ ] Lots and lots of tests Backend (J)
- [ ] Lots and lots of tests Frontend (K)
- [ ] Dynamically pull data based on column position during parsing to deal with RTD columns being not consistent (J)
- [ ] DNS (K)
- [ ] Build script needs to work (J)
- [ ] Fix initial load (J)
- [ ] Updating deployment currently requires killing the service after a git pull and restarting it. Need to find a better solution, and be able to deploy from local (J / K)
- [ ] Logging > printlns (J)
- [ ] Footer (K)
- [ ] Like, FB, Tweet (K)
- [ ] Contact us / Feedback / About (K)
- [ ] Ad? (J / K)
- [ ] Word light rail (K)
- [ ] SEO (J / K)
- [ ] MySQL? (J / K)
- [ ] Explain Query ? (J / K)
- [ ] Metrics (J)
- [ ] Clean up parser (J)
- [ ] Parser writes to temp table and replaces instead (J)
- [ ] Indexes (J / K)
- [ ] Show arrival time to the _TO_ station (K)
- [ ] Clean up frontend code (K)
- [ ] New react style (K) 
- [ ] Create account, allowing "frequently visited stations" (J + K)
- [ ] Validate data (J / K)
- [ ] Automatically notate if we send null data for times (J)

## Parsing
- [X] Parse the data
- [ ] Re-Pull data and parse after X duration (J)

## Saving
- [X] Save the data to the DB
- [X] Create queries for saving to DB
- [X] Create generic function for saving to DB

## Retrieving Data
- [X] Create an API to serve the data
- [X] Create endpoint to serve all stop locations: GET /stations
- [X] Endpoint to serve all locations with connections
  - [ ] Connections should only be provided if `connections=true` /stations&connections=true (J + K)
- [X] Endpoint to serve all locatons, with connections, with next X incoming times
  - [ ] If `to` isn't provided, show all directions (J + K)
  - [ ] Better error codes and error json responses (J + K)
- [X] Get Times, provide to and from station

## Abstraction
- [ ] Create abstraction layer _AFTER_ completing all of the above

## UI



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

