package parser

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
    "time"
    "archive/zip"
    "path/filepath"
    "net/http"
    "strconv"

	m "github.com/kimmyfek/next_rtd/models"
	log "github.com/sirupsen/logrus"
	set "gopkg.in/fatih/set.v0"
)

// Col Names
var rID = "route_id"
var rShort = "route_short_name"
var rLong = "route_long_name"
var rDesc = "route_desc"

var tID = "trip_id"
var servID = "service_id"
var dirID = "direction_id"

var aTime = "arrival_time"
var dTime = "departure_time"
var sDate = "start_date"
var eDate = "end_date"

var sID = "stop_id"
var sCode = "stop_code"
var sName = "stop_name"
var sDesc = "stop_desc"

// File header columns
var routeC = []string{rID, rShort, rLong, rDesc}
var tripC = []string{rID, servID, tID, dirID}
var sTimeC = []string{tID, aTime, dTime, sID}
var stopC = []string{sID, sCode, sName, sDesc}
var calC = []string{servID, "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday", sDate, eDate}

type db interface {
	CreateTables(bool) error
	SaveRoutes(bool, map[string]m.Route) error
	SaveTrips(bool, map[string]m.Trip) error
	SaveStops(bool, map[string]m.Stop) error
	SaveStopTimes(bool, map[string][]m.StopTime) error
	SaveCalendar(bool, []m.Calendar) error
	SwapTables() error
	DeleteBackupTables() error
	CreateIndices(bool) error
}

// Parser type allows parsing of a filedir
type Parser struct {
	DB      db
	Logger  *log.Entry
}

const (
    scheduleDir = "schedule"
    scheduleFile = "schedule/google_transit.zip"
    scheduleUrl = "http://www.rtd-denver.com/GoogleFeeder/google_transit.zip"
)

// ParseData parses the data files and saves them to the data store.
func (p *Parser) ParseData() {
    // download google transit zip
    downloadSchedule(scheduleUrl, scheduleFile)
    // unzip
    unzipSchedule(scheduleFile, scheduleDir)

    // Parse files
	p.Logger.Info("Parsing Routes File")
	routes := parseRoutes(scheduleDir, "routes")
	p.Logger.Info("Parsing Trips File")
	trips := parseTrips(routes, scheduleDir, "trips")
	p.Logger.Info("Parsing StopTimes File")
	stopTimes, stopIDs := parseStopTimes(trips, scheduleDir, "stop_times")
	p.Logger.Info("Parsing Stops File")
	stops := parseStops(stopIDs, scheduleDir, "stops")
	p.Logger.Info("Parsing Calendar File")
	calData := parseCalendar(scheduleDir, "calendar")

	p.Logger.Info("Creating Temp Tables")
	if err := p.DB.CreateTables(true); err != nil {
		panic(err)
	}

	p.Logger.Info("Saving Routes file")
	if err := p.DB.SaveRoutes(true, routes); err != nil {
		panic(fmt.Sprintf("Unable to save routes to db: %s", err))
	}

	p.Logger.Info("Saving Trips file")
	if err := p.DB.SaveTrips(true, trips); err != nil {
		panic(fmt.Sprintf("Unable to save trips to db: %s", err))
	}

	p.Logger.Info("Saving Stops file")
	if err := p.DB.SaveStops(true, stops); err != nil {
		panic(fmt.Sprintf("Unable to save stops to db: %s", err))
	}

	p.Logger.Info("Saving StopsTimes file")
	if err := p.DB.SaveStopTimes(true, stopTimes); err != nil {
		panic(fmt.Sprintf("Unable to save stop_times to db: %s", err))
	}

	p.Logger.Info("Saving Calendar file")
	if err := p.DB.SaveCalendar(true, calData); err != nil {
		panic(fmt.Sprintf("Unable to save calendar to db: %s", err))
	}

	p.Logger.Info("Swapping Temp and Prod tables")
	if err := p.DB.SwapTables(); err != nil {
		panic(err)
	}

	p.Logger.Info("Deleting prod tables")
	if err := p.DB.DeleteBackupTables(); err != nil {
		panic(err)
	}

	p.Logger.Info("Creating new prod indices")
	if err := p.DB.CreateIndices(false); err != nil {
		panic(err)
	}
}

func getColPos(cols []string, r []string) (map[string]int, error) {
	pos := make(map[string]int)
	for i, c := range r {
		for _, name := range cols {
			if name == c {
				pos[c] = i
			}
		}
		if len(pos) == len(cols) {
			return pos, nil
		}
	}

	return pos, fmt.Errorf("num columns (%s) != num positions (%v) -- row (%s)", cols, pos, r)

}

func parseRoutes(path, filename string) map[string]m.Route {
	filePath := fmt.Sprintf("%s/%s.txt", path, filename)
	f, err := os.Open(filePath)
	if err != nil {
		panic("Unable to open routes.txt")
	}
	defer f.Close()

	routes := make(map[string]m.Route)
	typeIdx := 1000
	r := csv.NewReader(f)

	colPos := make(map[string]int)
	for {
		if route, err := r.Read(); err != nil {
			if err == io.EOF {
				break
			} else {
				panic(fmt.Sprintf("Unable to parse routes: %s", err))
			}
		} else {
			if len(colPos) == 0 {
				colPos, err = getColPos(routeC, route)
				for idx, val := range route {
					if val == "route_type" {
						typeIdx = idx
					}
				}
				if typeIdx == 1000 {
					panic("Type index is not set")
				}
				continue
			}

			if route[typeIdx] == "2" || route[typeIdx] == "0" {
				routes[route[colPos[rID]]] = m.Route{
					RouteID:        route[colPos[rID]],
					RouteShortName: route[colPos[rShort]],
					RouteLongName:  route[colPos[rLong]],
					RouteDesc:      route[colPos[rDesc]],
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

	colPos := make(map[string]int)
	for {
		if trip, err := r.Read(); err != nil {
			if err == io.EOF {
				break
			} else {
				panic(fmt.Sprintf("Unable to parse trips: %s", err))
			}
		} else {
			if len(colPos) == 0 {
				colPos, err = getColPos(tripC, trip)
				continue
			}
			if _, ok := routes[trip[colPos[rID]]]; ok == true {
				trips[trip[colPos[tID]]] = m.Trip{
					RouteID:     trip[colPos[rID]],
					ServiceID:   trip[colPos[servID]],
					TripID:      trip[colPos[tID]],
					DirectionID: trip[colPos[dirID]],
				}
			}
		}
	}
	return trips
}

func parseStopTimes(trips map[string]m.Trip, path, filename string) (map[string][]m.StopTime, set.Interface) {
	filePath := fmt.Sprintf("%s/%s.txt", path, filename)
	f, err := os.Open(filePath)
	if err != nil {
		panic("Unable to open stop_times.txt")
	}
	defer f.Close()

	stopTimes := make(map[string][]m.StopTime)
	stopIDs := set.New(set.ThreadSafe)
	r := csv.NewReader(f)

	colPos := make(map[string]int)
	for {
		if stopTime, err := r.Read(); err != nil {
			if err == io.EOF {
				break
			} else {
				panic(fmt.Sprintf("Unable to parse stop_times: %s", err))
			}
		} else {
			if len(colPos) == 0 {
				colPos, err = getColPos(sTimeC, stopTime)
				continue
			}

			if _, ok := trips[stopTime[colPos[tID]]]; ok == true {
				stopTimes[stopTime[colPos[tID]]] = append(stopTimes[stopTime[colPos[tID]]], m.StopTime{
					TripID:        stopTime[colPos[tID]],
					ArrivalTime:   stopTime[colPos[aTime]],
					DepartureTime: stopTime[colPos[dTime]],
					StopID:        stopTime[colPos[sID]],
				})
				stopIDs.Add(stopTime[colPos[sID]])
			}
		}
	}
	return stopTimes, stopIDs
}

func parseStops(stopIDs set.Interface, path, filename string) map[string]m.Stop {
	filePath := fmt.Sprintf("%s/%s.txt", path, filename)
	f, err := os.Open(filePath)
	if err != nil {
		panic("Unable to open stops.txt")
	}
	defer f.Close()

	stops := make(map[string]m.Stop)
	r := csv.NewReader(f)

	shortNames := map[string]string{
		"Union Station":                                  "Union Station",
		"38th & Blake Station":                           "38th & Blake Station",
		"40th & Colorado Station":                        "40th & Colorado Station",
		"40th Ave & Airport Blvd - Gateway Park Station": "40th Ave & Airport Blvd - Gateway Park Station",
		"61st & Pena Station":                            "61st & Pena Station",
		"Central Park Station":                           "Central Park Station",
		"Peoria Station":                                 "Peoria Station",
		"Theatre District":                               "Theatre District / Convention Center",
		"Westminster Station":                            "Westminster Station",
	}

	colPos := make(map[string]int)
	for {
		if stop, err := r.Read(); err != nil {
			if err == io.EOF {
				break
			} else {
				panic(fmt.Sprintf("Unable to parse stops: %s", err))
			}
		} else {
			if len(colPos) == 0 {
				colPos, err = getColPos(stopC, stop)
				continue
			}

			stopName := stop[colPos[sName]]
			for test, short := range shortNames {
				if strings.Contains(stopName, test) {
					stopName = short
					break
				}
			}
			if ok := stopIDs.Has(stop[colPos[sID]]); ok == true {
				stops[stop[colPos[sID]]] = m.Stop{
					StopID:   stop[colPos[sID]],
					StopCode: stop[colPos[sCode]],
					StopName: stopName,
					StopDesc: stop[colPos[sDesc]],
				}
			}
		}
	}
	return stops
}

func parseCalendar(path, filename string) (cal []m.Calendar) {
    const (
        layoutISO = "20060102 15:04"
    )

	filePath := fmt.Sprintf("%s/%s.txt", path, filename)
	f, err := os.Open(filePath)
	if err != nil {
		panic(fmt.Sprintf("Unable to open %s/%s", path, filename))
	}
	defer f.Close()

	r := csv.NewReader(f)

	colPos := make(map[string]int)
	for {
		if day, err := r.Read(); err != nil {
			if err == io.EOF {
				break
			} else {
				panic(fmt.Sprintf("Unable to parse calendar: %s", err))
			}
		} else {
			if len(colPos) == 0 {
				colPos, err = getColPos(calC, day)
				continue
			}

            // Load MST timezone and make sure datetime is 4 AM
            loc, _ := time.LoadLocation("America/Denver")
            startDate, _ := time.ParseInLocation(layoutISO, day[colPos[sDate]]+" 4:00", loc)
            startUnix := startDate.Unix()
            endDate, _ := time.ParseInLocation(layoutISO, day[colPos[eDate]]+" 4:00", loc)
            // add 24 hours for end date
            endUnix := endDate.Unix() + 86400

			cal = append(cal, m.Calendar{
				ServiceID: day[colPos[sID]],
				Monday:    day[colPos["monday"]],
				Tuesday:   day[colPos["tuesday"]],
				Wednesday: day[colPos["wednesday"]],
				Thursday:  day[colPos["thursday"]],
				Friday:    day[colPos["friday"]],
				Saturday:  day[colPos["saturday"]],
				Sunday:    day[colPos["sunday"]],
				StartDate: strconv.FormatInt(startUnix, 10),
				EndDate:   strconv.FormatInt(endUnix, 10),
			})
		}
	}
	return cal
}


func downloadSchedule(url string, fileName string) error {
    // Get the data
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // Create the file
    out, err := os.Create(fileName)
    if err != nil {
        return err
    }
    defer out.Close()

    // Write the body to file
    _, err = io.Copy(out, resp.Body)
    return err
}


func unzipSchedule(src string, dest string) error {
    var filenames []string
    r, err := zip.OpenReader(src)

    if err != nil {
        return err
    }

    defer r.Close()

    for _, f := range r.File {

        // Store filename/path for returning and using later on
        fpath := filepath.Join(dest, f.Name)

        if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
            return fmt.Errorf("%s: illegal file path", fpath)
        }

        filenames = append(filenames, fpath)

        if f.FileInfo().IsDir() {
            // Make Folder
            os.MkdirAll(fpath, os.ModePerm)
            continue
        }

        // Make File
        if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
            return err
        }

        outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
        if err != nil {
            return err
        }

        rc, err := f.Open()
        if err != nil {
            return err
        }

        _, err = io.Copy(outFile, rc)

        // Close the file without defer to close before next iteration of loop
        outFile.Close()
        rc.Close()

        if err != nil {
            return err
        }
    }
    return nil

}
