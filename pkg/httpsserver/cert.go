package httpsserver

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	log "github.com/Kran001/basic-auth/pkg/logging"
	"github.com/kabukky/httpscerts"
)

type Cert struct {
	CaCert        []byte
	CaCertPool    *x509.CertPool
	CaCertificate tls.Certificate
}

func LoadCert(certFile, keyFile string) (*Cert, error) {
	var err error
	resultCert := new(Cert)
	resultCert.CaCertificate, err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Logger.Error("Failed loading PEM files")
		return resultCert, err
	}

	resultCert.CaCert, err = ioutil.ReadFile(certFile)

	if err != nil {
		log.Logger.Error("Failed loading certificate")
		return resultCert, err
	}

	resultCert.CaCertPool = x509.NewCertPool()
	resultCert.CaCertPool.AppendCertsFromPEM(resultCert.CaCert)

	return resultCert, nil
}

func (c *Cert) CreateTLSConfig() *tls.Config {
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{c.CaCertificate},
		RootCAs:      c.CaCertPool,
	}

	return tlsConfig
}

func CheckCerts(certPath, keyPath, host string) error {
	err := httpscerts.Check(certPath, keyPath)
	if err != nil {
		return httpscerts.Generate(certPath, keyPath, host)
	}

	return nil
}
