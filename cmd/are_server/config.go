package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/blacksfk/are_server/mongodb"
)

// Application configuration parameters.
type config struct {
	// mongodb connection parameters
	MongoDB *mongodb.Params

	// address for the server to listen on. Eg. ":6060".
	Address string

	// allow requests originating from this domain. Eg. "example.com", "*".
	AllowOrigin string
}

// Unmarshal file as JSON into a config struct. This function is only intended to be
// called from the main function and therefore dies if it encounters an error
// reading file or processing file's bytes as JSON.
func load(file string) *config {
	// read the entire file
	bytes, e := os.ReadFile(file)

	if e != nil {
		// reading failed so die
		log.Fatal("Error loading %s:", e)
	}

	// unmarshal the bytes
	conf := &config{}
	e = json.Unmarshal(bytes, conf)

	if e != nil {
		// unmarshalling failed so die
		log.Fatalf("Error processing %s as JSON:", e)
	}

	return conf
}
