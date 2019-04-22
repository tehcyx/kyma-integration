package handler

import "net/http"

// Param map of path and Func's to be exposed by the server
type Param map[string]Func

// Func represents a default http handleFunc
type Func func(http.ResponseWriter, *http.Request)
