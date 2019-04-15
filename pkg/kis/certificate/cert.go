package certificate

import (
	"bytes"
	"context"
	"crypto/x509/pkix"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	cert "github.com/tehcyx/kyma-github-integration/internal/certificate"
	certificate "github.com/tehcyx/kyma-github-integration/internal/certificate"
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

func (ks *KymaIntegrationServer) generateKeysAndCertificate(subject string) *cert.CACertificate {
	var cert *certificate.CACertificate
	cert = new(certificate.CACertificate)

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

			genCert, err := certificate.GenerateCSR(subject, time.Duration(1200), 2048)
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

func (ks *KymaIntegrationServer) sendCSRResponse(ctx context.Context, responseURL, subject string) string {
	ks.cert = ks.generateKeysAndCertificate(subject)

	var jsonStr = []byte(fmt.Sprintf("{\"csr\":\"%s\"}", base64.StdEncoding.EncodeToString([]byte(ks.cert.csr))))
	req, err := http.NewRequest("POST", responseURL, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)

	resp, err := ks.httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	return bodyString
}

func (ks *KymaIntegrationServer) performGetOnURL(ctx context.Context, queryURL string) string {
	req, err := http.NewRequest("GET", queryURL, nil)
	req = req.WithContext(ctx)

	resp, err := ks.httpClient.Do(req)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		return bodyString
	}
	return fmt.Sprintf("Status: %d, Message: %s", resp.StatusCode, "error")
}
