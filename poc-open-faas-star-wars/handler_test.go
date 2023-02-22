package function

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func mockController(w http.ResponseWriter, r *http.Request) {
	log.Println("controller mock ok")

	w.WriteHeader(http.StatusOK)

}

func init() {
	routes["/starwar"] = mockController
}

func TestHandlerResponseHttpStatusOk(t *testing.T) {

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/starwar?q=data&page=1",
		nil,
	)
	rec := httptest.NewRecorder()

	Handle(rec, req)
	resp := rec.Result()

	assert.EqualValues(t, resp.StatusCode, http.StatusOK, "Http status code (%d) is not equals to 200", resp.StatusCode)
}

func TestHandlerResponseHttpStatusNotFound(t *testing.T) {

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/noteconosco?q=data&page=1",
		nil,
	)
	rec := httptest.NewRecorder()

	Handle(rec, req)
	resp := rec.Result()

	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(resp.Body)
	if err != nil {
		log.Println("error reading response body from controller mock")
		return
	}

	assert.EqualValues(t, resp.StatusCode, http.StatusNotFound, "Http status code (%d) is not equals to 404", resp.StatusCode)
	assert.EqualValues(t, buf.String(), "Url: /api/v1/noteconosco Not Found", "The body response is not equals to Url: /api/v1/noteconosco")
}
