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
- [X] Change "arrival_from" to departure time
- [X] Fix time showing in wrong TZ on deployed server by changing TZ.
- [X] Fix bug causing panic. (Bug was caused by extra space coming in from date formatting. Decided to change from `time.Stamp` to `time.RubyDate` which shouldn't add another space, messing up parsing.
- [X] Cache Station Data
- [X] Indexes (J)
- [X] Explain Query ? (J / K)
- [X] Startup flags should be pointing to prod dirs (J)
- [x] Automatically notate if we send null data for times (J)
- [X] Not changing time when querying for the next day
- [X] DNS (K)
- [x] Only show connection stations of that direction (J)
- [x] Logging > printlns (J)
- [x] Log status of each request (J)
- [x] Graceful shutdown (J)
- [x] Theatre district stop should display as "Theatre District / Convention Center" (J)
- [x] Build script needs to work (J)
- Building
 - [ ] Build JS from container
 - [ ] Docker image for apache + go + js
 - [ ] Script to deploy sql + go to remote
- [ ] Add some logs to DB rollover scripts (J)
- [X] Footer (K)
- [ ] Like, FB, Tweet (K)
- [ ] Word light rail (K)
- [X] Show arrival time to the _TO_ station (K)
- [ ] Contact us / Feedback / About (K / J)
- [X] Queries hitting the backend twice (K)
- [X] Make bundle.css smaller (K)
- [X] Make app.js smaller (K)
- [ ] Ad? (J / K)
- [ ] SEO (J / K)
- [ ] Validate data (J / K)
- [ ] HTTP Server like Apache (J / K)
- [ ] Rename everything to rtdgo.co (J / K)
- Performance Optimizations
  - MySQL? (J)
    - [x] Run DB script
    - [x] Replace sqlite with mysql in db.go
    - [x] Validate queries are working appropriately
    - [x] Move indexes to outside of check db.
    - [x] Build ConnString on CLI / ENV
    - [x] Password secret on CLI / ENV for service
    - [ ] Password secret for Docker run
	- [x] Explain query
- Parser Improvements (J)
  - [x] Parser writes to temp table and replaces instead (J)
  - [ ] Stuttering on parser.parser
  - Choose one
	- [ ] Check calendar start-end date, and make sure today falls in that value. If not, parse
	- [ ] Re-Pull data and parse each day at 4am
  - [ ] Determine which file to download & unzip
  - [ ] Update DB live
  - [ ] Refresh cache on reload
  - [x] Parser has logger
  - [x] Dynamically pull data based on column position during parsing to deal with RTD columns being not consistent

## Notes
- Start over button more apparent?
- Scroll to top after clicking an option or grey out the station you clicked and all non-connection stations?

## Server Cleanup
- [X] UI Should Sort stations as |1 2| |3 4| instead of |1 3| |2 4| (K)
- Interfaces
  - [ ] db.go -> sql interface ?
  - [ ] Create log package and work on interfacing
- [ ] Package renaming
- [ ] Now is a really bad name for time in db.go (J)
- [ ] Move rtdtime struct (J)
- [ ] Lots and lots of tests Backend (J)
- [ ] Clean up parser (J)
- [ ] Metrics (J)
- [ ] DB area is a mess of constants and functions that could likely be broken into funcs

## UI Cleanup
- [ ] Clean up frontend code (K)
- [X] New react style (K)
- [ ] Lots and lots of tests Frontend (K)

## Proper API Handling
- More RESTful
  - [ ] If `to` isn't provided, show all directions (J + K)
  - [ ] Better error codes and error json responses (J + K)
  - [ ] Connections should only be provided if `connections=true` /stations&connections=true (J + K)

## Advanced Features
- [ ] Create account, allowing "frequently visited stations" (J + K)

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

