package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

type KymaIntegrationServer struct {
	cert                                            *CACertificate
	dir, csrPath, pubPath, privPath, serverCertPath string
	httpClient                                      *http.Client
	httpSecureClient                                *http.Client
	listenNoTLS, listenTLS                          net.Listener
	appInfo                                         *ApplicationConnectResponse
}

func NewKymaIntegrationServer() *KymaIntegrationServer {
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

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	return &KymaIntegrationServer{
		dir:            envDir,
		csrPath:        path.Join(envDir, "request.csr"),
		pubPath:        path.Join(envDir, "client.crt"),
		privPath:       path.Join(envDir, "client.key"),
		serverCertPath: path.Join(envDir, "server.crt"),
		httpClient:     &http.Client{Transport: tr},
	}
}

func (ks *KymaIntegrationServer) Start() {
	http.HandleFunc("/", ks.indexHandler)
	http.HandleFunc("/github_callback", ks.gitHubCallbackHandler)
	http.HandleFunc("/connect", ks.connectHandler)
	http.HandleFunc("/register-service", ks.registerServiceHandler)

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", "127.0.0.1", "8080"))
	ks.listenNoTLS = listener
	if err != nil {
		log.Fatal(err)
	}

	if ks.tlsCertExists() {
		ks.startListenTLS()
	}
	fmt.Println("Listening on 8080")
	http.Serve(ks.listenNoTLS, nil)
}

func (ks *KymaIntegrationServer) indexHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "hello world")
}

func (ks *KymaIntegrationServer) startListenTLS() {
	if ks.listenTLS != nil {
		fmt.Println("Gracefully closing 443 to restart with new certificate.")
		ks.listenTLS.Close()
	}
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", "127.0.0.1", "8443"))
	ks.listenTLS = listener
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		fmt.Println("updating http client with certificate")
		// Load client cert
		cert, err := tls.LoadX509KeyPair(ks.serverCertPath, ks.privPath)
		if err != nil {
			log.Fatal(err)
		}

		// Load CA cert
		caCert, err := ioutil.ReadFile(ks.serverCertPath)
		if err != nil {
			log.Fatal(err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
			Certificates:       []tls.Certificate{cert},
			RootCAs:            caCertPool,
		}
		tlsConfig.BuildNameToCertificate()

		tr := &http.Transport{
			TLSClientConfig: tlsConfig,
		}
		ks.httpSecureClient = &http.Client{Transport: tr}
		fmt.Println("Listening on 8443")
		http.ServeTLS(ks.listenTLS, nil, ks.serverCertPath, ks.privPath)
	}()
}
