package main

import (
	"flag"

	"github.com/jauntyward/alertd/alertql"
	"github.com/jauntyward/alertd/api"
	"github.com/jauntyward/alertd/config"
	"github.com/jauntyward/alertd/db"
	"github.com/jauntyward/alertd/engine"
)

func main() {
	configFilePath := flag.String("config", "", "configuration file")
	flag.Parse()

	configFile, err := config.ReadConfig(*configFilePath)

	if err != nil {
		panic("Unable to read config")
	}

	topLevelConfig := config.ParseConfig(configFile)

	engineInstance := engine.NewAlertEngine(&topLevelConfig.AlertEngineConfig)
	dbscheduler := db.NewScheduler(&topLevelConfig.InfluxDBConfig, *engineInstance)
	parser := alertql.NewParser(engineInstance, dbscheduler)
	apiInstance := api.NewAPI(engineInstance, parser)

	go dbscheduler.Schedule()
	apiInstance.ServeAPI()

}
