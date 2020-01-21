package api

import "net/http"

// TestHandler handler function just returning ok.
func TestHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
