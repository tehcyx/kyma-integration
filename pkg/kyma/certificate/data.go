package certificate

// CertConnectResponse response the app gets from Kyma, when CSR is successful.
type CertConnectResponse struct {
	Cert string `json:"crt,omitempty"`
}

// ApplicationConnectResponse response the app gets from Kyma, when connect is successful.
type ApplicationConnectResponse struct {
	CsrURL      string  `json:"csrUrl,omitempty"`
	API         APIData `json:"api,omitempty"`
	Certificate Data    `json:"certificate,omitempty"`
}

// APIData response holding the metadata url that is important for service registry.
type APIData struct {
	MetadataURL     string `json:"metadataUrl,omitempty"`
	EventsURL       string `json:"eventsUrl,omitempty"`
	CertificatesURL string `json:"certificatesUrl,omitempty"`
}

// Data part of ApplicationConnectResponse holding key parameters for cert generation.
type Data struct {
	Subject      string `json:"subject,omitempty"`
	Extensions   string `json:"extensions,omitempty"`
	KeyAlgorithm string `json:"key-algorithm,omitempty"`
}

// CACertificate app struct to hold information about keys used to connect to Kyma.
type CACertificate struct {
	PrivateKey string
	PublicKey  string
	Csr        string
	ServerCert string
}
