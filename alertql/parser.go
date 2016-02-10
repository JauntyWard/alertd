package alertql

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jauntyward/alertd/db"
	"github.com/jauntyward/alertd/engine"
)

type (
	//Parser represents an instance of the CCPAlertQL parser
	Parser struct {
		Scheduler    *db.Scheduler
		Engine       *engine.AlertEngine
		RuleFilePath string
	}

	//Result represents any outcome of parsing which should be returned to the user
	Result struct {
		OK         bool
		ResultList []string
	}
)

//NewParser returns a new instance of the CCPAlertQL parser
func NewParser(engine *engine.AlertEngine, scheduler *db.Scheduler) *Parser {
	return &Parser{Scheduler: scheduler, Engine: engine}
}

//Parse identifies the query and calls the apppropriate parser function
func (p *Parser) Parse(query string) (Result, error) {
	var err error
	var result *Result

	if len(query) == 0 {
		err = fmt.Errorf("Unable to parse query")
		result = &Result{OK: false}
		return *result, err
	}

	switch strings.Fields(query)[0] {
	case "ALERT":
		var rule engine.Rule
		rule, err = p.ParseAlertStatement(query)
		p.Engine.AddRule(rule)
	case "SCHEDULE":
		var key, dbname, influxQuery string
		key, dbname, influxQuery, err = p.ParseScheduleStatement(query)

		//Check the the encapsulted InfluxDB query is valid
		if _, err = p.Scheduler.ExecuteQuery(influxQuery, dbname); err != nil {
			p.Scheduler.AddQuery(key, dbname, query)
			p.WriteQuery(query)
		}
	case "SHOW":
		var showRequest TokenType
		showRequest, err = p.ParseShowStatement(query)

		switch showRequest {
		case ALERTS:

		}
	}

	if err != nil {
		result = &Result{OK: false}
	} else {
		//make query persistent
		p.WriteQuery(query)
		result = &Result{OK: true}
	}

	return *result, nil
}

//WriteQuery writes a validated query to a file in order to make it persistent
func (p *Parser) WriteQuery(query string) {
	f, err := os.OpenFile(p.RuleFilePath, os.O_APPEND, 0666)
	if err == nil {
		f.WriteString(query)
		f.Close()
	}
}

//ParseScheduleStatement parses a schedule query and schedules the contained InfluxDB query
//A schedule statement takes the form of:
//SCHEDULE <ID> INFLUXDB <influxdb query> ON <Db name?
//To give examples:
//SCHEDULE alert1 INFLUXDB "SELECT last(value) from myseries" ON public
func (p *Parser) ParseScheduleStatement(scheduleStatment string) (string, string, string, error) {
	scanner := NewScanner(scheduleStatment)
	tokens := scanner.scan()

	if tokens[0].tokenType != SCHEDULE {
		err := fmt.Errorf("found %q, expected SCHEDULE", tokens[0].literal)
		return "", "", "", err
	}

	if tokens[1].tokenType != IDENTIFIER {
		err := fmt.Errorf("found %q, expected IDENTIFIER", tokens[0].literal)
		return "", "", "", err
	}
	key := tokens[1].literal

	if tokens[2].tokenType != INFLUXDB {
		err := fmt.Errorf("found %q, expected INFLUXDB", tokens[0].literal)
		return "", "", "", err
	}

	if tokens[3].tokenType != STRING {
		err := fmt.Errorf("found %q, expected INFLUXDB", tokens[0].literal)
		return "", "", "", err
	}

	query := tokens[3].literal

	if tokens[3].tokenType != STRING {
		err := fmt.Errorf("found %q, expected INFLUXDB", tokens[0].literal)
		return "", "", "", err
	}

	if tokens[4].tokenType != ON {
		err := fmt.Errorf("found %q, expected ON", tokens[0].literal)
		return "", "", "", err
	}

	if tokens[5].tokenType != IDENTIFIER {
		err := fmt.Errorf("found %q, expected db name", tokens[0].literal)
		return "", "", "", err
	}

	dbname := tokens[5].literal

	if len(tokens) > 6 {
		err := fmt.Errorf("trailing characters %q", tokens[8].literal)
		return "", "", "", err
	}

	return key, dbname, query, nil
}

//ParseAlertStatement takes a raw alert statement query and parses it to a Rule struct
//An alert statement stakes the form:
//ALERT <alert name> IF <metric name> <operator> <threshold value> TEXT <description of alert>
//To give examples:
//ALERT cpuOnFireAlert IF superImportantServer.cpuUsage > 100 TEXT "Critical production server is heavily loaded"
//ALERT noplayers IF tq.currentPlayers == 0 TEXT "something has gone badly wrong"
func (p *Parser) ParseAlertStatement(alertStatement string) (engine.Rule, error) {
	scanner := NewScanner(alertStatement)
	tokens := scanner.scan()
	newRule := new(engine.Rule)

	if tokens[0].tokenType != ALERT {
		err := fmt.Errorf("found %q, expected ALERT", tokens[0].literal)
		return engine.Rule{}, err
	}

	if tokens[1].tokenType == IDENTIFIER {
		newRule.Name = tokens[1].literal
	} else {
		err := fmt.Errorf("found %q, expected identifier", tokens[1].literal)
		return engine.Rule{}, err
	}

	if tokens[2].tokenType != IF {
		err := fmt.Errorf("found %q, expected IF", tokens[2].literal)
		return engine.Rule{}, err
	}

	if tokens[3].tokenType == IDENTIFIER {
		newRule.MetricKey = tokens[3].literal
	} else {
		err := fmt.Errorf("found %q, expected identifier", tokens[3].literal)
		return engine.Rule{}, err
	}

	if tokens[4].tokenType != OP {
		err := fmt.Errorf("found %q, expected <,> or ==", tokens[4].literal)
		return engine.Rule{}, err
	}

	if tokens[5].tokenType != VALUE {
		err := fmt.Errorf("found %q, expected value", tokens[5].literal)
		return engine.Rule{}, err
	}

	threshold, err := strconv.ParseFloat(tokens[5].literal, 64)

	if err != nil {
		return engine.Rule{}, err
	}

	condition, err := NewCondition(tokens[4].literal, threshold)

	if err == nil {
		newRule.Condition = condition
	} else {
		return engine.Rule{}, err
	}

	if tokens[6].tokenType != TEXT {
		err := fmt.Errorf("found %q, expected TEXT", tokens[6].literal)
		return engine.Rule{}, err
	}

	if tokens[7].tokenType != STRING {
		err := fmt.Errorf("found %q, expected string", tokens[7].literal)
		return engine.Rule{}, err
	}

	newRule.Text = tokens[7].literal

	if len(tokens) > 8 {
		err := fmt.Errorf("trailing characters %q", tokens[8].literal)
		return engine.Rule{}, err
	}

	return *newRule, nil
}

//ParseShowStatement parses a show statement and
func (p *Parser) ParseShowStatement(showStatement string) (TokenType, error) {
	scanner := NewScanner(showStatement)
	tokens := scanner.scan()

	if tokens[0].tokenType != SHOW {
		err := fmt.Errorf("expected SHOW found %q", tokens[0].literal)
		return ILLEGAL, err
	}

	if tokens[1].tokenType == ALERTS {
		return ALERTS, nil
	} else if tokens[1].tokenType == SCHEDULED {
		return SCHEDULE, nil
	} else {
		err := fmt.Errorf("expected ALERTS or SCHEDULED found %q", tokens[1].literal)
		return ILLEGAL, err
	}

}
