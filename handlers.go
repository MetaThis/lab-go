package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/xeipuuv/gojsonschema"
)

// Run contains the data needed for a "run" on a given lab instrument.
type Run struct {
	InstrumentID int
	Samples      []Sample
}

type Sample struct {
	ID int `json:"id"`
}

type ValidationError struct {
	Errors []string `json:"errors"`
}

type RunSuccessResponse struct {
	RunID int `json:"runId"`
}

// Handlers is the receiver for methods that handle http requests.
// It will provide dependencies. Use NewHandlers() for construction.
type Handlers struct {
	DB     DB
	Schema Schema
}

// NewHandlers is the constructor for Handlers. Unit tests may provide
// an alternative constructor to "inject" alternative dependencies.
func NewHandlers(db DB, schema Schema) Handlers {
	h := Handlers{DB: db, Schema: schema}
	return h
}

// RespondAsJSON has the boilerplate for the most typical responses.
func RespondAsJSON(payload interface{}, statusCode int, w http.ResponseWriter) {
	json, err := json.Marshal(payload)
	if err != nil {
		log.Println("Error marshalling payload", err.Error())
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(json)
}

func (h Handlers) Post(w http.ResponseWriter, r *http.Request) {
	// Extract instrumentID from URL.
	vars := mux.Vars(r)
	unparsedInstrumentID := vars["instrument_id"]
	instrumentID, err := strconv.Atoi(unparsedInstrumentID)
	if err != nil {
		v := ValidationError{[]string{"Instrument ID in URL must be an integer."}}
		RespondAsJSON(v, 400, w)
		return
	}

	// Read body
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Validate
	jsonLoader := gojsonschema.NewBytesLoader(body)
	result, err := h.Schema.Samples.Validate(jsonLoader)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if !result.Valid() {
		errors := make([]string, len(result.Errors()))
		for i, desc := range result.Errors() {
			errors[i] = desc.Description()
		}
		v := ValidationError{errors}
		RespondAsJSON(v, 400, w)
		return
	}

	// Data is good, we can unmarshal.
	samples := make([]Sample, 0)
	err = json.Unmarshal(body, &samples)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Execute run (just saving to db for the scope of this exercise).
	run := Run{instrumentID, samples}
	runID, err := h.DB.SaveRun(run)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Respond with newly created run ID.
	success := RunSuccessResponse{runID}
	RespondAsJSON(success, 200, w)
}
