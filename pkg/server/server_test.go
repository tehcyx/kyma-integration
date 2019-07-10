package server

import (
	"context"
	"net"
	"net/http"
	"reflect"
	"testing"

	"github.com/tehcyx/kyma-integration/pkg/kyma/certificate"
	"github.com/tehcyx/kyma-integration/internal/handler"
)

func TestNew(t *testing.T) {
	type args struct {
		host       string
		port       string
		securePort string
		handlers   handler.Param
	}
	tests := []struct {
		name string
		args args
		want *Server
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.host, tt.args.port, tt.args.securePort, tt.args.handlers); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_AddHandlers(t *testing.T) {
	type fields struct {
		Context       context.Context
		Host          string
		Port          string
		SecurePort    string
		Handlers      handler.Param
		ListenerNoTLS net.Listener
		ListenerTLS   net.Listener
		Client        *http.Client
		SecureClient  *http.Client
		Certificate   *certificate.CACertificate
		TLSPath       string
		AppName       string
	}
	type args struct {
		handlers handler.Param
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &Server{
				Context:       tt.fields.Context,
				Host:          tt.fields.Host,
				Port:          tt.fields.Port,
				SecurePort:    tt.fields.SecurePort,
				Handlers:      tt.fields.Handlers,
				ListenerNoTLS: tt.fields.ListenerNoTLS,
				ListenerTLS:   tt.fields.ListenerTLS,
				Client:        tt.fields.Client,
				SecureClient:  tt.fields.SecureClient,
				Certificate:   tt.fields.Certificate,
				TLSPath:       tt.fields.TLSPath,
				AppName:       tt.fields.AppName,
			}
			if err := srv.AddHandlers(tt.args.handlers); (err != nil) != tt.wantErr {
				t.Errorf("Server.AddHandlers() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServer_Run(t *testing.T) {
	type fields struct {
		Context       context.Context
		Host          string
		Port          string
		SecurePort    string
		Handlers      handler.Param
		ListenerNoTLS net.Listener
		ListenerTLS   net.Listener
		Client        *http.Client
		SecureClient  *http.Client
		Certificate   *certificate.CACertificate
		TLSPath       string
		AppName       string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &Server{
				Context:       tt.fields.Context,
				Host:          tt.fields.Host,
				Port:          tt.fields.Port,
				SecurePort:    tt.fields.SecurePort,
				Handlers:      tt.fields.Handlers,
				ListenerNoTLS: tt.fields.ListenerNoTLS,
				ListenerTLS:   tt.fields.ListenerTLS,
				Client:        tt.fields.Client,
				SecureClient:  tt.fields.SecureClient,
				Certificate:   tt.fields.Certificate,
				TLSPath:       tt.fields.TLSPath,
				AppName:       tt.fields.AppName,
			}
			srv.Run()
		})
	}
}

func TestServer_CertExists(t *testing.T) {
	type fields struct {
		Context       context.Context
		Host          string
		Port          string
		SecurePort    string
		Handlers      handler.Param
		ListenerNoTLS net.Listener
		ListenerTLS   net.Listener
		Client        *http.Client
		SecureClient  *http.Client
		Certificate   *certificate.CACertificate
		TLSPath       string
		AppName       string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &Server{
				Context:       tt.fields.Context,
				Host:          tt.fields.Host,
				Port:          tt.fields.Port,
				SecurePort:    tt.fields.SecurePort,
				Handlers:      tt.fields.Handlers,
				ListenerNoTLS: tt.fields.ListenerNoTLS,
				ListenerTLS:   tt.fields.ListenerTLS,
				Client:        tt.fields.Client,
				SecureClient:  tt.fields.SecureClient,
				Certificate:   tt.fields.Certificate,
				TLSPath:       tt.fields.TLSPath,
				AppName:       tt.fields.AppName,
			}
			if got := srv.CertExists(); got != tt.want {
				t.Errorf("Server.CertExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_GenerateKeysAndCertificate(t *testing.T) {
	type fields struct {
		Context       context.Context
		Host          string
		Port          string
		SecurePort    string
		Handlers      handler.Param
		ListenerNoTLS net.Listener
		ListenerTLS   net.Listener
		Client        *http.Client
		SecureClient  *http.Client
		Certificate   *certificate.CACertificate
		TLSPath       string
		AppName       string
	}
	type args struct {
		subject string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *certificate.CACertificate
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &Server{
				Context:       tt.fields.Context,
				Host:          tt.fields.Host,
				Port:          tt.fields.Port,
				SecurePort:    tt.fields.SecurePort,
				Handlers:      tt.fields.Handlers,
				ListenerNoTLS: tt.fields.ListenerNoTLS,
				ListenerTLS:   tt.fields.ListenerTLS,
				Client:        tt.fields.Client,
				SecureClient:  tt.fields.SecureClient,
				Certificate:   tt.fields.Certificate,
				TLSPath:       tt.fields.TLSPath,
				AppName:       tt.fields.AppName,
			}
			if got := srv.GenerateKeysAndCertificate(tt.args.subject); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Server.GenerateKeysAndCertificate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_StartListen(t *testing.T) {
	type fields struct {
		Context       context.Context
		Host          string
		Port          string
		SecurePort    string
		Handlers      handler.Param
		ListenerNoTLS net.Listener
		ListenerTLS   net.Listener
		Client        *http.Client
		SecureClient  *http.Client
		Certificate   *certificate.CACertificate
		TLSPath       string
		AppName       string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &Server{
				Context:       tt.fields.Context,
				Host:          tt.fields.Host,
				Port:          tt.fields.Port,
				SecurePort:    tt.fields.SecurePort,
				Handlers:      tt.fields.Handlers,
				ListenerNoTLS: tt.fields.ListenerNoTLS,
				ListenerTLS:   tt.fields.ListenerTLS,
				Client:        tt.fields.Client,
				SecureClient:  tt.fields.SecureClient,
				Certificate:   tt.fields.Certificate,
				TLSPath:       tt.fields.TLSPath,
				AppName:       tt.fields.AppName,
			}
			srv.StartListen()
		})
	}
}

func TestServer_StartListenTLS(t *testing.T) {
	type fields struct {
		Context       context.Context
		Host          string
		Port          string
		SecurePort    string
		Handlers      handler.Param
		ListenerNoTLS net.Listener
		ListenerTLS   net.Listener
		Client        *http.Client
		SecureClient  *http.Client
		Certificate   *certificate.CACertificate
		TLSPath       string
		AppName       string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &Server{
				Context:       tt.fields.Context,
				Host:          tt.fields.Host,
				Port:          tt.fields.Port,
				SecurePort:    tt.fields.SecurePort,
				Handlers:      tt.fields.Handlers,
				ListenerNoTLS: tt.fields.ListenerNoTLS,
				ListenerTLS:   tt.fields.ListenerTLS,
				Client:        tt.fields.Client,
				SecureClient:  tt.fields.SecureClient,
				Certificate:   tt.fields.Certificate,
				TLSPath:       tt.fields.TLSPath,
				AppName:       tt.fields.AppName,
			}
			srv.StartListenTLS()
		})
	}
}

func Test_getTLSPath(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTLSPath(); got != tt.want {
				t.Errorf("getTLSPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
