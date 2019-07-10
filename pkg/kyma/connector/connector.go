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
	"os"

	"github.com/tehcyx/kyma-integration/internal/handler"

	"github.com/tehcyx/kyma-integration/pkg/kyma/certificate"
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
		log.Fatalf("Url Param 'url' is missing")
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
		log.Fatalf("couldn't write server cert: %s", errCert)
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
	serviceDescription.Documentation = new(ServiceDocumentation)
	serviceDescription.Documentation.DisplayName = "Test"
	serviceDescription.Documentation.Description = "test decsription"
	serviceDescription.Documentation.Tags = []string{"Tag1", "Tag2"}
	serviceDescription.Documentation.Type = "Test Type"

	serviceDescription.Description = "API Description"
	serviceDescription.ShortDescription = "API Short Description"

	serviceDescription.Provider = "Kyma example"
	serviceDescription.Name = "Kyma example service"

	serviceDescription.API = new(ServiceAPI)
	if envIP := os.Getenv("INSTANCE_IP"); envIP != "" {
		serviceDescription.API.TargetURL = fmt.Sprintf("%s:8080", envIP)
	} else {
		serviceDescription.API.TargetURL = "http://localhost:8080"
	}
	serviceDescription.API.Spec = json.RawMessage(`{
		"swagger":"2.0",
		"info":{  
		   "description":"Kyma API example",
		   "version":"1.0",
		   "title":"Kyma example",
		   "contact":{  
			  "name":"Daniel Roth",
			  "email":"email@email.com"
		   }
		},
		"host":"localhost:8080",
		"basePath":"/",
		"tags":[  
		   {  
			  "name":"kyma-integration-golang-api",
			  "description":"Kyma integration Golang Api"
		   }
		],
		"paths":{  
		   "/api/v1/test":{  
			  "get":{  
				 "tags":[  
					"kyma-integration-golang-api"
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

	serviceDescription.Events = new(ServiceEvent)
	serviceDescription.Events.Spec = json.RawMessage(`{
		"asyncapi": "1.0.0",
		"info": {
			"title": "Example Events",
			"version": "1.0.0",
			"description": "Description of all the example events"
		},
		"baseTopic": "example.events.com",
		"topics": {
			"exampleEvent.v1": {
				"subscribe": {
					"summary": "Example event",
					"payload": {
						"type": "object",
						"properties": {
							"myObject": {
								"type": "object",
								"required": [
									"id"
								],
								"example": {
									"id": "4caad296-e0c5-491e-98ac-0ed118f9474e"
								},
								"properties": {
									"id": {
										"title": "Id",
										"description": "Resource identifier",
										"type": "string"
									}
								}
							}
						}
					}
				}
			}
		}
	}`)

	jsonBytes, err := json.Marshal(serviceDescription)
	if err != nil {
		log.Fatalf("JSON marshal failed: %s", err)
		return
	}

	if kc.AppInfo == nil || kc.AppInfo.API.MetadataURL == "" {
		log.Fatalf("%s", fmt.Errorf("metadata url is missing, cannot proceed"))
	}

	req, err := http.NewRequest("POST", kc.AppInfo.API.MetadataURL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Fatalf("Couldn't create request to register service: %s", err)
	}
	req.WithContext(ctx)

	resp, err := kc.Serving.SecureClient.Do(req)
	if err != nil {
		log.Fatalf("Couldn't register service: %s", err)
	}
	dump, err := httputil.DumpResponse(resp, true)
	defer resp.Body.Close() // close body after using it
	if err != nil {
		log.Fatalf("could not dump response: %v", err)
	}
	fmt.Printf("%s\n", dump)
	bodyString := string(dump)

	if resp.StatusCode == http.StatusOK {
		log.Printf("Successfully registered service with")
		fmt.Fprintf(w, bodyString)
	} else {
		fmt.Fprintf(w, "Status: %d >%s< \n on URL: %s", resp.StatusCode, bodyString, kc.AppInfo.API.MetadataURL)
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
