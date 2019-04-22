package certificate

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
	PrivateKey string
	PublicKey  string
	Csr        string

	PrivateKeyPath string
	PublicKeyPath  string
	CsrPath        string

	ServerCertPath string
}
