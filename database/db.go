package database

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gopkg.in/fatih/set.v0"

	_ "github.com/go-sql-driver/mysql" // SQL doesn't need a name
	m "github.com/kimmyfek/next_rtd/models"
)

const (
	chunkSize      int    = 500
	routesTable    string = "routes"
	stopsTable     string = "stops"
	stopTimesTable string = "stop_times"
	tripsTable     string = "trips"
	calendarTable  string = "calendar"
)

var serviceIDMap = map[string]string{
	time.Weekday(0).String(): "sunday",
	time.Weekday(1).String(): "monday",
	time.Weekday(2).String(): "tuesday",
	time.Weekday(3).String(): "wednesday",
	time.Weekday(4).String(): "thursday",
	time.Weekday(5).String(): "friday",
	time.Weekday(6).String(): "saturday",
}

// AccessLayer is the object meant to be used to access the DB
type AccessLayer struct {
	AL *sql.DB
}

// NewAccessLayer is the function provided to instantiate a new instance of the
// AccessLayer object and connect to the DB.
func NewAccessLayer() *AccessLayer {
	return &AccessLayer{}
}

// Open begins the connection with the db
func (al *AccessLayer) Open() error {
	db, err := sql.Open("mysql", "root@/rtd")
	if err != nil {
		return err
	}
	al.AL = db

	if !al.tableExists(routesTable) {
		// TODO rename table columns to not have table name in it
		_, err := al.AL.Exec(fmt.Sprintf(`
			CREATE TABLE %s(
				route_id 		 VARCHAR(10) NOT NULL PRIMARY KEY,
				route_short_name VARCHAR(5) NOT NULL,
				route_long_name  VARCHAR(255) NOT NULL,
				route_desc		 VARCHAR(255) NOT NULL
			)
		`, routesTable))

		if err != nil {
			fmt.Printf("Error creating %s table", routesTable)
			return err
		}
		fmt.Println("Created routes table")
	}

	if !al.tableExists(stopsTable) {
		_, err := al.AL.Exec(fmt.Sprintf(`
			CREATE TABLE %s(
				stop_id 	    INT NOT NULL PRIMARY KEY,
				stop_code		INT NOT NULL,
				stop_name		VARCHAR(55) NOT NULL,
				stop_desc		VARCHAR(55) NOT NULL
			)
		`, stopsTable))

		if err != nil {
			fmt.Printf("Error creating %s table", stopsTable)
			return err
		}
		fmt.Println("Created stops table")
	}

	if !al.indexExists(stopsTable, "idx_stop_name") {
		_, err = al.AL.Exec("CREATE INDEX idx_stop_name ON stops(stop_name)")
		if err != nil {
			fmt.Printf("Error creating %s index", stopsTable)
			return err
		}
		fmt.Println("Created idx_stop_name index")
	}
	if !al.tableExists(stopTimesTable) {
		_, err := al.AL.Exec(fmt.Sprintf(`
			CREATE TABLE %s(
				trip_id			INT NOT NULL,
				arrival_time	VARCHAR(9) NOT NULL,
				departure_time	VARCHAR(9) NOT NULL,
				stop_id			INT NOT NULL
			)
		`, stopTimesTable))

		if err != nil {
			fmt.Printf("Error creating %s table", stopTimesTable)
			return err
		}
		fmt.Println("Created stop times table")
	}

	if !al.indexExists(stopTimesTable, "idx_trip_id") {
		_, err = al.AL.Exec("CREATE INDEX idx_trip_id ON stop_times(trip_id)")
		if err != nil {
			fmt.Printf("Error creating %s index", stopTimesTable)
			return err
		}
		fmt.Println("Created idx_trip_id index on StopTimesTable")
	}

	if !al.indexExists(stopTimesTable, "idx_arrival_time") {
		_, err = al.AL.Exec("CREATE INDEX idx_arrival_time ON stop_times(arrival_time)")
		if err != nil {
			fmt.Printf("Error creating %s index", stopTimesTable)
			return err
		}
		fmt.Println("Created idx_arrival_time index")
	}

	if !al.tableExists(tripsTable) {
		_, err := al.AL.Exec(fmt.Sprintf(`
			CREATE TABLE %s(
				route_id 	 VARCHAR(10) NOT NULL,
				service_id	 VARCHAR(15) NOT NULL,
				trip_id		 INT NOT NULL,
				direction_id INT NOT NULL
			)
		`, tripsTable))

		if err != nil {
			fmt.Printf("Error creating %s table", tripsTable)
			return err
		}
		fmt.Println("Created trips table")
	}

	if !al.indexExists(tripsTable, "idx_trip_id") {
		_, err = al.AL.Exec("CREATE INDEX idx_trip_id ON trips(trip_id)")
		if err != nil {
			fmt.Printf("Error creating %s index", tripsTable)
			return err
		}
		fmt.Println("Created idx_trip_id index")
	}

	if !al.indexExists(tripsTable, "idx_service_id") {
		_, err = al.AL.Exec("CREATE INDEX idx_service_id ON trips(service_id)")
		if err != nil {
			fmt.Printf("Error creating %s index", tripsTable)
			return err
		}
		fmt.Println("Created idx_service_id index")
	}

	if !al.tableExists(calendarTable) {
		_, err := al.AL.Exec(fmt.Sprintf(`
            CREATE TABLE %s(
                service_id VARCHAR(15) NOT NULL,
                monday     VARCHAR(15) NOT NULL,
                tuesday    VARCHAR(15) NOT NULL,
                wednesday  VARCHAR(15) NOT NULL,
                thursday   VARCHAR(15) NOT NULL,
                friday     VARCHAR(15) NOT NULL,
                saturday   VARCHAR(15) NOT NULL,
                sunday     VARCHAR(15) NOT NULL,
                start_date VARCHAR(10) NOT NULL,
                end_date   VARCHAR(10) NOT NULL
            )
        `, calendarTable))

		if err != nil {
			fmt.Printf("Error creating %s table", calendarTable)
			return err
		}
		fmt.Println("Created calendar table")

	}
	return nil
}

// Close executes the close functionality against the DB.
func (al *AccessLayer) Close() error {
	return al.AL.Close()
}

func (al *AccessLayer) tableExists(name string) bool {
	var n string
	err := al.AL.QueryRow(`
		SELECT TABLE_NAME
		FROM information_schema.tables
		WHERE table_schema = 'rtd'
		  AND table_name = ?
		LIMIT 1`, name).Scan(&n)
	if err == nil && n == name {
		return true
	}
	return false
}

func (al *AccessLayer) indexExists(table, index string) bool {
	var i string
	err := al.AL.QueryRow(`
	SELECT DISTINCT INDEX_NAME
	FROM INFORMATION_SCHEMA.STATISTICS
	WHERE TABLE_SCHEMA = 'rtd'
		AND TABLE_NAME = ?
		AND INDEX_NAME = ?`, table, index).Scan(&i)
	if err == nil && i == index {
		return true
	}
	return false
}

// GetStations retrieves a list of stations from the DB and returns them
func (al *AccessLayer) GetStations() []string {
	return []string{}
}

// ReplaceData deletes all of the data that existed prior to a dump and
// replaces that data with the newly dumped data. In order for that to work,
// ReplaceData must do the following.
// 1. Accepts lists of all of the newly dumped objects
//   - If any of the lists are empty, the dump errors out
// 2. Creates temporary tables for each of the newly dumped fields that are
// 		duplicates of the existing tables.
// 3. Dumps the data to the temp tables.
// 4. Lock the live tables, rename them, rename temp tables to live names, remove lock
func (al *AccessLayer) ReplaceData() error {
	return nil
}

// SaveRoutes stores Route m to the DB for each entry in the provided list
func (al *AccessLayer) SaveRoutes(table string, data map[string]m.Route) error {
	if len(data) == 0 {
		return fmt.Errorf("Unable to save routes, empty list provided")
	}
	stmt := `INSERT INTO routes (
		route_id, route_short_name, route_long_name, route_desc) VALUES `
	values := []string{}

	for _, val := range data {
		values = append(values, fmt.Sprintf(
			"('%s', '%s', '%s', '%s') ",
			val.RouteID,
			val.RouteShortName,
			val.RouteLongName,
			val.RouteDesc,
		))
	}
	rowsAffected, err := al.exec(stmt, values)
	if err != nil {
		return err
	}
	if rowsAffected != int64(len(data)) {
		return fmt.Errorf("%d rows inserted, yet %d records provided",
			rowsAffected,
			len(data))
	}
	return nil
}

// SaveTrips stores Trip m to the DB for each entry in the provided list
func (al *AccessLayer) SaveTrips(table string, data map[string]m.Trip) error {
	if len(data) == 0 {
		return fmt.Errorf("Unable to save trips, empty list provided")
	}
	stmt := `INSERT INTO trips (
		route_id, service_id, trip_id, direction_id) VALUES `
	values := []string{}

	for _, val := range data {
		values = append(values, fmt.Sprintf(
			"('%s', '%s', '%s', '%s') ",
			val.RouteID,
			val.ServiceID,
			val.TripID,
			val.DirectionID,
		))
	}

	rowsAffected, err := al.exec(stmt, values)
	if err != nil {
		return err
	}
	if rowsAffected != int64(len(data)) {
		return fmt.Errorf("%d rows inserted, yet %d records provided",
			rowsAffected,
			len(data))
	}
	return nil
}

// SaveStops stores Trip m to the DB for each entry in the provided list
func (al *AccessLayer) SaveStops(table string, data map[string]m.Stop) error {
	if len(data) == 0 {
		return fmt.Errorf("Unable to save stop times, empty list provided")
	}
	stmt := `INSERT INTO stops(
		stop_id, stop_code, stop_name, stop_desc) VALUES `
	values := []string{}

	for _, val := range data {
		stopID, err := strconv.Atoi(val.StopID)
		if err != nil {
			return err
		}
		stopCode, err := strconv.Atoi(val.StopCode)
		if err != nil {
			return err
		}
		values = append(values, fmt.Sprintf(
			"(%d, %d, '%s', '%s') ",
			stopID,
			stopCode,
			val.StopName,
			val.StopDesc,
		))
	}

	rowsAffected, err := al.exec(stmt, values)
	if err != nil {
		return err
	}
	if rowsAffected != int64(len(data)) {
		return fmt.Errorf("%d rows inserted, yet %d records provided",
			rowsAffected,
			len(data))
	}
	return nil
}

// SaveStopTimes stores Trip m to the DB for each entry in the provided list
func (al *AccessLayer) SaveStopTimes(table string, data map[string][]m.StopTime) error {
	if len(data) == 0 {
		return fmt.Errorf("Unable to save stops, empty list provided")
	}
	stmt := `INSERT INTO stop_times(
		trip_id, arrival_time, departure_time, stop_id) VALUES `
	values := []string{}

	var numProvided int64
	for _, stopTimes := range data {
		for _, val := range stopTimes {
			tripID, err := strconv.Atoi(val.TripID)
			if err != nil {
				return err
			}
			stopID, err := strconv.Atoi(val.StopID)
			if err != nil {
				return err
			}
			values = append(values, fmt.Sprintf(
				"(%d, '%s', '%s', %d) ",
				tripID,
				val.ArrivalTime,
				val.DepartureTime,
				stopID,
			))
			numProvided++
		}
	}

	rowsAffected, err := al.exec(stmt, values)
	if err != nil {
		return err
	}
	if rowsAffected != numProvided {
		return fmt.Errorf("%d rows inserted, yet %d records provided",
			rowsAffected,
			numProvided)
	}
	return nil
}

// SaveCalendarData stores Calendar m to the DB for each entry in the provided list
func (al *AccessLayer) SaveCalendarData(table string, data []m.Calendar) error {
	if len(data) == 0 {
		return fmt.Errorf("Unable to calendar data, empty list provided")
	}
	stmt := `INSERT INTO calendar(
		service_id, monday, tuesday, wednesday, thursday, friday,
		saturday, sunday, start_date, end_date) VALUES `
	values := []string{}

	var numProvided int64
	for _, day := range data {
		values = append(values, fmt.Sprintf(
			"('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s') ",
			day.ServiceID,
			day.Monday,
			day.Tuesday,
			day.Wednesday,
			day.Thursday,
			day.Friday,
			day.Saturday,
			day.Sunday,
			day.StartDate,
			day.EndDate,
		))
		numProvided++
	}

	rowsAffected, err := al.exec(stmt, values)
	if err != nil {
		return err
	}
	if rowsAffected != numProvided {
		return fmt.Errorf("%d rows inserted, yet %d records provided",
			rowsAffected,
			numProvided)
	}
	return nil
}

func (al *AccessLayer) exec(query string, values []string) (int64, error) {
	var affectedRows int64
	for chunk := 0; chunk <= len(values); chunk += chunkSize {
		endChunk := chunk + chunkSize
		if endChunk > len(values) {
			endChunk = len(values)
		}
		stmt := query + strings.Join(values[chunk:endChunk], ", ")

		result, err := al.AL.Exec(stmt)
		if err != nil {
			return 0, err
		}
		affected, err := result.RowsAffected()
		if err != nil {
			return 0, err
		}
		affectedRows += affected
	}
	return affectedRows, nil
}

// GetStationsAndConnections returns all stations along with their connections
func (al *AccessLayer) GetStationsAndConnections() ([]m.Station, error) {
	query := `
		SELECT DISTINCT
			s.stop_name, s2.stop_name
		FROM stops s
		INNER JOIN stop_times st
			ON st.stop_id = s.stop_id
		INNER JOIN trips t
			ON t.trip_id = st.trip_id
		INNER JOIN routes r
			ON r.route_id = t.route_id
		INNER JOIN (
			SELECT DISTINCT
				s.stop_name,
				r.route_short_name,
				r.route_id,
				st.arrival_time,
				t.trip_id
			FROM stops s
			INNER JOIN stop_times st
				ON st.stop_id = s.stop_id
			INNER JOIN trips t
				ON t.trip_id = st.trip_id
			INNER JOIN routes r
				ON r.route_id = t.route_id
		) s2
		ON r.route_id = s2.route_id
		  AND s.stop_name != s2.stop_name
		  AND t.trip_id = s2.trip_id
		  AND st.arrival_time < s2.arrival_time
	`

	rows, err := al.AL.Query(query)
	if err != nil {
		return nil, err
	}

	connSet := make(map[string]*set.Set)
	defer rows.Close()
	for rows.Next() {
		var stop, connection string
		if err := rows.Scan(&stop, &connection); err != nil {
			return nil, err
		}

		if conns, ok := connSet[stop]; ok {
			conns.Add(connection)
		} else {
			connSet[stop] = set.New(connection)
		}
	}

	connections := []m.Station{}
	for stop, conns := range connSet {
		connList := []m.Station{}
		conns.Each(func(item interface{}) bool {
			connList = append(connList, m.Station{Name: item.(string)})
			return true
		})
		connections = append(connections, m.Station{
			Name:        stop,
			Connections: connList,
		})
	}

	return connections, nil
}

func (al *AccessLayer) getStationTimes(from, to, now, day string, numTimes int) (*sql.Rows, error) {
	query := fmt.Sprintf(`
		SELECT DISTINCT
			s.stop_name, -- from
			s2.stop_name, -- to
			st.arrival_time AS departure_time, -- from time
			s2.arrival_time, -- to time
			r.route_short_name
		FROM stops s
		INNER JOIN stop_times st
			ON st.stop_id = s.stop_id
		INNER JOIN trips t
			ON t.trip_id = st.trip_id
		INNER JOIN routes r
			ON r.route_id = t.route_id
		INNER JOIN calendar c
			ON c.service_id = t.service_id
		INNER JOIN (
			SELECT DISTINCT
				s.stop_name,
				st.arrival_time,
				t.trip_id,
				r.route_short_name,
				r.route_id
			FROM stops s
			INNER JOIN stop_times st ON st.stop_id = s.stop_id
			INNER JOIN trips t ON t.trip_id = st.trip_id
			INNER JOIN routes r ON r.route_id = t.route_id
			INNER JOIN calendar c ON c.service_id = t.service_id
		) s2
			ON r.route_id = s2.route_id
			AND s.stop_name != s2.stop_name
			AND t.trip_id = s2.trip_id

		WHERE
			s.stop_name = ?
			AND s2.stop_name = ?
			AND departure_time < s2.arrival_time -- From time < To time
			AND departure_time > ?
			AND c.%s = 1
		ORDER BY departure_time ASC
		LIMIT ?
	`, day)
	return al.AL.Query(query, from, to, now, numTimes)
}

type rtdTime struct {
	h string
	m string
	s string
}

func newRTDTime(t string) *rtdTime {
	c := strings.Split(t, ":")
	return &rtdTime{h: c[0], m: c[1], s: c[2]}
}

func (r *rtdTime) toStringNextDay() string {
	intH, _ := strconv.Atoi(r.h)
	var h, m string
	if intH <= 23 {
		h = "00"
		m = "00"
	} else {
		h = r.h
		m = r.m
	}
	return fmt.Sprintf("%s:%s:%s", h, m, r.s)
}

func (r *rtdTime) toStringRTDTime() string {
	intH, _ := strconv.Atoi(r.h)
	if intH <= 4 {
		intH += 24
	}
	h := strconv.Itoa(intH)
	return fmt.Sprintf("%s:%s:%s", h, r.m, r.s)
}

func (r *rtdTime) hourDelta(delta int) {
	// This should be incapable of error since it's an hour time
	intH, _ := strconv.Atoi(r.h)
	intH += delta
	r.h = strconv.Itoa(intH)
}

// GetTimesForStations returns a list of time slots between two train stations
// Since times are weird in RTD data (24:00 under the Friday slot would actually
// 	be midnight on Saturday), some time manipulation needs to be done in order
// 	to get the appropriate results from the database.
// The database is queried for 'RTD Time' (if < 5am, add 24 and day == yesterday)
// if fewer than numTimes results are found, a second query is executed to
// get the remaining times
func (al *AccessLayer) GetTimesForStations(from, to, now string, numTimes int) ([]m.Time, error) {
	var times []m.Time
	t := newRTDTime(now)
	// NOTE: Day might be something that's passed in
	day := al.getServiceIDFromDay(0)
	rows, err := al.getStationTimes(from, to, t.toStringRTDTime(), day, numTimes)
	if err != nil {
		return nil, err
	}
	times = append(times, parseStationTimeRows(rows)...)

	if len(times) < numTimes {
		day = al.getServiceIDFromDay(24 * time.Hour)
		rows, err := al.getStationTimes(from, to, t.toStringNextDay(), day, numTimes-len(times))
		if err != nil {
			return nil, err
		}
		times = append(times, parseStationTimeRows(rows)...)

	}

	if len(times) < numTimes {
		fmt.Println(fmt.Sprintf("%d times requested, only %d provided", numTimes, len(times)))
	}
	return times, nil
}

func parseStationTimeRows(rows *sql.Rows) []m.Time {
	var times []m.Time
	defer rows.Close()
	for rows.Next() {
		var to, from, arrivalTime, departureTime, route string
		rows.Scan(&from, &to, &departureTime, &arrivalTime, &route)
		times = append(times, m.Time{
			From:          from,
			To:            to,
			DepartureTime: departureTime,
			ArrivalTime:   arrivalTime,
			Route:         route,
		})
	}

	return times
}

// getServiceIDFromDay determines the day of the week based on now, adjusted by delta.
// Due to a day by RTD standards going up to a maximum of 27 o'clock, if today
// is still before 5 AM, the clock gets rolled back to yesterday.
func (al *AccessLayer) getServiceIDFromDay(delta time.Duration) string {
	l, _ := time.LoadLocation("MST")
	now := time.Now().In(l).Add(delta)
	if now.Hour() <= 4 {
		now = now.Add(-5 * time.Hour)
	}
	return serviceIDMap[now.Weekday().String()]
}
