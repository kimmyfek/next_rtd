package main

import (
	"fmt"

	"github.com/nursejason/next-rtd/database"
	"github.com/nursejason/next-rtd/web"
)

func main() {
	fmt.Println("Starting application")

	// TODO Flags -- args
	al := database.NewAccessLayer("/opt/next-rtd/next-rtd.db")
	err := al.Open()
	if err != nil {
		panic(fmt.Sprintf("Unable to create and open database: %s", err))
	}

	rh := web.NewRestHandler(al, 3000)
	rh.Init()
}
