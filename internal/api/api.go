package api

import (
	"context"
	"net"
	"net/http"

	"github.com/tehcyx/kyma-integration/pkg/kyma/certificate"
)

// KymaIntegration interface for api functions
type KymaIntegration interface {
	Start()
	IndexHandler(w http.ResponseWriter, r *http.Request)
	GitHubCallbackHandler(w http.ResponseWriter, r *http.Request)
	ConnectHandler(w http.ResponseWriter, r *http.Request)
	RegisterServiceHandler(w http.ResponseWriter, r *http.Request)
	GetURL(ctx context.Context, queryURL string) string
	SendCSRResponse(ctx context.Context, responseURL, subject string) string
	GenerateKeysAndCertificate(subject string) *certificate.CACertificate
	TlsCertExists() bool
	StartListenTLS()
}

// KymaIntegrationServer struct for the server information
type KymaIntegrationServer struct {
	KymaIntegration
	Context                                         context.Context
	Cert                                            *certificate.CACertificate
	Dir, CsrPath, PubPath, PrivPath, ServerCertPath string
	Client                                          *http.Client
	SecureClient                                    *http.Client
	ListenNoTLS, ListenTLS                          net.Listener
	AppInfo                                         *certificate.ApplicationConnectResponse
}
