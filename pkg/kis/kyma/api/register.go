package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"

	"log"
)

// Service kyma service struct
type Service struct {
	Provider         string                `json:"provider,omitempty"`
	Name             string                `json:"name,omitempty"`
	Description      string                `json:"description,omitempty"`
	ShortDescription string                `json:"shortDescription,omitempty"`
	Labels           *ServiceLabel         `json:"labels,omitempty"`
	API              *ServiceAPI           `json:"api,omitempty"`
	Events           *ServiceEvent         `json:"events,omitempty"`
	Documentation    *ServiceDocumentation `json:"documentation,omitempty"`
}

// ServiceLabel kyma service labels
type ServiceLabel map[string]string

// ServiceAPI kyma service api definition
type ServiceAPI struct {
	TargetURL   string              `json:"targetUrl,omitempty"`
	Spec        json.RawMessage     `json:"spec,omitempty"`
	Credentials *ServiceCredentials `json:"credentials,omitempty"`
}

// ServiceCredentials kyma service credentials definition
type ServiceCredentials struct {
	Basic *ServiceBasicCredentials `json:"basic,omitempty"`
	OAuth *ServiceOAuthCredentials `json:"oauth,omitempty"`
}

// ServiceBasicCredentials kyma basic auth service credentials
type ServiceBasicCredentials struct {
	ClientID string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// ServiceOAuthCredentials kyma oauth service credentials
type ServiceOAuthCredentials struct {
	ClientID     string `json:"clientId,omitempty"`
	ClientSecret string `json:"clientSecret,omitempty"`
}

// ServiceEvent kyma service event definition
type ServiceEvent struct {
	Spec *ServiceEventSpec `json:"spec,omitempty"`
}

// ServiceEventSpec kyma service event spec definition
type ServiceEventSpec struct {
	AsyncAPI  string                            `json:"asyncapi,omitempty"`
	Info      *ServiceEventSpecInfo             `json:"info,omitempty"`
	BaseTopic string                            `json:"baseTopic,omitempty"`
	Topics    map[string]*ServiceEventSpecTopic `json:"topics,omitempty"`
}

// ServiceEventSpecInfo kyma service event spec info definition
type ServiceEventSpecInfo struct {
	Title       string `json:"title,omitempty"`
	Version     string `json:"version,omitempty"`
	Description string `json:"description,omitempty"`
}

// ServiceEventSpecTopic kyma service event spec topic definition
type ServiceEventSpecTopic map[string]*ServiceEventSpecTopicDetail

// ServiceEventSpecTopicDetail kyma service event spec topic detail definition
type ServiceEventSpecTopicDetail struct {
	Summary string                 `json:"summary,omitempty"`
	Payload map[string]interface{} `json:"payload,omitempty"`
}

// ServiceDocumentation kyma service documentation definition
type ServiceDocumentation struct {
	DisplayName string                     `json:"displayName,omitempty"`
	Description string                     `json:"description,omitempty"`
	Type        string                     `json:"type,omitempty"`
	Tags        []string                   `json:"tags,omitempty"`
	Docs        []*ServiceDocumentationDoc `json:"docs,omitempty"`
}

// ServiceDocumentationDoc kyma service documentation doc definition
type ServiceDocumentationDoc struct {
	Title  string `json:"title,omitempty"`
	Type   string `json:"type,omitempty"`
	Source string `json:"source,omitempty"`
}

func (ks *KymaIntegrationServer) registerServiceHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if ks.appInfo == nil {
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
	serviceDescription.API.TargetURL = "http://kis.local:8080"
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
		"host":"kis.local:8080",
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

	// acquire NodePort to modify URL locally
	// 30019
	// https://gateway.kyma.local:31635/github-test/v1/metadata/services
	// ks.appInfo.API.MetadataURL = "https://gateway.kyma.local:31654/github-test/v1/metadata/services"

	req, err := http.NewRequest("POST", "https://gateway.kyma.local:31635/github-test/v1/metadata/services", bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Printf("Couldn't register service: %s", err)
	}
	req.WithContext(ctx)

	resp, err := ks.httpSecureClient.Do(req)
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
		fmt.Fprintf(w, "Status: %d >%s< \n on URL: %s", resp.StatusCode, bodyString, "https://gateway.kyma.local:31635/github-test/v1/metadata/services")
	}
}
