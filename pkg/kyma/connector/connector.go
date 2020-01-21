package connector

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/tehcyx/kyma-integration/internal/handler"
	"github.com/tehcyx/kyma-integration/pkg/kyma/certificate"
	"github.com/tehcyx/kyma-integration/pkg/kyma/config"
	"github.com/tehcyx/kyma-integration/pkg/server"
)

// KymaConnector holds all information and functionality regarding Kyma
type KymaConnector struct {
	Serving       *server.Server
	AppInfo       *certificate.ApplicationConnectResponse
	AppConfig     config.Config
	servicePrefix string
}

var serviceDescription *Service

// init generates a new service description on package import
func init() {
	serviceDescription = new(Service)

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
}

// New Kyma one time init factory.
func New(srv *server.Server, prefix string) *KymaConnector {
	kc := &KymaConnector{
		Serving:       srv,
		servicePrefix: prefix,
	}
	handlers := make(handler.Param)

	handlers[fmt.Sprintf("%s%s", prefix, "/connect")] = kc.connectHandler
	handlers[fmt.Sprintf("%s%s", prefix, "/connect/auto")] = kc.autoConnectHandler
	handlers[fmt.Sprintf("%s%s", prefix, "/register-service")] = kc.registerServiceHandler

	kc.Serving.AddHandlers(handlers)
	kc.AppConfig = config.New()

	return kc
}

func (kc *KymaConnector) getResponseBodyWithContext(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	req = req.WithContext(ctx)

	resp, err := kc.Serving.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		return bodyString, nil
	}
	return "", fmt.Errorf("response was not 200 as expected but %d instead", resp.StatusCode)
}

func (kc *KymaConnector) connectHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	err := kc.connectApplicationPOST(ctx, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Errorf("failed to connect application: %w", err).Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Connected successfully: \n%v", kc.AppInfo)))
}

func (kc *KymaConnector) autoConnectHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log.Debugf("Connect request coming in via %s", r.Method)
	switch r.Method {
	case http.MethodPost:
		err := kc.connectApplicationPOST(ctx, r)
		if err != nil {
			log.Printf("failed to connect application: %w", err)
			// redirect back to referer and mark with error
			redirectURL := fmt.Sprintf("%s?error", r.Header.Get("Referer"))
			http.Redirect(w, r, redirectURL, 302)
			return
		}
	default:
		log.Printf("failed to connect application: method not supported")
		// redirect back to referer and mark with error
		redirectURL := fmt.Sprintf("%s?error", r.Header.Get("Referer"))
		http.Redirect(w, r, redirectURL, 302)
		return
	}
	message, err := kc.registerService(ctx)
	if err != nil {
		log.Printf("failed to register service: %w", err)
		// redirect back to referer and mark with error
		redirectURL := fmt.Sprintf("%s?error", r.Header.Get("Referer"))
		http.Redirect(w, r, redirectURL, 302)
		return
	}

	resp := RegisterResponse{}
	log.Println(string(message))
	jsonErr := json.Unmarshal(message, &resp)
	if jsonErr != nil {
		log.Printf("failed to unmarshal service id: %w", jsonErr)
		// redirect back to referer and mark with error
		redirectURL := fmt.Sprintf("%s?error", r.Header.Get("Referer"))
		http.Redirect(w, r, redirectURL, 302)
	}

	// save ID
	kc.AppConfig.UpdateAppID(resp.ID)

	// redirect back to referer and set this current page as referer
	redirectURL := fmt.Sprintf("%s?redirect", r.Header.Get("Referer"))
	http.Redirect(w, r, redirectURL, 302)
}

func (kc *KymaConnector) registerServiceHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	message, err := kc.registerService(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Errorf("failed to register service: %w", err).Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))
}

// SendCSRResponse sends a POST request with a newly generated certificate signing request response to the passed URL.
func (kc *KymaConnector) SendCSRResponse(ctx context.Context, responseURL, subject string) (string, error) {
	kc.Serving.Certificate = kc.AppConfig.GenerateKeysAndCertificate(subject)

	var jsonStr = []byte(fmt.Sprintf("{\"csr\":\"%s\"}", base64.StdEncoding.EncodeToString([]byte(kc.Serving.Certificate.Csr))))
	req, err := http.NewRequest("POST", responseURL, bytes.NewBuffer(jsonStr))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)

	resp, err := kc.Serving.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("csr failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	return bodyString, nil
}

func (kc *KymaConnector) connectApplicationPOST(ctx context.Context, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return fmt.Errorf("ParseForm() failed to parse POST form data: %w", err)
	}
	log.Println(r.Form)
	urlVar := r.FormValue("url")

	if urlVar == "" {
		return fmt.Errorf("Url Param 'url' is missing")
	}
	parsedURL, urlParseErr := url.Parse(urlVar)
	if urlParseErr != nil {
		return fmt.Errorf("need url param: %w", urlParseErr)
	}

	return kc.connectApplication(ctx, parsedURL.String())
}

func (kc *KymaConnector) connectApplication(ctx context.Context, url string) error {
	resp, err := kc.getResponseBodyWithContext(ctx, url)
	if err != nil {
		return fmt.Errorf("connect application request failed: %w", err)
	}

	appData := &certificate.ApplicationConnectResponse{}

	unmarshalInfoErr := json.Unmarshal([]byte(resp), appData)
	if unmarshalInfoErr != nil {
		return fmt.Errorf("failed to unmarshal json response: %w", unmarshalInfoErr)
	}
	kc.AppInfo = appData

	resp, csrErr := kc.SendCSRResponse(ctx, appData.CsrURL, appData.Certificate.Subject)
	if csrErr != nil {
		return fmt.Errorf("csr failed %w", csrErr)
	}

	certData := &certificate.CertConnectResponse{}

	unmarshalCertErr := json.Unmarshal([]byte(resp), certData)
	if unmarshalCertErr != nil {
		fmt.Errorf("could not parse response: %w", unmarshalCertErr)
	}

	decodedCert, decodeErr := base64.StdEncoding.DecodeString(certData.Cert)
	if decodeErr != nil {
		fmt.Errorf("something went wrong decoding the response: %w", decodeErr)
	}
	certData.Cert = string(decodedCert)
	// store cert in config
	kc.AppConfig.UpdateServerCert(certData.Cert)
	return nil
}

func (kc *KymaConnector) registerService(ctx context.Context) ([]byte, error) {
	if kc.AppInfo == nil {
		return []byte{}, fmt.Errorf("remote application data not in memory, can't register service")
	}
	if kc.AppInfo == nil || kc.AppInfo.API.MetadataURL == "" {
		return []byte{}, fmt.Errorf("metadata url is missing, cannot proceed")
	}
	jsonBytes, err := json.Marshal(serviceDescription)
	if err != nil {
		return []byte{}, fmt.Errorf("JSON marshal failed: %w", err)
	}

	req, err := http.NewRequest("POST", kc.AppInfo.API.MetadataURL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return []byte{}, fmt.Errorf("couldn't create request to register service: %w", err)
	}
	req.WithContext(ctx)

	client, cliErr := secureClientInit(kc.AppConfig)
	if cliErr != nil {
		return []byte{}, fmt.Errorf("error creating secure client: %w", cliErr)
	}
	kc.Serving.SecureClient = client

	resp, err := kc.Serving.SecureClient.Do(req)
	if err != nil {
		return []byte{}, fmt.Errorf("couldn't register service: %w", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close() // close body after using it
	if err != nil {
		return []byte{}, fmt.Errorf("could not read response: %w", err)
	}

	if resp.StatusCode == http.StatusOK {
		log.Debug("Successfully registered service with")
		return body, nil
	}
	return []byte{}, fmt.Errorf("status: %d >%s< \n on URL: %s", resp.StatusCode, string(body), kc.AppInfo.API.MetadataURL)
}

func secureClientInit(cfg config.Config) (*http.Client, error) {
	tr, transErr := createTransport(cfg)
	if transErr != nil {
		return nil, fmt.Errorf("error creating secure transport config: %w", transErr)
	}
	return &http.Client{Transport: tr}, nil
}

func createTransport(cfg config.Config) (*http.Transport, error) {
	clientCert, x509Err := tls.X509KeyPair([]byte(cfg.ServerCert), []byte(cfg.PrivateKey))
	if x509Err != nil {
		return nil, fmt.Errorf("loading x509 key pair failed: %w", x509Err)
	}

	serverCert := []byte(cfg.ServerCert)

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(serverCert)

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
