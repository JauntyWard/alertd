package db

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/influxdb/influxdb/client"
	"github.com/jauntyward/alertd/engine"
)

type (
	//Scheduler represents an instance of the query scheduler
	Scheduler struct {
		InfluxDBConfig *InfluxDBConfig
		Engine         engine.AlertEngine
		Stop           chan struct{}
		Queries        map[ScheduleKey]string
	}

	//InfluxDBConfig provides the config required to pull metrics from InfluxDB
	InfluxDBConfig struct {
		InfluxDBHost     string
		InfluxDBPort     int
		InfluxDBUsername string
		InfluxDBPassword string
		InfluxDBDB       string
	}

	//ScheduleKey key indexes a map of scheduled DB queries
	ScheduleKey struct {
		MetricKey string
		Database  string
	}
)

//NewScheduler returns a new instance of DBScheduler
func NewScheduler(c *InfluxDBConfig, engine engine.AlertEngine) *Scheduler {
	s := &Scheduler{InfluxDBConfig: c, Engine: engine}
	s.Queries = make(map[ScheduleKey]string)
	s.Stop = make(chan struct{})
	return s
}

//AddQuery adds a query to the scheduler
func (db *Scheduler) AddQuery(metricKey string, database string, query string) {
	key := &ScheduleKey{MetricKey: metricKey, Database: database}
	db.Queries[*key] = query
}

//Schedule periodically executes predefined InfluxDB queries
func (db *Scheduler) Schedule() {
	ticker := time.NewTicker(5 * time.Second)

	for {
		select {
		case <-ticker.C:
			for key, query := range db.Queries {
				go db.scheduledCheck(key, query)
			}
		case <-db.Stop:
			ticker.Stop()
			return
		}
	}
}

func (db *Scheduler) scheduledCheck(key ScheduleKey, query string) {
	value, err := db.ExecuteQuery(query, key.Database)
	if err == nil {
		db.Engine.Check(key.MetricKey, value)
	}
}

//ExecuteQuery executes an InfluxDB query and returns the resultant value
func (db *Scheduler) ExecuteQuery(query string, database string) (float64, error) {
	host, err := url.Parse(fmt.Sprintf("http://%s:%d", db.InfluxDBConfig.InfluxDBHost, db.InfluxDBConfig.InfluxDBPort))
	if err != nil {
		log.Fatal(err)
	}
	con, err := client.NewClient(client.Config{
		URL: *host,
		//Username: db.InfluxDBConfig.InfluxDBUsername,
		//Password: db.InfluxDBConfig.InfluxDBPassword,
	})

	if err != nil {
		return 0, err
	}

	q := client.Query{
		Command:  query,
		Database: database,
	}

	response, err := con.Query(q)

	if err == nil {
		var jsonValue json.Number
		//This somewhat unpleasant looking line goes through several arrays nested structs
		//to get to the actual value.
		responseValue := response.Results[0].Series[0].Values[0][1]
		//the value is encoded as a JSON number as it comes from the web API
		jsonValue = responseValue.(json.Number)
		//parse a float from the json value
		value, _ := jsonValue.Float64()

		if response.Err == nil {
			return value, nil
		}
		return 0, fmt.Errorf("Unable to parse value from InfluxDB query, ensure that query returns a single value and that the series contains data")
	}
	return 0, err
}
