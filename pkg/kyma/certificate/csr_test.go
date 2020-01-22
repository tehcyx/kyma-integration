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
		want    interface{}
		wantErr bool
	}{ //"OU=OrgUnit,O=Organization,L=Waldorf,ST=Waldorf,C=DE,CN=api-test"
		{name: "GenerateCSR should not fail", want: &CACertificate{}, wantErr: false, args: args{names: pkix.Name{CommonName: "hello", OrganizationalUnit: []string{"unit"}, Organization: []string{"org"}, Locality: []string{"city"}, Province: []string{"state"}, Country: []string{"country"}}, expiration: time.Second * 4, size: 2048}},
		{name: "GenerateCSR should fail if key size not big enough", want: &CACertificate{}, wantErr: true, args: args{names: pkix.Name{CommonName: "hello", OrganizationalUnit: []string{"unit"}, Organization: []string{"org"}, Locality: []string{"city"}, Province: []string{"state"}, Country: []string{"country"}}, expiration: time.Second * 4, size: 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateCSR(tt.args.names, tt.args.expiration, tt.args.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateCSR() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("GenerateCSR() = %v, want %v", reflect.TypeOf(got), reflect.TypeOf(tt.want))
			}
		})
	}
}
