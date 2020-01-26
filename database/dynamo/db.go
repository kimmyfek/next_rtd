package dynamo

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/kimmyfek/next_rtd/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const (
	maxItems       = 4
	maxRetries     = 3
	connectionsKey = "connections"
)

// TableName is the exported table name for dynamo
const TableName = "rtdgo"

var (
	// ErrUnexpectedTimeRangeFmt err case when data in db is incorrect
	ErrUnexpectedTimeRangeFmt = "time ranges in unexpected format"
	// ErrBadTimeRange is when provided time range doesn't match db time range
	ErrBadTimeRange = "unable to find provided time %s in stored data"
)

// AccessLayer defines the structure that can actually interact with Dynamo
type AccessLayer struct {
	tableName string
	Session   dynamodb.DynamoDB

	routes    map[string]*models.Route
	trips     map[string]*models.Trip
	stops     map[string]*models.Stop
	calendar  map[string]*models.Calendar
	stopTimes map[string][]*models.StopTime
}

// New returns a new AccessLayer
func New(tn string, cfg *aws.Config) (*AccessLayer, error) {
	if tn == "" {
		return nil, fmt.Errorf("Table name must not be empty")
	}

	s := dynamodb.New(session.Must(session.NewSession()), cfg)

	return &AccessLayer{
		tableName: tn,
		Session:   *s,
	}, nil
}

// CreateTables creates new dynamodb tables
func (a *AccessLayer) CreateTables(bool) error { return nil }

// SaveRoutes will, provided a list of routes
// store them on the accesslayer struct to eventually insert them into Dynamo.
// * Note * This has no effect on the actual DB
func (a *AccessLayer) SaveRoutes(_ bool, r map[string]*models.Route) error {
	a.routes = r
	return nil
}

// SaveTrips will, provided a list of trips,
// store them on the accesslayer struct to eventually insert them into Dynamo.
// * Note * This has no effect on the actual DB
func (a *AccessLayer) SaveTrips(_ bool, t map[string]*models.Trip) error {
	a.trips = t
	return nil
}

// SaveStops will, provided a list of stops,
// store them on the accesslayer struct to eventually insert them into Dynamo.
// * Note * This has no effect on the actual DB
func (a *AccessLayer) SaveStops(_ bool, s map[string]*models.Stop) error {
	a.stops = s
	return nil
}

// SaveStopTimes will, provided a list of stop times,
// store them on the accesslayer struct to eventually insert them into Dynamo.
// * Note * This has no effect on the actual DB
func (a *AccessLayer) SaveStopTimes(_ bool, s map[string][]*models.StopTime) error {
	a.stopTimes = s
	return nil
}

// SaveCalendar will, provided a list of calendars,
// store them on the accesslayer struct to eventually insert them into Dynamo.
// * Note * This has no effect on the actual DB
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
	if a.tableName == "" {
		return fmt.Errorf("Dynamo AccessLayer table name is not set and is required")
	}

	var err error
	err = a.saveStops()
	if err != nil {
		return fmt.Errorf("error saving stops: %s", err)
	}

	err = a.saveConnections()
	if err != nil {
		return fmt.Errorf("error saving connections: %s", err)
	}

	return nil
}

// TOOD clean this up
func (a *AccessLayer) saveStops() error {
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

	items := make(map[string][]*dynamodb.WriteRequest)
	for k, v := range d {
		for _, v2 := range v {
			for _, v3 := range v2 {
				sort.Sort(v3)
				for range v3 {
				}
			}
		}

		att, err := dynamodbattribute.MarshalMap(v)
		if err != nil {
			return fmt.Errorf("unable to marshal items into dynamo map: %s", err)
		}

		items[a.tableName] = append(items[a.tableName], &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: map[string]*dynamodb.AttributeValue{
					"path": {
						S: aws.String(k),
					},
					"times": {
						M: att,
					},
					"lastUpdated": {
						S: aws.String(time.Now().Format(time.RFC3339)),
					},
				},
			},
		})

		if len(items[a.tableName]) > maxItems {
			out, err := a.Session.BatchWriteItem(&dynamodb.BatchWriteItemInput{
				RequestItems: items,
			})
			if err != nil {
				return fmt.Errorf("unable to write batches to dynamo: %s", err)
			}

			items = make(map[string][]*dynamodb.WriteRequest)
			if len(out.UnprocessedItems) > 0 {
				if _, ok := out.UnprocessedItems[a.tableName]; ok {
					if len(out.UnprocessedItems[a.tableName]) == maxItems {
						return fmt.Errorf("all items during item processing went unprocessed")
					}
					items[a.tableName] = append(items[a.tableName], out.UnprocessedItems[a.tableName]...)
				}
			}
		}
	}

	for i := 0; i <= maxRetries; i++ {
		if len(items[a.tableName]) > 0 {
			out, err := a.Session.BatchWriteItem(&dynamodb.BatchWriteItemInput{
				RequestItems: items,
			})
			if err != nil {
				fmt.Println(err) // TODO
			}
			if len(out.UnprocessedItems) > 0 {
				if i == maxRetries {
					return fmt.Errorf("error processing last batch of documents")
				}
				time.Sleep(1 * time.Second)
				items = make(map[string][]*dynamodb.WriteRequest)
				items[a.tableName] = append(items[a.tableName], out.UnprocessedItems[a.tableName]...)
				continue
			}
		}
		break
	}

	return nil
}

func (a *AccessLayer) saveConnections() error {
	consTmp := make(map[string]map[string]interface{})
	for _, times := range a.stopTimes {
		for i := 0; i < len(times)-1; i++ {
			for j := i + 1; j < len(times); j++ {
				dST := times[i]
				dStop := a.stops[dST.StopID]

				aST := times[j]
				aStop := a.stops[aST.StopID]

				if _, ok := consTmp[dStop.StopName]; !ok {
					consTmp[dStop.StopName] = make(map[string]interface{})
				}
				consTmp[dStop.StopName][aStop.StopName] = nil

				//if _, ok := consTmp[aStop.StopName]; !ok {
				//	consTmp[aStop.StopName] = make(map[string]interface{})
				//}
				//consTmp[aStop.StopName][dStop.StopName] = nil
			}
		}
	}

	cons := []*models.Station{}
	for s1, mapCons := range consTmp {
		conns := models.Connections{}
		for s2 := range mapCons {
			conns = append(conns, &models.Station{Name: s2})
		}
		sort.Sort(conns)
		c := &models.Station{
			Name:        s1,
			Connections: conns,
		}
		cons = append(cons, c)
	}

	request := make(map[string][]*dynamodb.WriteRequest)
	att, err := dynamodbattribute.MarshalList(cons)
	if err != nil {
		return fmt.Errorf("unable to marshal items into dynamo map: %s", err)
	}

	request[a.tableName] = append(request[a.tableName], &dynamodb.WriteRequest{
		PutRequest: &dynamodb.PutRequest{
			Item: map[string]*dynamodb.AttributeValue{
				"path": {
					S: aws.String(connectionsKey),
				},
				"connections": {
					L: att,
				},
				"lastUpdated": {
					S: aws.String(time.Now().Format(time.RFC3339)),
				},
			},
		},
	})

	out, err := a.Session.BatchWriteItem(&dynamodb.BatchWriteItemInput{
		RequestItems: request,
	})
	if err != nil {
		return fmt.Errorf("unable to write connections to dynamo: %s", err)
	}
	if len(out.UnprocessedItems) > 0 {
		return fmt.Errorf("error writing connections, items unprocessed")
	}

	return nil
}

// DeleteBackupTables removes temp tables
// TODO
func (a *AccessLayer) DeleteBackupTables() error { return nil }

// CreateIndices will create table indexes
// NOOP
func (a *AccessLayer) CreateIndices(bool) error { return nil }

// GetStationsAndConnections retrieves a list of stations from Dynamo
func (a *AccessLayer) GetStationsAndConnections() ([]*models.Station, error) {
	item, err := a.Session.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(a.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"path": {
				S: aws.String(connectionsKey),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	cons := []*models.Station{}
	err = dynamodbattribute.Unmarshal(item.Item["connections"], &cons)
	if err != nil {
		return nil, err
	}

	return cons, nil
}

// GetTimesForStations retrieves a list of upcoming times based on a from station, to station,
// 		the current time, and how many times to retrieve.
func (a *AccessLayer) GetTimesForStations(from, to, now string, numTimes int) ([]*models.Time, error) {
	key := fmt.Sprintf("%s-%s", from, to)

	item, err := a.Session.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(a.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"path": {
				S: aws.String(key),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	var itemTimes map[string]*dynamodb.AttributeValue
	for k := range item.Item["times"].M {
		timeRange := strings.Split(k, "-")
		if len(timeRange) != 2 {
			// TODO should have a log here, don't want to send this to use
			return nil, fmt.Errorf(ErrUnexpectedTimeRangeFmt)
		}
		if now >= timeRange[0] && timeRange[1] > now {
			itemTimes = itemTimes[k].M
			break
		}
	}
	if len(itemTimes) == 0 {
		return nil, fmt.Errorf(ErrBadTimeRange, now)
	}

	for k := range itemTimes {
		fmt.Printf("Key: %s\n", k)
	}

	times := []*models.Time{}

	return times, nil
}
