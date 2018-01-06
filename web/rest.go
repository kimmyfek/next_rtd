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

	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"

	"github.com/kimmyfek/next_rtd/database"
	"github.com/kimmyfek/next_rtd/models"
)

// RestHandler allows dependency injection for REST calls
type RestHandler struct {
	DB       *database.AccessLayer
	Port     int
	stations []models.Station
	Logger   *log.Entry
}

// NewRestHandler returns a new instance of the RestHandler obj
func NewRestHandler(db *database.AccessLayer, port int, logger *log.Entry) *RestHandler {
	return &RestHandler{
		DB:     db,
		Port:   port,
		Logger: logger,
	}
}

// Init defines the API endpoints and starts the HTTP server
func (rh *RestHandler) Init() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	http.Handle("/", http.FileServer(http.Dir(fmt.Sprintf("%s/next", dir))))

	rh.handleFuncWrapper("/stations", rh.GetStations)
	rh.handleFuncWrapper("/times", rh.GetTimes)

	// Cache stations -- This should be something we can reset
	if st, err := rh.DB.GetStationsAndConnections(); err != nil {
		panic(err)
	} else {
		rh.stations = st
	}

	http.ListenAndServe(fmt.Sprintf(":%d", rh.Port), nil)
}

func (rh *RestHandler) handleFuncWrapper(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	h := rh.loggedHandler(http.HandlerFunc(handler))
	h = rh.metricsHandler(h)
	http.Handle(pattern, h)
}

// loggedHandler wraps handlers to log details about each request
func (rh *RestHandler) loggedHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uid := uuid.NewV4()
		logger := rh.Logger.WithFields(log.Fields{
			"context": uid,
			"URI":     r.RequestURI,
		})
		logger.Debug("Incoming request")
		s := time.Now()
		h.ServeHTTP(w, r)
		logger.Debugf("Request duration: %s", time.Now().Sub(s))
	})
}

// metricsHandler wraps handlers to store metrics for each request
func (rh *RestHandler) metricsHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}

// GetStations queries the DB for a list of all stations and then returns them
// to the caller.
// Adding argument "connections=true" will provide all connecting stations.
func (rh *RestHandler) GetStations(w http.ResponseWriter, r *http.Request) {
	// TODO If connections param is true
	j, err := json.Marshal(rh.stations)
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Write(j)
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
	l, _ := time.LoadLocation("MST")
	now := formatTime(time.Now().In(l))
	times, err := rh.DB.GetTimesForStations(from, to, now, numTimes)
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
	return strings.Split(t.Format(time.RubyDate), " ")[3]
}
