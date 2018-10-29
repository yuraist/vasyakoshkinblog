package main

import (
	"fmt"
	"github.com/yuraist/vasyakoshkin/blogservice/rest"
	"github.com/yuraist/vasyakoshkin/dblayer"
)

func main() {
	fmt.Println("Connecting to database...")
	dbHandler, err := mongolayer.NewMongoLayer()

	if err != nil {
		fmt.Printf("An %v error was occured at the programm beginning", err)
	}

	rest.ServeRestAPI("localhost:8000", *dbHandler)
}