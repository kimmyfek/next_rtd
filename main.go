package main

import (
	"flag"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/kimmyfek/next_rtd/database"
	"github.com/kimmyfek/next_rtd/parser"
	"github.com/kimmyfek/next_rtd/web"
)

var parse = flag.Bool("parse", false, "If provided will parse data and add to DB. Default: Disabled")
var sourceDir = flag.String("sourceDir", "/runboard", "Dir where source RTD txt files are located. NOT NEEDED IF PARSE=false")
var level = flag.String("level", "info", "Log level. Valid options: 'Info', 'Debug'")

func main() {
	l, err := log.ParseLevel(*level)
	if err != nil {
		panic("Incorrect logging level")
	}
	log.SetLevel(l)

	logger := log.WithFields(log.Fields{
		"app": "rtdGO",
	})
	logger.Info("Application Initialization Begin...")

	flag.Parse()

	al := database.NewAccessLayer()
	if err = al.Open(); err != nil {
		panic(fmt.Sprintf("Unable to create and open database: %s", err))
	}
	defer func() {
		if err := al.Close(); err != nil {
			logger.Error(fmt.Sprintf("Error shutting down database connection: %s", err))
		}
	}()

	if *parse {
		logger.Warning("Parsing data into DB.")
		logger.Warning("---------- WARNING ----------")
		logger.Warning("Be sure the DB is new and empty, else this will error or create dup data")
		logger.Warning("---------- WARNING ----------")
		parser.ParseData(al, *sourceDir)
	}

	logger.Warning("Application Init complete. Running...")
	rh := web.NewRestHandler(al, 3000, logger)
	rh.Init()
}
