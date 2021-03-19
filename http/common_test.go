package http

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	uf "github.com/blacksfk/microframework"
	"github.com/julienschmidt/httprouter"
)

func test404(t *testing.T, method, url string, body io.Reader, h uf.Handler, params ...httprouter.Param) {
	r, e := http.NewRequest(method, url, body)

	if e != nil {
		t.Fatal(e)
	}

	// embed parameters
	uf.EmbedParams(r, params...)

	// create a new response recorder and call the handler
	w := httptest.NewRecorder()
	e = h(w, r)

	if e == nil {
		t.Fatal("Expected: 404 Not Found error. Actual: nil.")
	}

	he, ok := e.(uf.HttpError)

	if !ok {
		t.Fatalf("Expected: 404 Not Found error. Actual: %+v", e)
	}

	if he.Code != http.StatusNotFound {
		t.Fatalf("Expected: %d. Actual: %d", http.StatusNotFound, he.Code)
	}
}
