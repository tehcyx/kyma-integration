package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/martian/log"
)

type Service struct {
	Provider         string               `json:"provider,omitempty"`
	Name             string               `json:"name,omitempty"`
	Description      string               `json:"description,omitempty"`
	ShortDescription string               `json:"shortDescription,omitempty"`
	Labels           ServiceLabel         `json:"labels,omitempty"`
	API              ServiceAPI           `json:"api,omitempty"`
	Events           ServiceEvent         `json:"events,omitempty"`
	Documentation    ServiceDocumentation `json:"documentation,omitempty"`
}

type ServiceLabel map[string]string

type ServiceAPI struct {
	TargetURL   string             `json:"targetUrl,omitempty"`
	Spec        map[string]string  `json:"spec,omitempty"`
	Credentials ServiceCredentials `json:"credentials,omitempty"`
}

type ServiceCredentials struct {
	Basic ServiceBasicCredentials `json:"basic,omitempty"`
	OAuth ServiceOAuthCredentials `json:"oauth,omitempty"`
}

type ServiceBasicCredentials struct {
	ClientID string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type ServiceOAuthCredentials struct {
	ClientID     string `json:"clientId,omitempty"`
	ClientSecret string `json:"clientSecret,omitempty"`
}

type ServiceEvent struct {
	Spec ServiceEventSpec `json:"spec,omitempty"`
}

type ServiceEventSpec struct {
	AsyncAPI  string                           `json:"asyncapi,omitempty"`
	Info      ServiceEventSpecInfo             `json:"info,omitempty"`
	BaseTopic string                           `json:"baseTopic,omitempty"`
	Topics    map[string]ServiceEventSpecTopic `json:"topics,omitempty"`
}

type ServiceEventSpecInfo struct {
	Title       string `json:"title,omitempty"`
	Version     string `json:"version,omitempty"`
	Description string `json:"description,omitempty"`
}

type ServiceEventSpecTopic map[string]ServiceEventSpecTopicDetail

type ServiceEventSpecTopicDetail struct {
	Summary string                 `json:"summary,omitempty"`
	Payload map[string]interface{} `json:"payload,omitempty"`
}

type ServiceDocumentation struct {
	DisplayName string                    `json:"displayName,omitempty"`
	Description string                    `json:"description,omitempty"`
	Type        string                    `json:"type,omitempty"`
	Tags        []string                  `json:"tags,omitempty"`
	Docs        []ServiceDocumentationDoc `json:"docs,omitempty"`
}

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

	serviceDescription.Documentation.DisplayName = "Test"
	serviceDescription.Documentation.Description = "test decsription"
	serviceDescription.Documentation.Tags = []string{"Tag1", "Tag2"}
	serviceDescription.Documentation.Type = "Test Type"
	serviceDescription.Documentation.Docs = []ServiceDocumentationDoc{ServiceDocumentationDoc{Title: "Test", Type: "Doc", Source: "google.com"}}

	serviceDescription.Description = "API Description"
	serviceDescription.ShortDescription = "API Short Description"

	serviceDescription.Provider = "Daniel"
	serviceDescription.Name = "Daniel's Service"
	serviceDescription.API.TargetURL = "https://localhost:8443"

	jsonBytes, err := json.Marshal(serviceDescription)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(jsonBytes))

	// acquire NodePort to modify URL locally
	// 31490
	// https://gateway.kyma.local:31490/github-test/v1/metadata/services
	ks.appInfo.API.MetadataURL = "https://gateway.kyma.local:31490/github-test/v1/metadata/services"

	req, err := http.NewRequest("POST", ks.appInfo.API.MetadataURL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Errorf("Couldn't register service: %s", err)
	}
	req.WithContext(ctx)

	resp, err := ks.httpClient.Do(req)
	if err != nil {
		log.Errorf("Couldn't register service: %s", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	if resp.StatusCode == http.StatusOK {
		log.Infof("Successfully registered service with ID %s", "id here")
		fmt.Fprintf(w, bodyString)
	}
	fmt.Fprintf(w, "Status: %d >%s< \n on URL: %s", resp.StatusCode, bodyString, ks.appInfo.API.MetadataURL)
}
