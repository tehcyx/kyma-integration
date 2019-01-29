package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type CertConnectResponse struct {
	Cert string `json:"crt,omitempty"`
}

type ApplicationConnectResponse struct {
	CsrURL      string          `json:"csrUrl,omitempty"`
	API         APIData         `json:"api,omitempty"`
	Certificate CertificateData `json:"certificate,omitempty"`
}

type APIData struct {
	MetadataURL     string `json:"metadataUrl,omitempty"`
	EventsURL       string `json:"eventsUrl,omitempty"`
	CertificatesURL string `json:"certificatesUrl,omitempty"`
}

type CertificateData struct {
	Subject      string `json:"subject,omitempty"`
	Extensions   string `json:"extensions,omitempty"`
	KeyAlgorithm string `json:"key-algorithm,omitempty"`
}

type CACertificate struct {
	privateKey string
	publicKey  string
	csr        string
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
