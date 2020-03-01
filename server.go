package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/xeipuuv/gojsonschema"
)

func main() {
	// In a real app, the connection information should be configuration
	// provided by the environment and pass into this constructor.
	db := NewDB("lab.db", true)
	schema := NewSchema()
	h := NewHandlers(db, schema)
	r := NewRouter(h)
	fmt.Println("Starting server")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// NewRouter is extracted to be reusable by tests. The Handlers object provides the dependencies.
func NewRouter(h Handlers) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/lab/instrument/{instrument_id}/samples", h.Post)
	// ... additional routes ...
	return r
}

// Schema holds loaded JSON schemas used for validation
type Schema struct {
	Samples *gojsonschema.Schema
}

func NewSchema() Schema {
	path, err := filepath.Abs("./schema/samples.json")
	if err != nil {
		log.Fatal(err)
	}
	samplesLoader := gojsonschema.NewReferenceLoader("file://" + path)
	samples, err := gojsonschema.NewSchema(samplesLoader)
	if err != nil {
		log.Fatal(err)
	}
	schema := Schema{Samples: samples}
	return schema
}
