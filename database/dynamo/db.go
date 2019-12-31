package dynamo

import (
	"fmt"
	"sort"
	"strings"

	"github.com/kimmyfek/next_rtd/models"
)

// AccessLayer defines the structure that can actually interact with Dynamo
type AccessLayer struct {
	routes    map[string]*models.Route
	trips     map[string]*models.Trip
	stops     map[string]*models.Stop
	calendar  map[string]*models.Calendar
	stopTimes map[string][]*models.StopTime
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
func (a *AccessLayer) SaveRoutes(_ bool, r map[string]*models.Route) error {
	a.routes = r
	return nil
}

// SaveTrips will, provided a list of trips,
// store them on the accesslayer struct to eventually insert them into Dynamo.
// * NOTE * This has no effect on the actual DB
func (a *AccessLayer) SaveTrips(_ bool, t map[string]*models.Trip) error {
	a.trips = t
	return nil
}

// SaveStops will, provided a list of stops,
// store them on the accesslayer struct to eventually insert them into Dynamo.
// * NOTE * This has no effect on the actual DB
func (a *AccessLayer) SaveStops(_ bool, s map[string]*models.Stop) error {
	a.stops = s
	return nil
}

// SaveStopTimes will, provided a list of stop times,
// store them on the accesslayer struct to eventually insert them into Dynamo.
// * NOTE * This has no effect on the actual DB
func (a *AccessLayer) SaveStopTimes(_ bool, s map[string][]*models.StopTime) error {
	a.stopTimes = s
	return nil
}

// SaveCalendar will, provided a list of calendars,
// store them on the accesslayer struct to eventually insert them into Dynamo.
// * NOTE * This has no effect on the actual DB
func (a *AccessLayer) SaveCalendar(_ bool, c []*models.Calendar) error {
	cal := make(map[string]*models.Calendar)
	for _, s := range c {
		cal[s.ServiceID] = s
	}
	a.calendar = cal
	return nil
}

type tripTime struct {
	Departure string `json:"d"`
	Arrival   string `json:"a"`
	Route     string `json:"r"`
}

type tripTimes []*tripTime

// Len is the number of elements in the collection.
func (t tripTimes) Len() int { return len(t) }

// Less reports whether the element with
// index i should sort before the element with index j.
func (t tripTimes) Less(i, j int) bool {
	ti := strings.Split(t[i].Departure, ":")
	tj := strings.Split(t[j].Departure, ":")

	if len(ti) != 3 || len(tj) != 3 {
		fmt.Println("Error parsing departure times")
	}

	for k := 0; k < 3; k++ {
		if ti[k] < tj[k] {
			return true
		} else if ti[k] > tj[k] {
			return false
		}
	}

	ti = strings.Split(t[i].Arrival, ":")
	tj = strings.Split(t[j].Arrival, ":")

	if len(ti) != 3 || len(tj) != 3 {
		fmt.Println("Error parsing arrival times")
	}

	for k := 0; k < 3; k++ {
		if ti[k] < tj[k] {
			return true
		} else if ti[k] > tj[k] {
			return false
		}
	}

	return false
}

// Swap swaps the elements with indexes i and j.
func (t tripTimes) Swap(i, j int) { t[i], t[j] = t[j], t[i] }

// SwapTables coalesces all of the data "Saved" to this access layer and submits
// it to the database. Generates a map with the following structure:
// {
// 		"<station1>-<station2>": {
// 			"<startTime>-<endTime>": {
// 				"<svcID>": [
//					{ "departure: "<dep>", "arrival": "<arrival>" },
//					{ "departure: "<dep>", "arrival": "<arrival>" },
//					{ "departure: "<dep>", "arrival": "<arrival>" }
//				]
//			}
//		}
//	}
func (a *AccessLayer) SwapTables() error {
	d := make(map[string]map[string]map[string]tripTimes)
	for _, times := range a.stopTimes {
		for i := 0; i < len(times)-1; i++ {
			for j := i + 1; j < len(times); j++ {
				dST := times[i]
				dStop := a.stops[dST.StopID]

				aST := times[j]
				aStop := a.stops[aST.StopID]

				key := fmt.Sprintf("%s-%s", dStop.StopName, aStop.StopName)

				trip := a.trips[dST.TripID]
				svc := a.calendar[trip.ServiceID]
				route := a.routes[trip.RouteID]

				activeTimes := fmt.Sprintf("%s-%s", svc.StartDate, svc.EndDate)

				if _, ok := d[key]; !ok {
					d[key] = make(map[string]map[string]tripTimes)
				}

				if _, ok := d[key][activeTimes]; !ok {
					d[key][activeTimes] = make(map[string]tripTimes)
				}

				d[key][activeTimes][svc.ServiceID] = append(d[key][activeTimes][svc.ServiceID], &tripTime{
					Departure: dST.DepartureTime,
					Arrival:   aST.ArrivalTime,
					Route:     route.RouteShortName,
				})
			}
		}
	}

	var i int
	for _, v := range d {
		for _, v2 := range v {
			for _, v3 := range v2 {
				sort.Sort(v3)
				for range v3 {
					i++
				}
			}
		}
	}

	//b, err := json.MarshalIndent(d, "", "  ")
	//if err != nil {
	//	return fmt.Errorf("Error formatting results as JSON")
	//}

	return nil
}

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
