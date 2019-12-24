package dynamo

import (
	"github.com/kimmyfek/next_rtd/models"
)

// AccessLayer defines the structure that can actually interact with Dynamo
type AccessLayer struct{}

// New returns a new AccessLayer
func New() (*AccessLayer, error) {
	return &AccessLayer{}, nil
}

// CreateTables creates new dynamodb tables
func (a *AccessLayer) CreateTables(bool) error { return nil }

// SaveRoutes will, provided a list of routes, insert them into Dynamo
func (a *AccessLayer) SaveRoutes(bool, map[string]models.Route) error { return nil }

// SaveTrips will, provided a list of trips, insert them into Dynamo
func (a *AccessLayer) SaveTrips(bool, map[string]models.Trip) error { return nil }

// SaveStops will, provided a list of stops, insert them into Dynamo
func (a *AccessLayer) SaveStops(bool, map[string]models.Stop) error { return nil }

// SaveStopTimes will, provided a list of stop times, insert them into Dynamo
func (a *AccessLayer) SaveStopTimes(bool, map[string][]models.StopTime) error { return nil }

// SaveCalendar will, provided a list of calendars, insert them into Dynamo
func (a *AccessLayer) SaveCalendar(bool, []models.Calendar) error { return nil }

// SwapTables replaces temp tables and live tables
func (a *AccessLayer) SwapTables() error { return nil }

// DeleteBackupTables removes temp tables
func (a *AccessLayer) DeleteBackupTables() error { return nil }

// CreateIndices will create table indexes
func (a *AccessLayer) CreateIndices(bool) error { return nil }

// GetStationsAndConnections retrieves a list of stations from Dynamo
func (a *AccessLayer) GetStationsAndConnections() ([]models.Station, error) { return nil, nil }

// GetTimesForStations retrieves a list of upcoming times based on a from station, to station,
// 		the current time, and how many times to retrieve.
func (a *AccessLayer) GetTimesForStations(from, to, now string, numTimes int) ([]models.Time, error) {
	return nil, nil
}
