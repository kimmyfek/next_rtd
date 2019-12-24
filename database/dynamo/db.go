package dynamo

import (
	"github.com/kimmyfek/next_rtd/models"
)

// AccessLayer defines the structure that can actually interact with Dynamo
type AccessLayer struct {
	routes    map[string]*models.Route
	stopTimes map[string][]*models.StopTime
	trips     map[string]*models.Trip
	stops     map[string]*models.Stop
	calendar  []*models.Calendar
}

// New returns a new AccessLayer
func New() (*AccessLayer, error) {
	return &AccessLayer{}, nil
}

// CreateTables creates new dynamodb tables
func (a *AccessLayer) CreateTables(bool) error { return nil }

// SaveRoutes will, provided a list of routes
// store them on the accesslayer struct to eventually insert them into Dynamo.
// * NOTE * This has no effect on the actual DB
func (a *AccessLayer) SaveRoutes(bool, r map[string]*models.Route) error {
	a.routes = r
	return nil
}

// SaveTrips will, provided a list of trips,
// store them on the accesslayer struct to eventually insert them into Dynamo.
// * NOTE * This has no effect on the actual DB
func (a *AccessLayer) SaveTrips(bool, t map[string]*models.Trip) error {
	a.trips = t
	return nil
}

// SaveStops will, provided a list of stops,
// store them on the accesslayer struct to eventually insert them into Dynamo.
// * NOTE * This has no effect on the actual DB
func (a *AccessLayer) SaveStops(bool, s map[string]*models.Stop) error {
	a.stops = s
	return nil
}

// SaveStopTimes will, provided a list of stop times,
// store them on the accesslayer struct to eventually insert them into Dynamo.
// * NOTE * This has no effect on the actual DB
func (a *AccessLayer) SaveStopTimes(bool, s map[string][]*models.StopTime) error {
	a.stopTimes = s
	return nil
}

// SaveCalendar will, provided a list of calendars,
// store them on the accesslayer struct to eventually insert them into Dynamo.
// * NOTE * This has no effect on the actual DB
func (a *AccessLayer) SaveCalendar(bool, c []*models.Calendar) error {
	a.calendar = c
	return nil
}

// SwapTables replaces temp tables and live tables
func (a *AccessLayer) SwapTables() error { return nil }

// DeleteBackupTables removes temp tables
func (a *AccessLayer) DeleteBackupTables() error { return nil }

// CreateIndices will create table indexes
func (a *AccessLayer) CreateIndices(bool) error { return nil }

// GetStationsAndConnections retrieves a list of stations from Dynamo
func (a *AccessLayer) GetStationsAndConnections() ([]models.Station, error) {
	return []models.Station{
		models.Station{Name: "Booty Stop"},
		models.Station{Name: "Alameda Station"},
		models.Station{Name: "Red Light Go!"},
		models.Station{Name: "Veridian Station"},
		models.Station{Name: "Blackwood Station"},
		models.Station{Name: "Wild"},
	}, nil
}

// GetTimesForStations retrieves a list of upcoming times based on a from station, to station,
// 		the current time, and how many times to retrieve.
func (a *AccessLayer) GetTimesForStations(from, to, now string, numTimes int) ([]models.Time, error) {
	return []models.Time{
		models.Time{
			To:            "Wild",
			From:          "Blackwood Station",
			ArrivalTime:   "3",
			DepartureTime: "4",
			Route:         "abc",
		},
	}, nil
}
