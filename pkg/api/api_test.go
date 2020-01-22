package api

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func assertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a != b {
		t.Errorf("%s: Received %v (type %v), expected %v (type %v)", message, a, reflect.TypeOf(a), b, reflect.TypeOf(b))
	}
}

func TestTestHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "localhost:8080/api/v1/test", nil)
	w := httptest.NewRecorder()

	TestHandler(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	assertEqual(t, string(body), "ok", "handler should return message 'ok'")
	assertEqual(t, resp.StatusCode, http.StatusOK, "handler should return status '200'")
}
