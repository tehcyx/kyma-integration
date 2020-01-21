package kyma

import (
	"github.com/tehcyx/kyma-integration/internal/handler"
	"github.com/tehcyx/kyma-integration/pkg/kyma/connector"
	"github.com/tehcyx/kyma-integration/pkg/server"
)

// Kyma struct holding the app data, enabling kyma connections and endpoint serving.
type Kyma struct {
	Serving   *server.Server
	Connector *connector.KymaConnector
}

// New creates a new kyma application
func New(host, port string, handlers handler.Param) *Kyma {
	srv := server.New(host, port, handlers)
	return &Kyma{
		Serving:   srv,
		Connector: connector.New(srv, "/kyma"),
	}
}

// Run start serving
func (k *Kyma) Run() {
	k.Serving.Run()
}
