package main

import (
	"github.com/tehcyx/kyma-integration/internal/handler"
	"github.com/tehcyx/kyma-integration/pkg/api"
	"github.com/tehcyx/kyma-integration/pkg/frontend"
	"github.com/tehcyx/kyma-integration/pkg/kyma"
)

func main() {
	var handlers handler.Param
	handlers = make(handler.Param)

	handlers["/"] = frontend.IndexHandler
	handlers["/api/v1/test"] = api.TestHandler

	kyma := kyma.New("0.0.0.0", "8080", handlers)

	kyma.Run()
}
