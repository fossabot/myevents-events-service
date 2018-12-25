package main

import (
	"flag"
	"fmt"
	"github.com/danielpacak/go-rest-api-seed/configuration"
	"github.com/danielpacak/go-rest-api-seed/dblayer"
	"github.com/danielpacak/go-rest-api-seed/rest"
	"log"
)

func main() {
	confPath := flag.String("conf", `.\configuration\config.json`, "flag to set the path to the configuration json file")
	flag.Parse()
	config, _ := configuration.ExtractConfiguration(*confPath)
	fmt.Println("Connecting to database")
	dbhandler, _ := dblayer.NewPersistenceLayer(config.Databasetype, config.DBConnection)

	log.Fatal(rest.ServeAPI(config.RestfulEndpoint, dbhandler))
}
