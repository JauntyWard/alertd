package config

import (
	"bytes"
	"io/ioutil"
	"log"

	"github.com/jauntyward/alertd/db"
	"github.com/jauntyward/alertd/engine"

	"github.com/spf13/viper"
)

type (
	//CCPAlertConfig is a struct representing the configuration for CCP Alert
	CCPAlertConfig struct {
		//InfluxDBConfig describes the configuration needed to communicate with InflxuDB
		InfluxDBConfig db.InfluxDBConfig
		//AlertEngineConfig represents the config for AlertEngine
		AlertEngineConfig engine.Config
	}
)

//ReadConfig takes a file path as a string and returns a string representing
//the contents of that file
func ReadConfig(configFile string) ([]byte, error) {
	//viper accepts config file without extension, so remove extension
	if configFile == "" {
		panic("No config file provided")
	}

	f, err := ioutil.ReadFile(configFile)

	if err != nil {
		log.Fatal(err)
	}

	return f, err
}

//ParseConfig parses a YAML config  file
func ParseConfig(rawConfig []byte) CCPAlertConfig {
	parsedConfig := new(CCPAlertConfig)

	viper.SetConfigType("yaml")
	viper.ReadConfig(bytes.NewBuffer(rawConfig))

	parsedConfig.AlertEngineConfig.PagerDutyAPIKey = viper.GetString("PagerDutyAPIKey")

	parsedConfig.AlertEngineConfig.EmailServer = viper.GetString("email.server")
	parsedConfig.AlertEngineConfig.EmailUsername = viper.GetString("email.username")
	parsedConfig.AlertEngineConfig.EmailPassword = viper.GetString("email.password")
	parsedConfig.AlertEngineConfig.EmailPort = viper.GetInt("email.port")
	parsedConfig.AlertEngineConfig.EmailRecipient = viper.GetString("email.recipient")

	parsedConfig.InfluxDBConfig = *new(db.InfluxDBConfig)
	parsedConfig.InfluxDBConfig.InfluxDBHost = viper.GetString("influx.host")
	parsedConfig.InfluxDBConfig.InfluxDBPort = viper.GetInt("influx.port")
	parsedConfig.InfluxDBConfig.InfluxDBUsername = viper.GetString("influx.username")
	parsedConfig.InfluxDBConfig.InfluxDBPassword = viper.GetString("influx.password")
	parsedConfig.InfluxDBConfig.InfluxDBDB = viper.GetString("influx.password")

	if (len(parsedConfig.InfluxDBConfig.InfluxDBHost)) == 0 {
		panic("InfluxDB host undefined")
	}

	if parsedConfig.InfluxDBConfig.InfluxDBPort == 0 {
		panic("InfluxDB port undefined")
	}

	if (len(parsedConfig.InfluxDBConfig.InfluxDBUsername)) == 0 {
		panic("InfluxDB username undefined")
	}

	if (len(parsedConfig.InfluxDBConfig.InfluxDBPassword)) == 0 {
		panic("InfluxDB password undefined")
	}

	if (len(parsedConfig.InfluxDBConfig.InfluxDBDB)) == 0 {
		panic("InfluxDB db undefined")
	}

	return *parsedConfig
}
