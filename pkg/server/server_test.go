package server

import (
	"context"
	"net"
	"net/http"
	"reflect"
	"testing"

	"github.com/tehcyx/kyma-integration/internal/handler"
	"github.com/tehcyx/kyma-integration/pkg/kyma/certificate"
)

func TestNew(t *testing.T) {
	type args struct {
		host     string
		port     string
		handlers handler.Param
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
			if got := New(tt.args.host, tt.args.port, tt.args.handlers); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_AddHandlers(t *testing.T) {
	type fields struct {
		Context      context.Context
		Host         string
		Port         string
		Handlers     handler.Param
		Listener     net.Listener
		Client       *http.Client
		SecureClient *http.Client
		Certificate  *certificate.CACertificate
		ConfigPath   string
		AppName      string
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
				Context:      tt.fields.Context,
				Host:         tt.fields.Host,
				Port:         tt.fields.Port,
				Handlers:     tt.fields.Handlers,
				Listener:     tt.fields.Listener,
				Client:       tt.fields.Client,
				SecureClient: tt.fields.SecureClient,
				Certificate:  tt.fields.Certificate,
				AppName:      tt.fields.AppName,
			}
			if err := srv.AddHandlers(tt.args.handlers); (err != nil) != tt.wantErr {
				t.Errorf("Server.AddHandlers() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServer_Run(t *testing.T) {
	type fields struct {
		Context      context.Context
		Host         string
		Port         string
		Handlers     handler.Param
		Listener     net.Listener
		Client       *http.Client
		SecureClient *http.Client
		Certificate  *certificate.CACertificate
		ConfigPath   string
		AppName      string
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
				Context:      tt.fields.Context,
				Host:         tt.fields.Host,
				Port:         tt.fields.Port,
				Handlers:     tt.fields.Handlers,
				Listener:     tt.fields.Listener,
				Client:       tt.fields.Client,
				SecureClient: tt.fields.SecureClient,
				Certificate:  tt.fields.Certificate,
				AppName:      tt.fields.AppName,
			}
			srv.Run()
		})
	}
}

func TestServer_StartListen(t *testing.T) {
	type fields struct {
		Context      context.Context
		Host         string
		Port         string
		Handlers     handler.Param
		Listener     net.Listener
		Client       *http.Client
		SecureClient *http.Client
		Certificate  *certificate.CACertificate
		ConfigPath   string
		AppName      string
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
				Context:      tt.fields.Context,
				Host:         tt.fields.Host,
				Port:         tt.fields.Port,
				Handlers:     tt.fields.Handlers,
				Listener:     tt.fields.Listener,
				Client:       tt.fields.Client,
				SecureClient: tt.fields.SecureClient,
				Certificate:  tt.fields.Certificate,
				AppName:      tt.fields.AppName,
			}
			srv.StartListen()
		})
	}
}
