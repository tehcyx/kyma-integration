package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func (ks *KymaIntegrationServer) connectHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("handler started")
	defer log.Println("handler ended")

	params, ok := r.URL.Query()["url"]

	if !ok || len(params[0]) < 1 {
		log.Println("Url Param 'url' is missing")
		return
	}
	urlParam, err := url.Parse(params[0])
	if err != nil {
		fmt.Fprintln(w, "error need url param")
	}

	resp := ks.performGetOnURL(ctx, urlParam.String())

	appData := &ApplicationConnectResponse{}

	unmarshalInfoErr := json.Unmarshal([]byte(resp), appData)
	if unmarshalInfoErr != nil {
		fmt.Fprintln(w, "could not parse response")
	}

	ks.appInfo = appData

	resp = ks.sendCSRResponse(ctx, appData.CsrURL, appData.Certificate.Subject)

	certData := &CertConnectResponse{}

	unmarshalCertErr := json.Unmarshal([]byte(resp), certData)
	if unmarshalCertErr != nil {
		fmt.Fprintln(w, "could not parse response")
	}

	decodedCert, decodeErr := base64.StdEncoding.DecodeString(certData.Cert)
	if decodeErr != nil {
		fmt.Fprintf(w, "something went wrong decoding the response")
	}
	certData.Cert = string(decodedCert)
	certBytes := []byte(certData.Cert)
	errCert := ioutil.WriteFile(ks.serverCertPath, certBytes, 0644)
	if errCert != nil {
		log.Fatal("couldn't write server cert")
	}

	ks.startListenTLS()

	fmt.Fprintf(w, "Connected successfully: \n%v", ks.appInfo)
}
