package validate

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/blacksfk/are_server"
	uf "github.com/blacksfk/microframework"
)

// Does it store valid channels in the request's context?
func TestStoreValid(t *testing.T) {
	// create mock data
	embedded := are_server.Channel{
		Name:     "R-Motorsport AMR",
		Password: "amr",
	}

	// marshal the data
	body, e := json.Marshal(embedded)

	if e != nil {
		t.Fatal(e)
	}

	// create a mock request
	r, e := http.NewRequest(http.MethodPost, "/channel", bytes.NewReader(body))

	if e != nil {
		t.Fatal(e)
	}

	// set content type
	r.Header.Set("Content-Type", "application/json")

	// create a validator and run the validation method
	v := NewChannel()
	e = v.Store(r)

	if e != nil {
		t.Fatal(e)
	}

	// confirm the channel was stored in the request
	extracted, e := are_server.ChannelFromCtx(r.Context())

	if e != nil {
		t.Fatal(e)
	}

	if extracted.Name != embedded.Name || extracted.Password != embedded.Password {
		t.Fatalf("Expected: %+v. Actual: %+v.", embedded, extracted)
	}
}

// Does it reject requests with incorrect content-types?
func TestStoreBadCT(t *testing.T) {
	// create mock data
	embedded := are_server.Channel{
		Name:     "R-Motorsport AMR",
		Password: "amr",
	}

	// marshal the data
	body, e := json.Marshal(embedded)

	if e != nil {
		t.Fatal(e)
	}

	// create a mock request
	r, e := http.NewRequest(http.MethodPost, "/channel", bytes.NewReader(body))

	if e != nil {
		t.Fatal(e)
	}

	// create a validator and run the validation method with no content-type set
	v := NewChannel()
	e = v.Store(r)

	// an error should be returned
	if e == nil {
		t.Fatal("Expected: 400 Bad Request error. Actual: nil.")
	}

	he, ok := e.(uf.HttpError)

	if !ok {
		t.Fatalf("Expected: 400 Bad Request error. Actual: %s.", e)
	}

	if he.Code != http.StatusBadRequest {
		t.Fatalf("Expected: %d. Actual: %d.", http.StatusBadRequest, he.Code)
	}
}

// Does it reject invalid channels with a 400 Bad Request error?
func TestStoreInvalid(t *testing.T) {
	// create mock data
	embedded := are_server.Channel{
		Name:     "",
		Password: "",
	}

	// marshal the data
	body, e := json.Marshal(embedded)

	if e != nil {
		t.Fatal(e)
	}

	// create a mock request
	r, e := http.NewRequest(http.MethodPost, "/channel", bytes.NewReader(body))

	if e != nil {
		t.Fatal(e)
	}

	r.Header.Set("Content-Type", "application/json")

	// create a validator and run the validation method
	v := NewChannel()
	e = v.Store(r)

	// an error should be returned
	if e == nil {
		t.Fatal("Expected: 400 Bad Request error. Actual: nil.")
	}

	he, ok := e.(uf.HttpError)

	if !ok {
		t.Fatalf("Expected: 400 Bad Request error. Actual: %s.", e)
	}

	if he.Code != http.StatusBadRequest {
		t.Fatalf("Expected: %d. Actual: %d.", http.StatusBadRequest, he.Code)
	}
}

// Does it reject malformed request bodies?
func TestStoreMalformed(t *testing.T) {
	// create garbage data
	body := []byte(`""\`)
	r, e := http.NewRequest(http.MethodPost, "/channel", bytes.NewReader(body))

	if e != nil {
		t.Fatal(e)
	}

	r.Header.Set("Content-Type", "application/json")

	// create a validator and run the validation method
	v := NewChannel()
	e = v.Store(r)

	// an error should be returned
	if e == nil {
		t.Fatal("Expected: 400 Bad Request error. Actual: nil.")
	}

	he, ok := e.(uf.HttpError)

	if !ok {
		t.Fatalf("Expected: 400 Bad Request error. Actual: %s.", e)
	}

	if he.Code != http.StatusBadRequest {
		t.Fatalf("Expected: %d. Actual: %d.", http.StatusBadRequest, he.Code)
	}
}
