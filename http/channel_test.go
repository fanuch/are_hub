package http

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/blacksfk/are_server"
	"github.com/blacksfk/are_server/mock"
	uf "github.com/blacksfk/microframework"
	"github.com/julienschmidt/httprouter"
)

// Test data.
var channels []are_server.Channel = []are_server.Channel{
	{Name: "Audi Sport Team WRT", Password: "abc123"},
	{Name: "Mercedes-AMG Team Black Falcon", Password: "lol123"},
}

// Does it call repo.All?
// Is the Content-Type header set to application/json?
// Does it return all channels?
func TestChannelIndex(t *testing.T) {
	// mock All function
	fn := func(_ context.Context) ([]are_server.Channel, error) {
		return channels, nil
	}

	// create the mock repo and controller
	repo := &mock.ChannelRepo{AllFunc: fn}
	controller := NewChannel(repo)

	// create a mock request
	req, e := http.NewRequest(http.MethodGet, "/channel", nil)

	if e != nil {
		t.Fatal(e)
	}

	// create a response recorder and run the controller method
	w := httptest.NewRecorder()
	e = controller.Index(w, req)

	if e != nil {
		t.Fatal(e)
	}

	// get the response
	res := w.Result()

	// check if the repo was hit
	if !repo.AllCalled {
		t.Error("Did not call repo.All")
	}

	// ensure the content type is application/json
	checkCT(res, t)

	// extract the body and confirm all data was returned
	defer res.Body.Close()
	body, e := io.ReadAll(res.Body)

	if e != nil {
		t.Fatal(e)
	}

	var received []are_server.Channel
	e = json.Unmarshal(body, &received)

	if e != nil {
		t.Fatal(e)
	}

	lr := len(received)
	lc := len(channels)

	// check that all channels were returned
	if lr != lc {
		t.Fatalf("Expected: %d channels. Actual: %d.", lc, lr)
	}

	// loop and ensure the data is correct
	for i := 0; i < lr; i++ {
		if received[i].Name != channels[i].Name {
			t.Fatalf("Expected: %s. Actual: %s.", channels[i].Name, received[i].Name)
		}

		if received[i].Password != channels[i].Password {
			t.Fatalf("Expected: %s. Actual: %s.", channels[i].Password, received[i].Password)
		}
	}
}

// Does it call repo.Store?
// Is the content type set to application/json?
// Does it return the new channel?
func TestChannelStore(t *testing.T) {
	// mock Insert function
	fn := func(_ context.Context, v are_server.Archetype) error {
		return nil
	}

	// create mock repo and controller
	repo := &mock.ChannelRepo{InsertFunc: fn}
	controller := NewChannel(repo)

	// create and embed a new channel
	msport := are_server.Channel{Name: "Bentley Team M-Sport", Password: "abc123"}

	// create a mock request
	req, e := http.NewRequest(http.MethodPost, "/channel", nil)

	if e != nil {
		t.Fatal(e)
	}

	// update the request's context with the channel
	req = req.WithContext(msport.ToCtx(req.Context()))

	// create a response recorder and run the controller method
	w := httptest.NewRecorder()
	e = controller.Store(w, req)

	if e != nil {
		t.Fatal(e)
	}

	// check if the repo was hit
	if !repo.InsertCalled {
		t.Error("Did not call repo.Insert")
	}

	// get the response
	res := w.Result()

	// ensure the content type is application/json
	checkCT(res, t)

	// extract the returned channel
	defer res.Body.Close()
	resBody, e := io.ReadAll(res.Body)

	if e != nil {
		t.Fatal(e)
	}

	// unmarshal the response body
	received := are_server.Channel{}
	e = json.Unmarshal(resBody, &received)

	if e != nil {
		t.Fatal(e)
	}

	// compare the sent and received channels
	if msport.Name != received.Name || msport.Password != received.Password {
		t.Fatalf("Expected: %+v. Actual: %+v", msport, received)
	}
}

// Does it call repo.FindID?
// Is the content type set to application/json?
// Does it return the correct channel?
// Does it return a 404 Not Found error for an invalid ID?
func TestChannelShow(t *testing.T) {
	// expecting this channel
	wrt := channels[0]

	// create the mock repo and controller
	repo := &mock.ChannelRepo{FindIDFunc: findChannelID}
	controller := NewChannel(repo)

	// create a mock request
	p := httprouter.Param{Key: "id", Value: "1"}
	req, e := http.NewRequest(http.MethodGet, "/channel/"+p.Value, nil)

	if e != nil {
		t.Fatal(e)
	}

	// embed the channel ID in the request's
	// context (necessary for controller.Show to function)
	uf.EmbedParams(req, httprouter.Param{Key: "id", Value: "1"})

	// create a response recorder and call the show method
	w := httptest.NewRecorder()
	e = controller.Show(w, req)

	if e != nil {
		t.Fatal(e)
	}

	// check the repo was hit
	if !repo.FindIDCalled {
		t.Error("Did not call repo.FindID")
	}

	res := w.Result()

	// ensure the content type is application/json
	checkCT(res, t)

	// read and unmarshal the body
	defer res.Body.Close()
	body, e := io.ReadAll(res.Body)

	if e != nil {
		t.Fatal(e)
	}

	received := are_server.Channel{}
	e = json.Unmarshal(body, &received)

	if e != nil {
		t.Fatal(e)
	}

	// compare the expected and received channels
	if received.Name != wrt.Name || received.Password != wrt.Password {
		t.Fatalf("Expected: %+v. Actual: %+v.", wrt, received)
	}

	// check show returns 404 for an invalid ID
	p = httprouter.Param{Key: "id", Value: "-1"}
	test404(t, http.MethodGet, "/channel/"+p.Value, nil, controller.Show, p)
}

// Does it call repo.UpdateID?
// Is the content type set to application/json?
// Does it return the updated channel?
func TestChannelUpdate(t *testing.T) {
	// mock UpdateID function
	fn := func(_ context.Context, str string, ch *are_server.Channel) error {
		_, e := findChannelID(nil, str)

		// the update itself has no bearing on the test so simply return
		// the error (if there was one)
		return e
	}

	// create mock repo and controller
	repo := &mock.ChannelRepo{UpdateIDFunc: fn}
	controller := NewChannel(repo)

	// mock channel
	wrt := are_server.Channel{Name: "Belgian Audi Club WRT", Password: "abc123"}

	// create mock request
	p := httprouter.Param{Key: "id", Value: "1"}
	req, e := http.NewRequest(http.MethodPut, "/channel/"+p.Value, nil)

	if e != nil {
		t.Fatal(e)
	}

	// embed the updated channel in the request's context
	req = req.WithContext(wrt.ToCtx(req.Context()))

	// embed parameters in the request's context
	uf.EmbedParams(req, p)

	// create a response recorder run the update method
	w := httptest.NewRecorder()
	e = controller.Update(w, req)

	if e != nil {
		t.Fatal(e)
	}

	res := w.Result()

	// check if repo was hit
	if !repo.UpdateIDCalled {
		t.Error("Did not call repo.UpdateID")
	}

	// ensure the content type is applicaton/json
	checkCT(res, t)

	// read and unmarshal the body
	defer res.Body.Close()
	resBody, e := io.ReadAll(res.Body)

	if e != nil {
		t.Fatal(e)
	}

	received := are_server.Channel{}
	e = json.Unmarshal(resBody, &received)

	if e != nil {
		t.Fatal(e)
	}

	// compare the sent and received channels
	if wrt.Name != received.Name || wrt.Password != received.Password {
		t.Fatalf("Expected: %+v. Actual: %+v", wrt, received)
	}

	// check if Update returns a 404 error on an invalid ID
	p = httprouter.Param{Key: "id", Value: "-1"}
	req, e = http.NewRequest(http.MethodPut, "/channel/"+p.Value, nil)

	if e != nil {
		t.Fatal(e)
	}

	// embed non-existant channel into the request's context
	gpx := are_server.Channel{Name: "Grand Prix Extreme", Password: "porsche"}
	req = req.WithContext(gpx.ToCtx(req.Context()))

	// embed parameters
	uf.EmbedParams(req, p)

	// create a new response recorder and call the update method
	w = httptest.NewRecorder()
	e = controller.Update(w, req)

	if e == nil {
		t.Fatal("Expected: 404 Not found error. Actual: nil")
	}

	he, ok := e.(uf.HttpError)

	if !ok {
		t.Fatalf("Expected: 404 Not Found error. Actual: %+v", e)
	}

	if he.Code != http.StatusNotFound {
		t.Fatalf("Expected: %d. Actual: %d", http.StatusNotFound, he.Code)
	}
}

// Does it call repo.DeleteID with the correct ID?
// Is the content type set to application/json?
// Does it return the deleted channel?
// Does it return a 404 Not Found error for an invalid ID?
func TestChannelDelete(t *testing.T) {
	// delete this channel
	wrt := channels[0]

	// create the mock repo and controller.
	// the deletion itself has no bearing on the test
	// so just use the findID function which has the the same signature
	// and performs the operation we need
	repo := &mock.ChannelRepo{DeleteIDFunc: findChannelID}
	controller := NewChannel(repo)

	// create a mock request
	p := httprouter.Param{Key: "id", Value: "1"}
	req, e := http.NewRequest(http.MethodDelete, "/channel/"+p.Value, nil)

	if e != nil {
		t.Fatal(e)
	}

	// embed params necessary for controller function
	uf.EmbedParams(req, p)

	// create a response recorder and call the delete method
	w := httptest.NewRecorder()
	e = controller.Delete(w, req)

	if e != nil {
		t.Fatal(e)
	}

	res := w.Result()

	// check if the repo was hit
	if !repo.DeleteIDCalled {
		t.Error("Did not call repo.DeleteID")
	}

	// ensure the content type is application/json
	checkCT(res, t)

	// extract the body and check the correct channel was returned
	defer res.Body.Close()
	body, e := io.ReadAll(res.Body)

	if e != nil {
		t.Fatal(e)
	}

	received := &are_server.Channel{}
	e = json.Unmarshal(body, received)

	if e != nil {
		t.Fatal(e)
	}

	if received.Name != wrt.Name || received.Password != wrt.Password {
		t.Fatalf("Expected: %v. Actual: %v.", wrt, received)
	}

	// check delete returns 404 for an invalid ID
	p = httprouter.Param{Key: "id", Value: "-1"}
	test404(t, http.MethodDelete, "/channel/"+p.Value, nil, controller.Delete, p)
}

func findChannelID(_ context.Context, str string) (*are_server.Channel, error) {
	id64, e := strconv.ParseInt(str, 10, 64)

	if e != nil {
		return nil, e
	}

	id := int(id64)
	idx := id - 1

	if idx < 0 || idx >= len(channels) {
		return nil, are_server.NewNoObjectsFound("channels", "id == "+str)
	}

	return &channels[id-1], nil
}

// Helper function to check the content type was set to "application/json"
func checkCT(res *http.Response, t *testing.T) {
	if ct := res.Header.Get("Content-Type"); ct != "application/json" {
		t.Fatalf("Incorrect content type. Expected: application/json. Actual: %s", ct)
	}
}
