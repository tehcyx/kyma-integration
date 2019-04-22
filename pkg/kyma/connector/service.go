package connector

import "encoding/json"

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
