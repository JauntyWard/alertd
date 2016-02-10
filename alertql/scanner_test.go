package alertql

import (
	"fmt"
	"testing"
)

//TestScanLength tests if scan returns the correct number of tokens
func TestScanLength(t *testing.T) {
	scanner := NewScanner("ALERT foo IF bar > 10 TEXT \"HELLO WOLRD\"")
	result := scanner.scan()

	if len(result) != 8 {
		fmt.Println(result)
		t.Error("Was expecting 8 tokens, got", len(result))
	}
}

//TestScanTokens tests if scan returns the correct type of tokens
func TestScanTokens(t *testing.T) {
	scanner := NewScanner("IF ALERT \"string\" INFLUXDB < > ==")
	result := scanner.scan()

	if result[0].tokenType != IF {
		t.Errorf("Was expecting IF token found %s", result[0].literal)
	}

	if result[1].tokenType != ALERT {
		t.Errorf("Was expecting ALERT token found %s", result[1].literal)
	}

	if result[2].tokenType != STRING {
		t.Errorf("Was expecting STRING token found %s", result[2].literal)
	}

	if result[3].tokenType != INFLUXDB {
		t.Errorf("Was expecting INFLUXDB token found %s", result[3].literal)
	}

	if result[4].tokenType != OP {
		t.Errorf("Was expecting OP token found %s", result[4].literal)
	}

	if result[5].tokenType != OP {
		t.Errorf("Was expecting OP token found %s", result[5].literal)
	}
	if result[6].tokenType != OP {
		t.Errorf("Was expecting OP token found %s", result[6].literal)
	}

}
