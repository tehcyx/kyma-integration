package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/tehcyx/kyma-integration/internal/handler"
	"github.com/tehcyx/kyma-integration/pkg/kyma/certificate"
)

// Server struct to handle http exposure and TLS
type Server struct {
	Context              context.Context
	Host, Port           string
	Handlers             handler.Param
	Listener             net.Listener
	Client, SecureClient *http.Client
	Certificate          *certificate.CACertificate
	AppName              string
}

// New creates a new server allowing you to expose rest endpoints
func New(host, port string, handlers handler.Param) *Server {
	for path, hndl := range handlers {
		http.HandleFunc(path, hndl)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &Server{
		Host:         host,
		Port:         port,
		Handlers:     handlers,
		Listener:     nil,
		Client:       &http.Client{Transport: tr},
		SecureClient: nil,
		Certificate:  &certificate.CACertificate{},
	}
}

// AddHandlers adds more handlers
func (srv *Server) AddHandlers(handlers handler.Param) error {
	for path := range handlers {
		if _, ok := srv.Handlers[path]; ok {
			return fmt.Errorf("Handler with path '%s' already exists. Not adding any handlers from this call", path)
		}
	}
	for path, hndl := range handlers {
		http.HandleFunc(path, hndl)
		srv.Handlers[path] = hndl
	}
	return nil
}

// Run starts up the server
func (srv *Server) Run() {
	srv.StartListen()
}

// StartListen starts exposure of service on port of choice via http://
func (srv *Server) StartListen() {
	var err error
	srv.Listener, err = net.Listen("tcp", fmt.Sprintf("%s:%s", srv.Host, srv.Port))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("🔓 Listening on %s:%s", srv.Host, srv.Port)
	http.Serve(srv.Listener, nil)
}
