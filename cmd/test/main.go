package main

import (
	"fmt"
	"net"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/tehcyx/kyma-integration/pkg/kyma/config"
)

func main() {

	go conf()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("index handler")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", "0.0.0.0", "8080"))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("🔓 Listening on %s:%s", "0.0.0.0", "8080")
	http.Serve(listener, nil)
}

func conf() {
	for {
		appConfig := config.New()

		log.Printf("--- appConfig:\n%+v\n\n", appConfig)

		time.Sleep(10 * time.Second)

	}
}
