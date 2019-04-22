package connector

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/tehcyx/kyma-integration/internal/handler"

	"github.com/tehcyx/kyma-integration/internal/certificate"
	"github.com/tehcyx/kyma-integration/pkg/server"
)

// KymaConnector holds all information and functionality regarding Kyma
type KymaConnector struct {
	Serving       *server.Server
	AppInfo       *certificate.ApplicationConnectResponse
	servicePrefix string
}

// New Kyma one time init factory
func New(srv *server.Server, prefix string) *KymaConnector {
	kc := &KymaConnector{
		Serving:       srv,
		servicePrefix: prefix,
	}
	handlers := make(handler.Param)

	handlers[fmt.Sprintf("%s%s", prefix, "/connect")] = kc.connectHandler
	handlers[fmt.Sprintf("%s%s", prefix, "/register-service")] = kc.registerServiceHandler

	kc.Serving.AddHandlers(handlers)

	return kc
}

func (kc *KymaConnector) getResponseBodyWithContext(ctx context.Context, url string) string {
	req, err := http.NewRequest("GET", url, nil)
	req = req.WithContext(ctx)

	resp, err := kc.Serving.Client.Do(req)
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

func (kc *KymaConnector) connectHandler(w http.ResponseWriter, r *http.Request) {
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

	resp := kc.getResponseBodyWithContext(ctx, urlParam.String())

	appData := &certificate.ApplicationConnectResponse{}

	unmarshalInfoErr := json.Unmarshal([]byte(resp), appData)
	if unmarshalInfoErr != nil {
		fmt.Fprintln(w, "could not parse response")
	}

	kc.AppInfo = appData

	resp = kc.SendCSRResponse(ctx, appData.CsrURL, appData.Certificate.Subject)

	certData := &certificate.CertConnectResponse{}

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
	errCert := ioutil.WriteFile(kc.Serving.Certificate.ServerCertPath, certBytes, 0644)
	if errCert != nil {
		log.Fatal("couldn't write server cert")
	}

	kc.Serving.StartListenTLS()

	fmt.Fprintf(w, "Connected successfully: \n%v", kc.AppInfo)
}

func (kc *KymaConnector) registerServiceHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if kc.AppInfo == nil {
		fmt.Fprintf(w, "It seems the server crashed since we connected, so currently you need to reconnect it to use this endpoint.")
	}

	serviceDescription := new(Service)

	// Documentation part of the serviceDescription broken: https://github.com/kyma-project/kyma/issues/3347
	// serviceDescription.Documentation = new(ServiceDocumentation)
	// serviceDescription.Documentation.DisplayName = "Test"
	// serviceDescription.Documentation.Description = "test decsription"
	// serviceDescription.Documentation.Tags = []string{"Tag1", "Tag2"}
	// serviceDescription.Documentation.Type = "Test Type"

	serviceDescription.Description = "API Description"
	serviceDescription.ShortDescription = "API Short Description"

	serviceDescription.Provider = "Daniel"
	serviceDescription.Name = "Daniel's Service"

	serviceDescription.API = new(ServiceAPI)
	serviceDescription.API.TargetURL = "http://10.48.60.12:8080"
	serviceDescription.API.Spec = json.RawMessage(`{
		"swagger":"2.0",
		"info":{  
		   "description":"API example",
		   "version":"1.0",
		   "title":"Github Kubernetes API",
		   "contact":{  
			  "name":"Daniel Roth",
			  "email":"daniel.roth02@sap.com"
		   }
		},
		"host":"10.48.60.12:8080",
		"basePath":"/",
		"tags":[  
		   {  
			  "name":"github-api",
			  "description":"Github Api"
		   }
		],
		"paths":{  
		   "/api/v1/test":{  
			  "get":{  
				 "tags":[  
					"github-api"
				 ],
				 "summary":"Test",
				 "description":"this is a test",
				 "operationId":"opId1",
				 "produces":[  
					"*/*"
				 ],
				 "responses":{  
					"200":{  
					   "description":"OK"
					}
				 },
				 "deprecated":false
			  }
		   }
		}
	 }`)

	jsonBytes, err := json.Marshal(serviceDescription)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(jsonBytes))

	if kc.Serving.AppName == "" {
		kc.Serving.AppName = "github-test"
	}

	// acquire NodePort to modify URL locally: kubectl -n kyma-system get svc application-connector-ingress-nginx-ingress-controller -o 'jsonpath={.spec.ports[?(@.port==443)].nodePort}'
	// ks.appInfo.API.MetadataURL = "https://gateway.kyma.local:31615/github-test/v1/metadata/services"
	metadataURL := fmt.Sprintf("https://gateway.kyma.local:31615/%s/v1/metadata/services", kc.Serving.AppName)
	// if ks.appInfo != nil && ks.appInfo.API.MetadataURL != "" {
	// 	metadataURL = ks.appInfo.API.MetadataURL
	// }

	req, err := http.NewRequest("POST", metadataURL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Printf("Couldn't register service: %s", err)
	}
	req.WithContext(ctx)

	resp, err := kc.Serving.SecureClient.Do(req)
	if err != nil {
		log.Printf("Couldn't register service: %s", err)
	}
	dump, err := httputil.DumpResponse(resp, true)
	defer resp.Body.Close() // close body after using it
	if err != nil {
		log.Printf("could not dump response: %v", err)
	}
	fmt.Printf("%s\n", dump)
	bodyString := string(dump)

	// bodyBytes, _ := ioutil.ReadAll(resp.Body)
	// bodyString := string(bodyBytes)

	if resp.StatusCode == http.StatusOK {
		log.Printf("Successfully registered service with ID %s", "id here")
		fmt.Fprintf(w, bodyString)
	} else {
		fmt.Fprintf(w, "Status: %d >%s< \n on URL: %s", resp.StatusCode, bodyString, metadataURL)
	}
}

func (kc *KymaConnector) SendCSRResponse(ctx context.Context, responseURL, subject string) string {
	kc.Serving.Certificate = kc.Serving.GenerateKeysAndCertificate(subject)

	var jsonStr = []byte(fmt.Sprintf("{\"csr\":\"%s\"}", base64.StdEncoding.EncodeToString([]byte(kc.Serving.Certificate.Csr))))
	req, err := http.NewRequest("POST", responseURL, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)

	resp, err := kc.Serving.Client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	return bodyString
}
