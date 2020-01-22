package kyma

import (
	"reflect"
	"testing"

	"github.com/tehcyx/kyma-integration/internal/handler"
	"github.com/tehcyx/kyma-integration/pkg/kyma/connector"
	"github.com/tehcyx/kyma-integration/pkg/server"
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
		want *Kyma
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

func TestKyma_Run(t *testing.T) {
	type fields struct {
		Serving   *server.Server
		Connector *connector.KymaConnector
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &Kyma{
				Serving:   tt.fields.Serving,
				Connector: tt.fields.Connector,
			}
			k.Run()
		})
	}
}
