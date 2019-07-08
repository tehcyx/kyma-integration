package certificate

import (
	"crypto/tls"
	"crypto/x509/pkix"
	"net/http"
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

func TestLoadClientCert(t *testing.T) {
	type args struct {
		cert *CACertificate
	}
	tests := []struct {
		name    string
		args    args
		want    tls.Certificate
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadClientCert(tt.args.cert)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadClientCert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadClientCert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadServerCertBytes(t *testing.T) {
	type args struct {
		cert *CACertificate
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadServerCertBytes(tt.args.cert)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadServerCertBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadServerCertBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateTLSConfig(t *testing.T) {
	type args struct {
		cert *CACertificate
	}
	tests := []struct {
		name    string
		args    args
		want    *http.Transport
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateTLSConfig(tt.args.cert)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTLSConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateTLSConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
