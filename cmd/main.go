package main

import (
	"eventhandler/internal/app"
	"eventhandler/util"
	"log"
)

func main() {
	cfg, err := util.LoadConfig("./")
	if err != nil {
		log.Fatal(err)
	}
	app.Run(cfg)
}
