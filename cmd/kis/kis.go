package main

import (
	"github.com/tehcyx/kyma-integration/internal/handler"
	"github.com/tehcyx/kyma-integration/pkg/kyma"
)

func main() {
	var handlers handler.Param
	handlers = make(handler.Param)

	kyma := kyma.New("0.0.0.0", "8080", "8443", handlers)

	kyma.Run()
}
