package github

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func (ks *KymaIntegrationServer) gitHubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		err := ctx.Err()
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	default:
		time.Sleep(5 * time.Second)
		fmt.Fprintln(w, "hello")
	}
}
