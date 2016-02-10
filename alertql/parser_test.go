package alertql

import (
	"strings"
	"testing"
)

//TestParseAlert tests alert parsing
func TestParseAlert(t *testing.T) {
	p := NewParser(nil, nil)
	result, _ := p.ParseAlertStatement("ALERT testAlert IF testKey > 10 TEXT \"HELLO WORLD\"")

	if result.Name != "testAlert" {
		t.Errorf("Expected key to be testAlert was %s", result.Name)
	}

	if result.MetricKey != "testKey" {
		t.Errorf("Expected key to be testKey was %s", result.MetricKey)
	}

	if !result.Condition(11) {
		t.Error("Expected alert to be triggered")
	}

	if result.Condition(1) {
		t.Error("Did not expect alert to be triggered")
	}

	if !strings.EqualFold(result.Text, "HELLO WORLD") {
		t.Errorf("Expected text to contain \"HELLO WORLD\" instead found \"%s\"", result.Text)
	}
}

//TestParseSchedule tests the parsing of Schedule statements
func TestParseSchedule(t *testing.T) {
	p := NewParser(nil, nil)
	key, dbname, query, err := p.ParseScheduleStatement("SCHEDULE alert1 INFLUXDB \"SELECT last(value) from myseries\" ON public")

	if err != nil {
		t.Error(err)
	}

	if key != "alert1" {
		t.Errorf("Expected key to be alert1 was %s", key)
	}

	if dbname != "public" {
		t.Errorf("Expected key to be public was %s", dbname)
	}

	if query != "SELECT last(value) from myseries" {
		t.Errorf("Expected query to be \"SELECT last(value) from myseries\" was \"%s\"", query)
	}

}
