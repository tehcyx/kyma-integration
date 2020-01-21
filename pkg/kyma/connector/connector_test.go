package connector

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/tehcyx/kyma-integration/pkg/kyma/certificate"
	"github.com/tehcyx/kyma-integration/pkg/server"
)

func TestNew(t *testing.T) {
	type args struct {
		srv    *server.Server
		prefix string
	}
	tests := []struct {
		name string
		args args
		want *KymaConnector
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.srv, tt.args.prefix); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKymaConnector_getResponseBodyWithContext(t *testing.T) {
	type fields struct {
		Serving       *server.Server
		AppInfo       *certificate.ApplicationConnectResponse
		servicePrefix string
	}
	type args struct {
		ctx context.Context
		url string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
		want2  error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kc := &KymaConnector{
				Serving:       tt.fields.Serving,
				AppInfo:       tt.fields.AppInfo,
				servicePrefix: tt.fields.servicePrefix,
			}
			if got, got2 := kc.getResponseBodyWithContext(tt.args.ctx, tt.args.url); got != tt.want {
				t.Errorf("KymaConnector.getResponseBodyWithContext() = %v, %v want %v, %v", got, tt.want, got2, tt.want2)
			}
		})
	}
}

func TestKymaConnector_connectHandler(t *testing.T) {
	type fields struct {
		Serving       *server.Server
		AppInfo       *certificate.ApplicationConnectResponse
		servicePrefix string
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kc := &KymaConnector{
				Serving:       tt.fields.Serving,
				AppInfo:       tt.fields.AppInfo,
				servicePrefix: tt.fields.servicePrefix,
			}
			kc.connectHandler(tt.args.w, tt.args.r)
		})
	}
}

func TestKymaConnector_registerServiceHandler(t *testing.T) {
	type fields struct {
		Serving       *server.Server
		AppInfo       *certificate.ApplicationConnectResponse
		servicePrefix string
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kc := &KymaConnector{
				Serving:       tt.fields.Serving,
				AppInfo:       tt.fields.AppInfo,
				servicePrefix: tt.fields.servicePrefix,
			}
			kc.registerServiceHandler(tt.args.w, tt.args.r)
		})
	}
}

func TestKymaConnector_SendCSRResponse(t *testing.T) {
	type fields struct {
		Serving       *server.Server
		AppInfo       *certificate.ApplicationConnectResponse
		servicePrefix string
	}
	type args struct {
		ctx         context.Context
		responseURL string
		subject     string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
		want2  error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kc := &KymaConnector{
				Serving:       tt.fields.Serving,
				AppInfo:       tt.fields.AppInfo,
				servicePrefix: tt.fields.servicePrefix,
			}
			if got, got2 := kc.SendCSRResponse(tt.args.ctx, tt.args.responseURL, tt.args.subject); got != tt.want || got2 != tt.want2 {
				t.Errorf("KymaConnector.SendCSRResponse() = %v, %v want %v, %v", got, tt.want, got2, tt.want2)
			}
		})
	}
}
