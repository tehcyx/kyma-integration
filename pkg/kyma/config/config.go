package config

import (
	"crypto/x509/pkix"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/tehcyx/kyma-integration/pkg/kyma/certificate"
	"gopkg.in/yaml.v3"
)

// Config holds the config for the app
type Config struct {
	App struct {
		ID     string `yaml:"id"`
		Name   string `yaml:"name"`
		Remote string `yaml:"remote"`
	} `yaml:"app"`
	CSR        string `yaml:"request.csr"`
	PublicKey  string `yaml:"client.crt"`
	PrivateKey string `yaml:"client.key"`
	ServerCert string `yaml:"server.crt"`
}

// CertExists certificate exists and can be used
func (cfg *Config) CertExists() bool {
	if cfg.CSR != "" && cfg.PublicKey != "" && cfg.PrivateKey != "" && cfg.ServerCert != "" {
		return true
	}
	return false
}

// GenerateKeysAndCertificate generates keys and certificates
func (cfg *Config) GenerateKeysAndCertificate(subject string) *certificate.CACertificate {
	var appCert *certificate.CACertificate
	appCert = new(certificate.CACertificate)

	if cfg.CertExists() {
		// read cert.csr
		appCert.Csr = cfg.CSR
		appCert.PublicKey = cfg.PublicKey
		appCert.PrivateKey = cfg.PrivateKey
	} else {
		location := "Palo Alto"
		province := "CA"
		country := "USA"
		organization := "Organization"
		organizationalUnit := "OrgUnit"
		if cfg.App.Name == "" {
			cfg.App.Name = "api-test"
		}
		commonName := cfg.App.Name

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
		// override config
		cfg.CSR = genCert.Csr
		cfg.PublicKey = genCert.PublicKey
		cfg.PrivateKey = genCert.PrivateKey

		// safe in appCert
		appCert = genCert
	}

	return appCert
}

func readConfig() Config {
	data, err := ioutil.ReadFile(getConfigPath())
	if err != nil {
		log.Printf("reading file resulted in an error: %w", err)
	}
	cfg := Config{}
	yaml.Unmarshal(data, cfg)
	return cfg
}

func getConfigPath() string {
	envDir := os.Getenv("CONFIG_DIR")
	if envDir == "" {
		// get current application directory
		currentDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Fatal(err)
		}
		envDir = currentDir
	}
	log.Printf("%s", envDir)
	configFilePath := filepath.Join(envDir, "config.yaml")
	return configFilePath
}

func (cfg *Config) saveConfig() {
	f, err := os.Create(getConfigPath())
	if err != nil {
		log.Printf("error creating file: %w", err)
	}
	defer f.Close()
	data, err := yaml.Marshal(cfg)
	if err != nil {
		log.Printf("couldn't transform: %w", err)
	}
	_, writeErr := f.Write(data)
	if writeErr != nil {
		log.Printf("write failed: %w", err)
	}
	f.Close()
}

// UpdateAppID updates the app config with the App ID received from remote and calls saveConfig().
func (cfg *Config) UpdateAppID(id string) {
	cfg.App.ID = id
	cfg.saveConfig()
}

// UpdateServerCert updates the app config with the certificate received from remote and calls saveConfig().
func (cfg *Config) UpdateServerCert(certData string) {
	cfg.ServerCert = certData
	cfg.saveConfig()
}

// New reads or creates a new config. If a config file is in place it will always return the files contents
func New() Config {
	configFilePath := getConfigPath()
	cfg := Config{}
	if _, err := os.Stat(configFilePath); err != nil {
		cfg.saveConfig()
		return cfg
	}
	cfg = readConfig()
	return cfg
}
