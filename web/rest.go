package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/kimmyfek/next_rtd/database"
)

func main() {
	fmt.Println("vim-go")
}

// RestHandler allows dependency injection for REST calls
type RestHandler struct {
	DB   *database.AccessLayer
	Port int
}

// NewRestHandler returns a new instance of the RestHandler obj
func NewRestHandler(db *database.AccessLayer, port int) *RestHandler {
	return &RestHandler{
		DB:   db,
		Port: port,
	}
}

// Init defines the API endpoints and starts the HTTP server
func (rh *RestHandler) Init() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	http.Handle("/", http.FileServer(http.Dir(fmt.Sprintf("%s/next", dir)))) //rh.Index)
	http.HandleFunc("/stations", rh.GetStations)
	http.HandleFunc("/times", rh.GetTimes)
	http.ListenAndServe(fmt.Sprintf(":%d", rh.Port), nil)
}

// GetStations queries the DB for a list of all stations and then returns them
// to the caller.
// Adding argument "connections=true" will provide all connecting stations.
func (rh *RestHandler) GetStations(w http.ResponseWriter, r *http.Request) {
	// TODO If connections param is true
	st, err := rh.DB.GetStationsAndConnections()
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		j, err := json.Marshal(st)
		if err != nil {
			w.Write([]byte(err.Error()))
		} else {
			w.Write(j)
		}
	}
}

// GetTimes has several different paths of execution based on arguments.
// Regardless of the arguments, TimeList will return the next four times for
// any given station.
// - No provided args; Will return all stations, a direction, and when the next
// 		four times are for train arrivals at that station.
// - One station; Will return both directions and the next four times for each
// 		of those directions.
// - Two stations; Since the direction is now known, TimesList will only return
//  	the next four times for that singular direction.
func (rh *RestHandler) GetTimes(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	to := vals.Get("to")
	from := vals.Get("from")
	limit := vals.Get("limit")

	if to == "" || from == "" {
		w.Write([]byte(`to and from params must be specified`))
	}

	var numTimes int
	if limit == "" {
		numTimes = 5
	} else {
		var err error
		numTimes, err = strconv.Atoi(limit)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
	}

	// TODO FromTime
	// TODO No from station
	// TODO Days
	now := formatTime(time.Now())
	times, err := rh.DB.GetTimesForStations(to, from, now, numTimes)
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		j, err := json.Marshal(times)
		if err != nil {
			w.Write([]byte(err.Error()))
		} else {
			w.Write(j)
		}
	}
}

func formatTime(t time.Time) string {
	return strings.Split(t.Format(time.Stamp), " ")[2]
}
