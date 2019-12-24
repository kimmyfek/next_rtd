package sql

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"gopkg.in/fatih/set.v0"

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
	AL     *sql.DB
	logger *log.Entry
	pass   string
	user   string
	host   string
	db     string
}

// NewAccessLayer is the function provided to instantiate a new instance of the
// AccessLayer object and connect to the DB.
func NewAccessLayer(logger *log.Entry, user, pass, host, db string) *AccessLayer {
	return &AccessLayer{
		logger: logger,
		user:   user,
		pass:   pass,
		host:   host,
		db:     db,
	}
}

// Open begins the connection with the db
func (al *AccessLayer) Open() error {
	c := mysql.Config{
		User:   al.user,
		Passwd: al.pass,
		DBName: al.db,
		Addr:   al.host,
	}
	if al.host != "" {
		c.Net = "tcp"
	}
	c.Params = map[string]string{"allowNativePasswords": "true"}
	fmt.Println(c.FormatDSN())
	db, err := sql.Open("mysql", c.FormatDSN())
	if err != nil {
		return err
	}
	db.SetConnMaxLifetime(5 * time.Minute)
	if err := db.Ping(); err != nil {
		return err
	}

	al.AL = db

	err = al.CreateTables(false)
	if err != nil {
		return err
	}

	err = al.CreateIndices(false)
	if err != nil {
		return err
	}

	return nil

}

// CreateTables instantiates all tables if they exist.
// If "temp" == true, all table names will have '_temp' appended to name
func (al *AccessLayer) CreateTables(temp bool) error {
	tableMap := map[string]string{
		routesTable:    createRoutesTable,
		stopsTable:     createStopsTable,
		stopTimesTable: createStopTimesTable,
		tripsTable:     createTripsTable,
		calendarTable:  createCalendarTable,
	}
	for t, q := range tableMap {
		t := getTableName(t, temp)
		if al.tableExists(t) {
			continue
		}
		query := fmt.Sprintf(q, t)
		_, err := al.AL.Exec(query)
		if err != nil {
			al.logger.Errorf("Error creating %s table: %s", t, query)
			return err
		}
		al.logger.Infof("Created %s table", t)
	}

	return nil
}

// CreateIndices creates the indexes for all tables
func (al *AccessLayer) CreateIndices(temp bool) error {

	idxMap := map[string][]string{
		stopsTable:     []string{"stop_name"},
		stopTimesTable: []string{"trip_id", "arrival_time"},
		tripsTable:     []string{"trip_id", "service_id"},
	}
	for t, indices := range idxMap {
		for _, c := range indices {
			t := getTableName(t, temp)
			i := fmt.Sprintf("idx_%s", c)
			if al.indexExists(t, i) {
				continue
			}
			q := fmt.Sprintf("CREATE INDEX %s ON %s(%s)", i, t, c)
			_, err := al.AL.Exec(q)
			if err != nil {
				al.logger.Errorf("Error creating index %s against %s", i, t)
				return err
			}
			al.logger.Infof("Created index %s against %s", i, t)
		}
	}

	return nil
}

func getTableName(t string, temp bool) string {
	if temp {
		t = fmt.Sprintf("temp_%s", t)
	}
	return t

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

// SaveRoutes stores Route m to the DB for each entry in the provided list
func (al *AccessLayer) SaveRoutes(temp bool, data map[string]m.Route) error {
	table := getTableName(routesTable, temp)
	if len(data) == 0 {
		return fmt.Errorf("Unable to save data to %s, empty list provided", table)
	}
	stmt := fmt.Sprintf(`INSERT INTO %s(
		route_id, route_short_name, route_long_name, route_desc) VALUES `, table)
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
func (al *AccessLayer) SaveTrips(temp bool, data map[string]m.Trip) error {
	table := getTableName(tripsTable, temp)
	if len(data) == 0 {
		return fmt.Errorf("Unable to save data to %s, empty list provided", table)
	}
	stmt := fmt.Sprintf(`INSERT INTO %s(
		route_id, service_id, trip_id, direction_id) VALUES `, table)
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
func (al *AccessLayer) SaveStops(temp bool, data map[string]m.Stop) error {
	table := getTableName(stopsTable, temp)
	if len(data) == 0 {
		return fmt.Errorf("Unable to save data to %s, empty list provided", table)
	}
	stmt := fmt.Sprintf(`INSERT INTO %s(
		stop_id, stop_code, stop_name, stop_desc) VALUES `, table)
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
func (al *AccessLayer) SaveStopTimes(temp bool, data map[string][]m.StopTime) error {
	table := getTableName(stopTimesTable, temp)
	if len(data) == 0 {
		return fmt.Errorf("Unable to save data to %s, empty list provided", table)
	}
	stmt := fmt.Sprintf(`INSERT INTO %s(
		trip_id, arrival_time, departure_time, stop_id) VALUES `, table)
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

// SaveCalendar stores Calendar m to the DB for each entry in the provided list
func (al *AccessLayer) SaveCalendar(temp bool, data []m.Calendar) error {
	table := getTableName(calendarTable, temp)
	if len(data) == 0 {
		return fmt.Errorf("Unable to save data to %s, empty list provided", table)
	}
	stmt := fmt.Sprintf(`INSERT INTO %s(
		service_id, monday, tuesday, wednesday, thursday, friday,
		saturday, sunday, start_date, end_date) VALUES `, table)
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

// SwapTables replaces all tables with their temp name equivelants
func (al *AccessLayer) SwapTables() error {
	query := "RENAME TABLE "
	for i, t := range []string{routesTable, stopsTable, stopTimesTable, tripsTable, calendarTable} {
		if i != 0 {
			query += ", "
		}
		temp := getTableName(t, true)
		bkup := fmt.Sprintf("%s_bkup", t)
		query += fmt.Sprintf(" %s TO %s, %s TO %s ", t, bkup, temp, t)
	}
	if _, err := al.AL.Exec(query); err != nil {
		return fmt.Errorf("Unable to swap temp tables for live tables: %s", err)
	}

	return nil
}

// DeleteBackupTables wipes out all backup tables
func (al *AccessLayer) DeleteBackupTables() error {
	query := fmt.Sprintf("DROP TABLES %s", strings.Join([]string{
		fmt.Sprintf("%s_bkup", routesTable),
		fmt.Sprintf("%s_bkup", stopsTable),
		fmt.Sprintf("%s_bkup", stopTimesTable),
		fmt.Sprintf("%s_bkup", tripsTable),
		fmt.Sprintf("%s_bkup", calendarTable)}, ", "))
	if _, err := al.AL.Exec(query); err != nil {
		return fmt.Errorf("Unable to delete temp tables: %s", err)
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

	connSet := make(map[string]set.Interface)
	defer rows.Close()
	for rows.Next() {
		var stop, connection string
		if err := rows.Scan(&stop, &connection); err != nil {
			return nil, err
		}

		if conns, ok := connSet[stop]; ok {
			conns.Add(connection)
		} else {
			s := set.New(set.ThreadSafe)
			s.Add(connection)
			connSet[stop] = s

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
	al.logger.Debugf("Query: %s", query)
	al.logger.Debugf("From: %s", from)
	al.logger.Debugf("To: %s", to)
	al.logger.Debugf("Now: %s", now)
	al.logger.Debugf("Numtimes: %d", numTimes)
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
	if intH < 10 {
		h = fmt.Sprintf("0%s", h)
	}
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
		al.logger.Warningf(fmt.Sprintf("%d times requested, only %d provided", numTimes, len(times)))
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
