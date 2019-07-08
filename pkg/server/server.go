package server

import (
	"context"
	"crypto/tls"
	"crypto/x509/pkix"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"time"

	"github.com/tehcyx/kyma-integration/internal/certificate"
	"github.com/tehcyx/kyma-integration/internal/handler"
)

// Server struct to handle http exposure and TLS
type Server struct {
	Context                    context.Context
	Host, Port, SecurePort     string
	Handlers                   handler.Param
	ListenerNoTLS, ListenerTLS net.Listener
	Client, SecureClient       *http.Client
	Certificate                *certificate.CACertificate
	TLSPath                    string
	AppName                    string
}

// New creates a new server allowing you to expose rest endpoints
func New(host, port, securePort string, handlers handler.Param) *Server {
	for path, hndl := range handlers {
		http.HandleFunc(path, hndl)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	envDir := getTLSPath()
	return &Server{
		Host:          host,
		Port:          port,
		SecurePort:    securePort,
		Handlers:      handlers,
		ListenerNoTLS: nil,
		ListenerTLS:   nil,
		Client:        &http.Client{Transport: tr},
		SecureClient:  nil,
		TLSPath:       envDir,
		Certificate: &certificate.CACertificate{
			CsrPath:        path.Join(envDir, "request.csr"),
			PublicKeyPath:  path.Join(envDir, "client.crt"),
			PrivateKeyPath: path.Join(envDir, "client.key"),
			ServerCertPath: path.Join(envDir, "server.crt"),
		},
	}
}

// AddHandlers adds more handlers
func (srv *Server) AddHandlers(handlers handler.Param) error {
	for path := range handlers {
		if _, ok := srv.Handlers[path]; ok {
			return fmt.Errorf("Handler with path '%s' already exists. Not adding any handlers from this call", path)
		}
	}
	for path, hndl := range handlers {
		http.HandleFunc(path, hndl)
		srv.Handlers[path] = hndl
	}
	return nil
}

// Run starts up the server, if the TLS certificates exist it also starts with TLS
func (srv *Server) Run() {
	if srv.CertExists() {
		srv.StartListenTLS()
	}
	srv.StartListen()
}

// CertExists certificate exists and can be used
func (srv *Server) CertExists() bool {
	_, errCert := os.Stat(srv.Certificate.ServerCertPath)
	if errCert == nil {
		return true
	} else if os.IsNotExist(errCert) {
		return false
	} else {
		log.Fatal("read error on cert file")
		return false
	}
}

// GenerateKeysAndCertificate generates keys and certificates
func (srv *Server) GenerateKeysAndCertificate(subject string) *certificate.CACertificate {
	var appCert *certificate.CACertificate
	appCert = new(certificate.CACertificate)

	appCert.CsrPath = srv.Certificate.CsrPath
	appCert.PrivateKeyPath = srv.Certificate.PrivateKeyPath
	appCert.PublicKeyPath = srv.Certificate.PublicKeyPath
	appCert.ServerCertPath = srv.Certificate.ServerCertPath

	if !srv.CertExists() {
		_, errCSR := os.Stat(srv.Certificate.CsrPath)
		_, errPub := os.Stat(srv.Certificate.PublicKeyPath)
		_, errPriv := os.Stat(srv.Certificate.PrivateKeyPath)

		// read cert.csr
		if errCSR == nil && errPub == nil && errPriv == nil {
			csrBytes, err := ioutil.ReadFile(srv.Certificate.CsrPath)
			if err != nil {
				log.Fatal("Read error on csr file")
			}
			appCert.Csr = string(csrBytes[:])
			pubKeyBytes, err := ioutil.ReadFile(srv.Certificate.PublicKeyPath)
			if err != nil {
				log.Fatal("Read error on pub file")
			}
			appCert.PublicKey = string(pubKeyBytes[:])
			privKeyBytes, err := ioutil.ReadFile(srv.Certificate.PrivateKeyPath)
			if err != nil {
				log.Fatal("Read error on priv file")
			}
			appCert.PrivateKey = string(privKeyBytes[:])
		} else if os.IsNotExist(errCSR) && os.IsNotExist(errPub) && os.IsNotExist(errPriv) {
			location := "Walldorf"
			province := "Walldorf"
			country := "DE"
			organization := "Organization"
			organizationalUnit := "OrgUnit"
			commonName := "api-test"

			if subject != "" {
				//TODO: add a more generic version of this, as it panics if the order of the elements in the subject line is changed
				subjectMatch := regexp.MustCompile("^O=(?P<o>.*),OU=(?P<ou>.*),L=(?P<l>.*),ST=(?P<st>.*),C=(?P<c>.*),CN=(?P<cn>.*)$")
				match := subjectMatch.FindStringSubmatch(subject)
				result := make(map[string]string)
				for i, name := range subjectMatch.SubexpNames() {
					if i != 0 && name != "" {
						result[name] = match[i]
					}
				}
				location = result["l"]
				province = result["st"]
				country = result["c"]
				organization = result["o"]
				organizationalUnit = result["ou"]
				commonName = result["cn"]
				srv.AppName = commonName
			}

			subject := pkix.Name{
				Locality:           []string{location},
				Province:           []string{province},
				Country:            []string{country},
				Organization:       []string{organization},
				OrganizationalUnit: []string{organizationalUnit},
				CommonName:         commonName,
				// ??:              []string{"OU=OrgUnit,O=Organization,L=Waldorf,ST=Waldorf,C=DE,CN=api-test"},
			}

			genCert, err := certificate.GenerateCSR(subject, time.Duration(1200), 2048)
			if err != nil {
				fmt.Println(err)
			}
			//write files here
			csrBytes := []byte(genCert.Csr)
			pubKeyBytes := []byte(genCert.PublicKey)
			privKeyBytes := []byte(genCert.PrivateKey)
			errCSR := ioutil.WriteFile(srv.Certificate.CsrPath, csrBytes, 0644)
			if errCSR != nil {
				log.Fatal("couldn't write csr")
			}
			errPub := ioutil.WriteFile(srv.Certificate.PublicKeyPath, pubKeyBytes, 0644)
			if errPub != nil {
				log.Fatal("couldn't write pub key")
			}
			errPriv := ioutil.WriteFile(srv.Certificate.PrivateKeyPath, privKeyBytes, 0644)
			if errPriv != nil {
				log.Fatal("couldn't write priv key")
			}
			genCert.CsrPath = srv.Certificate.CsrPath
			genCert.PrivateKeyPath = srv.Certificate.PrivateKeyPath
			genCert.PublicKeyPath = srv.Certificate.PublicKeyPath
			genCert.ServerCertPath = srv.Certificate.ServerCertPath
			appCert = genCert
		} else {
			log.Fatal("cert not readable or does not exist")
		}
	}

	return appCert
}

// StartListen starts exposure of service on port of choice via http://
func (srv *Server) StartListen() {
	var err error
	srv.ListenerNoTLS, err = net.Listen("tcp", fmt.Sprintf("%s:%s", srv.Host, srv.Port))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("🔓 Listening on %s:%s\n", srv.Host, srv.Port)
	http.Serve(srv.ListenerNoTLS, nil)
}

// StartListenTLS starts exposure of service on port of choice via https://
func (srv *Server) StartListenTLS() {
	if srv.ListenerTLS != nil {
		log.Printf("Gracefully closing %s to restart with new certificate.", srv.SecurePort)
		srv.ListenerTLS.Close()
	}
	var err error
	srv.ListenerTLS, err = net.Listen("tcp", fmt.Sprintf("%s:%s", srv.Host, srv.SecurePort))
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		log.Println("updating http client with certificate")
		tr, err := certificate.CreateTLSConfig(srv.Certificate)
		if err != nil {
			log.Fatalf("TLSConfig could not be created: %s\n", err.Error())
		}
		srv.SecureClient = &http.Client{Transport: tr}
		log.Printf("🔐 Listening on %s:%s\n", srv.Host, srv.SecurePort)
		http.ServeTLS(srv.ListenerTLS, nil, srv.Certificate.ServerCertPath, srv.Certificate.PrivateKeyPath)
	}()
}

func getTLSPath() string {
	envDir := os.Getenv("KEY_DIR")
	if envDir == "" {
		// get current application directory
		currentDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Fatal(err)
		}
		envDir = path.Join(currentDir, ".key/")
		_ = os.MkdirAll(envDir, os.ModePerm)
	}
	return envDir
}
