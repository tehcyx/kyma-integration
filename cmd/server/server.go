package main

import (
	"fmt"
	"net/http"
)

// func main() {

// 	var handlers handler.Param
// 	handlers = make(handler.Param)

// 	handlers["/"] = indexHandler

// 	srv := server.New("127.0.0.1", "8080", "8443", handlers)
// 	srv.Run()
// }

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world")
}
