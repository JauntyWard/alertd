package alertql

import (
	"fmt"

	"github.com/jauntyward/alertd/engine"
)

//NewCondition returns a function for triggering an alert from a query
func NewCondition(operator string, threshold float64) (engine.AlertCondition, error) {
	var alertRule engine.AlertCondition

	if operator == "=" {
		alertRule = func(value float64) bool {
			if value == threshold {
				return true
			}
			return false
		}
	} else if operator == ">" {
		alertRule = func(value float64) bool {
			if value > threshold {
				return true
			}
			return false
		}
	} else if operator == "<" {
		alertRule = func(value float64) bool {
			if value < threshold {
				return true
			}
			return false
		}
	} else {
		return *new(engine.AlertCondition), fmt.Errorf("Invalid operator")
	}

	return alertRule, nil
}
