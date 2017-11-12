package main

import (
	"flag"
	"fmt"

	"github.com/kimmyfek/next_rtd/database"
	"github.com/kimmyfek/next_rtd/parser"
	"github.com/kimmyfek/next_rtd/web"
)

var parse = flag.Bool("parse", false, "If provided will parse data and add to DB. Default: Disabled")
var dbPath = flag.String("dbPath", "/opt/next-rtd/next-rtd.db", "Path to sqlite DB file.")

// TODO -- Flag for parse file path

func main() {
	fmt.Println("Application Initialization Begin...")

	flag.Parse()
	fmt.Println(fmt.Sprintf("Accessing DB @ %s", *dbPath))
	al := database.NewAccessLayer(*dbPath)
	err := al.Open()
	if err != nil {
		panic(fmt.Sprintf("Unable to create and open database: %s", err))
	}

	if *parse {
		fmt.Println("Parsing data into DB.")
		fmt.Println("---------- WARNING ----------")
		fmt.Println("Be sure the DB is new and empty, else this will error or create dup data")
		fmt.Println("---------- WARNING ----------")
		parser.ParseData(al)
	}

	fmt.Println("Application Init complete. Running...")
	rh := web.NewRestHandler(al, 3000)
	rh.Init()
}
