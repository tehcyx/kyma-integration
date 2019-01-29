package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"time"
)

func (ks *KymaIntegrationServer) tlsCertExists() bool {
	_, errCert := os.Stat(ks.serverCertPath)
	if errCert == nil {
		return true
	} else if os.IsNotExist(errCert) {
		return false
	} else {
		log.Fatal("read error on cert file")
		return false
	}
}

func (ks *KymaIntegrationServer) generateKeysAndCertificate(subject string) *CACertificate {
	var cert *CACertificate
	cert = new(CACertificate)

	if !ks.tlsCertExists() {
		_, errCSR := os.Stat(ks.csrPath)
		_, errPub := os.Stat(ks.pubPath)
		_, errPriv := os.Stat(ks.privPath)

		// read cert.csr
		if errCSR == nil && errPub == nil && errPriv == nil {
			csrBytes, err := ioutil.ReadFile(ks.csrPath)
			if err != nil {
				log.Fatal("Read error on csr file")
			}
			cert.csr = string(csrBytes[:])
			pubKeyBytes, err := ioutil.ReadFile(ks.pubPath)
			if err != nil {
				log.Fatal("Read error on pub file")
			}
			cert.publicKey = string(pubKeyBytes[:])
			privKeyBytes, err := ioutil.ReadFile(ks.privPath)
			if err != nil {
				log.Fatal("Read error on priv file")
			}
			cert.privateKey = string(privKeyBytes[:])
		} else if os.IsNotExist(errCSR) && os.IsNotExist(errPub) && os.IsNotExist(errPriv) {
			subject := pkix.Name{
				Locality:           []string{"Waldorf"},
				Province:           []string{"Waldorf"},
				Country:            []string{"DE"},
				Organization:       []string{"Organization"},
				OrganizationalUnit: []string{"OrgUnit"},
				CommonName:         "github-test",
				// ??:              []string{"OU=OrgUnit,O=Organization,L=Waldorf,ST=Waldorf,C=DE,CN=github-test"},
			}

			// {pkix.Name{
			// 	CommonName:         "Steve Kille",
			// 	Organization:       []string{"Isode Limited"},
			// 	OrganizationalUnit: []string{"RFCs"},
			// 	Locality:           []string{"Richmond"},
			// 	Province:           []string{"Surrey"},
			// 	StreetAddress:      []string{"The Square"},
			// 	PostalCode:         []string{"TW9 1DT"},
			// 	SerialNumber:       "RFC 2253",
			// 	Country:            []string{"GB"},
			// }, "SERIALNUMBER=RFC 2253,CN=Steve Kille,OU=RFCs,O=Isode Limited,POSTALCODE=TW9 1DT,STREET=The Square,L=Richmond,ST=Surrey,C=GB"},

			genCert, err := generateCSR(subject, time.Duration(1200), 2048)
			if err != nil {
				fmt.Println(err)
			}
			//write files here
			csrBytes := []byte(genCert.csr)
			pubKeyBytes := []byte(genCert.publicKey)
			privKeyBytes := []byte(genCert.privateKey)
			errCSR := ioutil.WriteFile(ks.csrPath, csrBytes, 0644)
			if errCSR != nil {
				log.Fatal("couldn't write csr")
			}
			errPub := ioutil.WriteFile(ks.pubPath, pubKeyBytes, 0644)
			if errPub != nil {
				log.Fatal("couldn't write pub key")
			}
			errPriv := ioutil.WriteFile(ks.privPath, privKeyBytes, 0644)
			if errPriv != nil {
				log.Fatal("couldn't write priv key")
			}
			cert = genCert
		} else {
			log.Fatal("cert not readable or does not exist")
		}
	}

	return cert
}

func generateCSR(names pkix.Name, expiration time.Duration, size int) (*CACertificate, error) {
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
		SignatureAlgorithm: x509.SHA512WithRSA,
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
		privateKey: privateKey.String(),
		publicKey:  request.String(),
		csr:        string(csr),
	}, nil
}
