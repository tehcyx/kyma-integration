package certificate

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"time"
)

//TODO: double check : https://www.socketloop.com/tutorials/golang-create-x509-certificate-private-and-public-keys

func GenerateCSR(names pkix.Name, expiration time.Duration, size int) (*CACertificate, error) {
	keys, err := rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		return nil, fmt.Errorf("unable to generate private keys, error: %s", err)
	}
	type basicConstraints struct {
		IsCA       bool `asn1:"optional"`
		MaxPathLen int  `asn1:"optional,default:-1"`
	}
	val, err := asn1.Marshal(basicConstraints{true, 0})
	if err != nil {
		return nil, err
	}
	// step: generate a csr template
	var csrTemplate = x509.CertificateRequest{
		Subject:            names,
		SignatureAlgorithm: x509.SHA256WithRSA,
		ExtraExtensions: []pkix.Extension{
			{
				Id:       asn1.ObjectIdentifier{2, 5, 29, 19},
				Value:    val,
				Critical: true,
			},
		},
	}
	// step: generate the csr request
	csrCertificate, err := x509.CreateCertificateRequest(rand.Reader, &csrTemplate, keys)
	if err != nil {
		return nil, err
	}
	csr := pem.EncodeToMemory(&pem.Block{
		Type: "CERTIFICATE REQUEST", Bytes: csrCertificate,
	})
	// step: generate a serial number
	serial, err := rand.Int(rand.Reader, (&big.Int{}).Exp(big.NewInt(2), big.NewInt(159), nil))
	if err != nil {
		return nil, err
	}

	now := time.Now()
	// step: create the request template
	template := x509.Certificate{
		SerialNumber:          serial,
		Subject:               names,
		NotBefore:             now.Add(-10 * time.Minute).UTC(),
		NotAfter:              now.Add(expiration).UTC(),
		BasicConstraintsValid: true,
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
	}
	// step: sign the certificate authority
	certificate, err := x509.CreateCertificate(rand.Reader, &template, &template, &keys.PublicKey, keys)
	if err != nil {
		return nil, fmt.Errorf("failed to generate certificate, error: %s", err)
	}

	var request bytes.Buffer
	var privateKey bytes.Buffer
	if err := pem.Encode(&request, &pem.Block{Type: "CERTIFICATE", Bytes: certificate}); err != nil {
		return nil, err
	}
	if err := pem.Encode(&privateKey, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(keys)}); err != nil {
		return nil, err
	}

	return &CACertificate{
		PrivateKey: privateKey.String(),
		PublicKey:  request.String(),
		Csr:        string(csr),
	}, nil
}

func LoadClientCert(cert *CACertificate) (tls.Certificate, error) {
	// Load client cert
	return tls.LoadX509KeyPair(cert.ServerCertPath, cert.PrivateKeyPath)
}

func LoadServerCertBytes(cert *CACertificate) ([]byte, error) {
	// Load CA cert
	return ioutil.ReadFile(cert.ServerCertPath)
}

func CreateTLSConfig(cert *CACertificate) (*http.Transport, error) {
	clientCert, err := LoadClientCert(cert)
	if err != nil {
		return nil, err
	}

	serverCertBytes, err := LoadServerCertBytes(cert)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(serverCertBytes)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{clientCert},
		RootCAs:            caCertPool,
	}
	tlsConfig.BuildNameToCertificate()

	return &http.Transport{
		TLSClientConfig: tlsConfig,
	}, nil

}
