package certificate

import (
	"crypto/x509/pkix"
	"reflect"
	"testing"
	"time"
)

func TestGenerateCSR(t *testing.T) {
	type args struct {
		names      pkix.Name
		expiration time.Duration
		size       int
	}
	tests := []struct {
		name    string
		args    args
		want    *CACertificate
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateCSR(tt.args.names, tt.args.expiration, tt.args.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateCSR() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateCSR() = %v, want %v", got, tt.want)
			}
		})
	}
}
