package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
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
