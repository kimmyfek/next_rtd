package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"

	"github.com/kimmyfek/next_rtd/database"
	"github.com/kimmyfek/next_rtd/parser"
	"github.com/kimmyfek/next_rtd/web"
)

var debug = flag.Bool("debug", false, "Debug mode & logs")
var port = flag.Int("port", 3000, "Port to host service on")

// Data Parsing Flags
var parse = flag.Bool("parse", false, "If provided will parse data and add to DB. Default: Disabled")
var sourceDir = flag.String("sourceDir", "/runboard", "Dir where source RTD txt files are located. NOT NEEDED IF PARSE=false")

// MySQL Flags
//var sqlPort = flag.String("sqlPort", "", "The for SQL Connstring)
var sqlPass = flag.Bool("sqlPass", false, "Password flag SQL Connstring")
var sqlUser = flag.String("sqlUser", "root", "The user for SQL Connstring")
var sqlHost = flag.String("sqlHost", "localhost", "The hostname for SQL Connstring")
var sqlDB = flag.String("sqlDB", "rtd", "The DB Name for SQL Connstring")

func main() {
	flag.Parse()

	l, err := log.ParseLevel(log.DebugLevel.String())
	if err != nil {
		panic("Incorrect logging level")
	}
	log.SetLevel(l)

	logger := log.WithFields(log.Fields{
		"app": "rtdGO",
	})
	logger.Info("Application Initialization Begin...")
	logger.Debug("Debug mode enabled")

	var sqlPassword string
	if *sqlPass {
		r := bufio.NewReader(os.Stdin)
		fmt.Print("Enter SQL Password: ")
		sqlPassword, err = r.ReadString('\n')
		if err != nil {
			logger.Fatalf("Unable to read password from stdin: %s", err)
		}
	}
	al := database.NewAccessLayer(logger, *sqlUser, sqlPassword, *sqlHost, *sqlDB)
	if err = al.Open(); err != nil {
		panic(fmt.Sprintf("Unable to create and open database: %s", err))
	}
	defer func() {
		logger.Info("Closing DB connection.")
		if err := al.Close(); err != nil {
			logger.Error(fmt.Sprintf("Error shutting down database connection: %s", err))
		}
		logger.Info("DB Connection closed!")
	}()

	if *parse {
		p := &parser.Parser{DB: al, FileDir: *sourceDir, Logger: logger}
		logger.Info("Parsing data and inserting into DB.")
		p.ParseData()
	}

	logger.Info("Init endpoints")
	rh := web.NewRESTHandler(al, logger)
	rh.Init()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		logger.Infof("Listening on %d", *port)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
			logger.Fatalf("Error serving HTTP: %s", err)
		}
	}()

	<-stop
	logger.Warning("Shutting down server!")
}
