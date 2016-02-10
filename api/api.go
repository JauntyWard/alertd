package api

import (
	"encoding/json"
	"net/http"

	"github.com/jauntyward/alertd/alertql"
	"github.com/jauntyward/alertd/engine"
)

type (
	//AlertdAPI represents an instance of the API
	AlertdAPI struct {
		Engine *engine.AlertEngine
		Parser *alertql.Parser
	}

	ruleRequest struct {
		RawAlertStatement string
	}

	checkRequest struct {
		Key   string
		Value float64
	}
)

//NewAPI returns a new isntance of AlertdAPI
func NewAPI(e *engine.AlertEngine, p *alertql.Parser) *AlertdAPI {
	return &AlertdAPI{Engine: e, Parser: p}
}

func (api *AlertdAPI) query(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var query ruleRequest
	err := decoder.Decode(&query)

	if err != nil {
		http.Error(w, "invalid rule", 500)
		return
	}

	_, err = api.Parser.Parse(query.RawAlertStatement)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

//ServeAPI serves the ccpalert API on port 8080
func (api *AlertdAPI) ServeAPI() {
	server := http.NewServeMux()
	server.HandleFunc("/query", api.query)
	http.ListenAndServe(":8080", server)
}
