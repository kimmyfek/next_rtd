package main

import (
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

var parse = flag.Bool("parse", false, "If provided will parse data and add to DB. Default: Disabled")
var sourceDir = flag.String("sourceDir", "/runboard", "Dir where source RTD txt files are located. NOT NEEDED IF PARSE=false")
var level = flag.String("level", "info", "Log level. Valid options: 'Info', 'Debug'")
var port = flag.Int("port", 3000, "Port to host service on")

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

	al := database.NewAccessLayer(logger)
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
		logger.Info("Parsing data and inserting into DB.")
		parser.ParseData(al, *sourceDir)
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
