package rest

import (
	"fmt"
	"net/http"
)

type API struct {
	V1 APIv1
}

type APIv1 struct {
}

func (a *APIv1) GetTestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "All good")
}
