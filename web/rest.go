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

	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"

	"github.com/kimmyfek/next_rtd/models"
)

type db interface {
	GetStationsAndConnections() ([]models.Station, error)
	GetTimesForStations(from, to, now string, numTimes int) ([]models.Time, error)
}

// RESTHandler allows dependency injection for REST calls
type RESTHandler struct {
	DB       db
	stations []models.Station
	Logger   *log.Entry
}

// NewRESTHandler returns a new instance of the RestHandler obj
func NewRESTHandler(d db, logger *log.Entry) *RESTHandler {
	return &RESTHandler{
		DB:     d,
		Logger: logger,
	}
}

// Init defines the API endpoints and starts the HTTP server
func (rh *RESTHandler) Init() {
	rh.Logger.Info("Begin endpoint init and serve files.")
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	http.Handle("/", http.FileServer(http.Dir(fmt.Sprintf("%s/next", dir))))

	rh.Logger.Info("Wrap all non-index endpoints")
	rh.handleFuncWrapper("/stations", rh.GetStations)
	rh.handleFuncWrapper("/times", rh.GetTimes)

	// Cache stations -- This should be something we can reset
	rh.Logger.Info("Cache station data.")
	if st, err := rh.DB.GetStationsAndConnections(); err != nil {
		panic(err)
	} else {
		rh.stations = st
	}
	rh.Logger.Info("All REST init complete!")
}

func (rh *RESTHandler) handleFuncWrapper(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	h := rh.loggedHandler(http.HandlerFunc(handler))
	h = rh.metricsHandler(h)
	http.Handle(pattern, h)
}

type metricsWriter struct {
	http.ResponseWriter
	code int
}

// newMetricsWriter returns a new wrapped metrics responsewriter
func newMetricsWriter(w http.ResponseWriter) *metricsWriter {
	return &metricsWriter{ResponseWriter: w}
}

func (m metricsWriter) WriteHeader(code int) {
	m.code = code
	m.ResponseWriter.WriteHeader(code)
}

func (m metricsWriter) getCode() int {
	if m.code == 0 {
		return 200
	}
	return m.code
}

// loggedHandler wraps handlers to log details about each request
func (rh *RESTHandler) loggedHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := time.Now()
		uid := uuid.NewV4()
		logger := rh.Logger.WithFields(log.Fields{
			"context": uid,
			"URI":     r.RequestURI,
		})

		mw := newMetricsWriter(w)

		logger.Debug("Incoming request")
		h.ServeHTTP(mw, r)
		logger.WithFields(log.Fields{
			"StatusCode": mw.getCode(),
		}).Debugf("Request duration: %s", time.Now().Sub(s))
	})
}

// metricsHandler wraps handlers to store metrics for each request
func (rh *RESTHandler) metricsHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO
		h.ServeHTTP(w, r)
	})
}

// GetStations queries the DB for a list of all stations and then returns them
// to the caller.
// Adding argument "connections=true" will provide all connecting stations.
func (rh *RESTHandler) GetStations(w http.ResponseWriter, r *http.Request) {
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
func (rh *RESTHandler) GetTimes(w http.ResponseWriter, r *http.Request) {
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
	l, _ := time.LoadLocation("America/Denver")
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
