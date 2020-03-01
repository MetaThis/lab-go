package main

import (
	"log"
	"path/filepath"

	"github.com/xeipuuv/gojsonschema"
)

// Schema holds loaded JSON schemas used for validation
type Schema struct {
	Samples *gojsonschema.Schema
}

// NewSchema initializes our schema(s) from the JSON file(s) that define them.
func NewSchema() Schema {
	path, err := filepath.Abs("./json-schemas/samples.json")
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
