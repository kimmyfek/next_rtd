package parser

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/kimmyfek/next_rtd/database"
	m "github.com/kimmyfek/next_rtd/models"
	set "gopkg.in/fatih/set.v0"
)

// ParseData parses the data files and saves them to the data store.
// TODO Merge union stations together
// TODO stops with only a couple options? look at that data
func ParseData(db *database.AccessLayer, fileDir string) {
	routes := parseRoutes(fileDir, "routes")
	trips := parseTrips(routes, fileDir, "trips")
	stopTimes, stopIDs := parseStopTimes(trips, fileDir, "stop_times") // TODO could switch back to stopid for key
	stops := parseStops(stopIDs, fileDir, "stops")
	calData := parseCalendar(fileDir, "calendar")

	// TODO fix the funcs
	err := db.SaveRoutes("routes", routes) // Replace with replacedata
	if err != nil {
		panic(fmt.Sprintf("Unable to save routes to db: %s", err))
	}
	err = db.SaveTrips("trips", trips) // Replace with replacedata
	if err != nil {
		panic(fmt.Sprintf("Unable to save trips to db: %s", err))
	}
	err = db.SaveStops("stops", stops) // Replace with replacedata
	if err != nil {
		panic(fmt.Sprintf("Unable to save stops to db: %s", err))
	}
	err = db.SaveStopTimes("stop_times", stopTimes) // Replace with replacedata
	if err != nil {
		panic(fmt.Sprintf("Unable to save stop_times to db: %s", err))
	}
	err = db.SaveCalendarData("calendar", calData)
	if err != nil {
		panic(fmt.Sprintf("Unable to save calendar to db: %s", err))
	}

}

func parseRoutes(path, filename string) map[string]m.Route {
	filePath := fmt.Sprintf("%s/%s.txt", path, filename)
	f, err := os.Open(filePath)
	if err != nil {
		panic("Unable to open routes.txt")
	}
	defer f.Close()

	routes := make(map[string]m.Route)
	var typeIdx int
	r := csv.NewReader(f)

	for {
		if route, err := r.Read(); err != nil {
			if err == io.EOF {
				break
			} else {
				panic(fmt.Sprintf("Unable to parse routes: %s", err))
			}
		} else {
			if route[0] == "route_id" {
				for idx, val := range route {
					if val == "route_type" {
						typeIdx = idx
					}
				}
				if typeIdx == 0 {
					panic("Type index is not set")
				}
				continue
			}
			if route[typeIdx] == "2" || route[typeIdx] == "0" {
				routes[route[0]] = m.Route{
					RouteID:        route[0],
					RouteShortName: route[1],
					RouteLongName:  route[2],
					RouteDesc:      route[3],
					RouteType:      route[4],
					RouteURL:       route[5],
					RouteColor:     route[6],
					RouteTextColor: route[7],
				}
			}
		}
	}

	return routes
}

func parseTrips(routes map[string]m.Route, path, filename string) map[string]m.Trip {
	filePath := fmt.Sprintf("%s/%s.txt", path, filename)
	f, err := os.Open(filePath)
	if err != nil {
		panic("Unable to open trips.txt")
	}
	defer f.Close()

	trips := make(map[string]m.Trip)
	r := csv.NewReader(f)

	for {
		if trip, err := r.Read(); err != nil {
			if err == io.EOF {
				break
			} else {
				panic(fmt.Sprintf("Unable to parse trips: %s", err))
			}
		} else {
			if _, ok := routes[trip[0]]; ok == true {
				trips[trip[2]] = m.Trip{
					RouteID:      trip[0],
					ServiceID:    trip[1],
					TripID:       trip[2],
					TripHeadsign: trip[3],
					DirectionID:  trip[4],
					BlockID:      trip[5],
					ShapeID:      trip[6],
				}
			}
		}
	}
	return trips
}

func parseStopTimes(trips map[string]m.Trip, path, filename string) (map[string][]m.StopTime, *set.Set) {
	filePath := fmt.Sprintf("%s/%s.txt", path, filename)
	f, err := os.Open(filePath)
	if err != nil {
		panic("Unable to open stop_times.txt")
	}
	defer f.Close()

	stopTimes := make(map[string][]m.StopTime)
	stopIDs := set.New()
	r := csv.NewReader(f)

	for {
		if stopTime, err := r.Read(); err != nil {
			if err == io.EOF {
				break
			} else {
				panic(fmt.Sprintf("Unable to parse stop_times: %s", err))
			}
		} else {
			if _, ok := trips[stopTime[0]]; ok == true {
				stopTimes[stopTime[0]] = append(stopTimes[stopTime[0]], m.StopTime{
					TripID:            stopTime[0],
					ArrivalTime:       stopTime[1],
					DepartureTime:     stopTime[2],
					StopID:            stopTime[3],
					StopSequence:      stopTime[4],
					StopHeadsign:      stopTime[5],
					PickupType:        stopTime[6],
					DropOffType:       stopTime[7],
					ShapeDistTraveled: stopTime[8],
				})
				stopIDs.Add(stopTime[3])
			}
		}
	}
	return stopTimes, stopIDs
}

func parseStops(stopIDs *set.Set, path, filename string) map[string]m.Stop {
	filePath := fmt.Sprintf("%s/%s.txt", path, filename)
	f, err := os.Open(filePath)
	if err != nil {
		panic("Unable to open stops.txt")
	}
	defer f.Close()

	stops := make(map[string]m.Stop)
	r := csv.NewReader(f)

	shortNames := []string{
		"Union Station",
		"38th & Blake Station",
		"40th & Colorado Station",
		"40th Ave & Airport Blvd - Gateway Park Station",
		"61st & Pena Station",
		"Central Park Station",
		"Peoria Station",
		"Theatre District",
		"Westminster Station",
	}

	for {
		if stop, err := r.Read(); err != nil {
			if err == io.EOF {
				break
			} else {
				panic(fmt.Sprintf("Unable to parse stops: %s", err))
			}
		} else {
			stopName := stop[2]
			for _, short := range shortNames {
				if strings.Contains(stopName, short) {
					stopName = short
					break
				}
			}
			if ok := stopIDs.Has(stop[0]); ok == true {
				stops[stop[0]] = m.Stop{
					StopID:             stop[0],
					StopCode:           stop[1],
					StopName:           stopName,
					StopDesc:           stop[3],
					StopLat:            stop[4],
					StopLon:            stop[5],
					ZoneID:             stop[6],
					StopURL:            stop[7],
					LocationType:       stop[8],
					ParentStation:      stop[9],
					StopTimezone:       stop[10],
					WheelchairBoarding: stop[11],
				}
			}
		}
	}
	return stops
}

func parseCalendar(path, filename string) (cal []m.Calendar) {
	filePath := fmt.Sprintf("%s/%s.txt", path, filename)
	f, err := os.Open(filePath)
	if err != nil {
		panic(fmt.Sprintf("Unable to open %s/%s", path, filename))
	}
	defer f.Close()

	r := csv.NewReader(f)

	for {
		if day, err := r.Read(); err != nil {
			if err == io.EOF {
				break
			} else {
				panic(fmt.Sprintf("Unable to parse calendar: %s", err))
			}
		} else {
			cal = append(cal, m.Calendar{
				ServiceID: day[0],
				Monday:    day[1],
				Tuesday:   day[2],
				Wednesday: day[3],
				Thursday:  day[4],
				Friday:    day[5],
				Saturday:  day[6],
				Sunday:    day[7],
				StartDate: day[8],
				EndDate:   day[9],
			})
		}
	}
	return cal
}
