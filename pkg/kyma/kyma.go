package kyma

import (
	"github.com/tehcyx/kyma-integration/internal/handler"
	"github.com/tehcyx/kyma-integration/pkg/kyma/connector"
	"github.com/tehcyx/kyma-integration/pkg/server"
)

type Kyma struct {
	Serving   *server.Server
	Connector *connector.KymaConnector
}

// New creates a new kyma application
func New(host, port, securePort string, handlers handler.Param) *Kyma {
	srv := server.New(host, port, securePort, handlers)
	return &Kyma{
		Serving:   srv,
		Connector: connector.New(srv, "/connector"),
	}
}

// Run start serving
func (k *Kyma) Run() {
	k.Serving.Run()
}
