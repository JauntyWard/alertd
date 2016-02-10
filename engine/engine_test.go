package engine

import "testing"

func TestBasicAlertRules(t *testing.T) {
	engineInstance := NewAlertEngine(&Config{})

	testCondition := func(value float64) bool {
		if value > 10 {
			return true
		}
		return false
	}

	testRule := Rule{Name: "rule1",
		MetricKey: "metric1",
		Condition: testCondition,
		Text:      "this is a test event",
	}

	engineInstance.AddRule(testRule)

	result1, _ := engineInstance.Check("metric1", 11)
	if !result1 {
		t.Error("Alert rule should have been triggered")
	}

	result2, _ := engineInstance.Check("metric1", 3)
	if result2 {
		t.Error("Alert rule should not have been triggered")
	}

}
