package models

// Route defines the data as provided by the route file
type Route struct {
	RouteID        string
	RouteShortName string
	RouteLongName  string
	RouteDesc      string
	RouteType      string
	RouteURL       string
	RouteColor     string
	RouteTextColor string
}

// Trip defines the data as provided by the trips file
type Trip struct {
	RouteID      string
	ServiceID    string
	TripID       string
	TripHeadsign string
	DirectionID  string
	BlockID      string
	ShapeID      string
}

// StopTime defines the data as provided by the StopTimes file
type StopTime struct {
	TripID            string
	ArrivalTime       string
	DepartureTime     string
	StopID            string
	StopSequence      string
	StopHeadsign      string
	PickupType        string
	DropOffType       string
	ShapeDistTraveled string
}

// Stop defines the data as provided by the stops file
type Stop struct {
	StopID             string
	StopCode           string
	StopName           string
	StopDesc           string
	StopLat            string
	StopLon            string
	ZoneID             string
	StopURL            string
	LocationType       string
	ParentStation      string
	StopTimezone       string
	WheelchairBoarding string
}

// Station provides the representation of a train station
type Station struct {
	Name        string    `json:"name"`
	Connections []Station `json:"connections,omitempty"`
}

// Time how a Time payload is represented
type Time struct {
	To          string `json:"to"`
	From        string `json:"from"`
	ArrivalTo   string `json:"arrival_to"`
	ArrivalFrom string `json:"arrival_from"`
	Route       string `json:"route"`
}
