package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"time"
)

type DomainCert struct {
	Domain string
	Certs  []*x509.Certificate
}

func NewDomainCert(domain string, certs []*x509.Certificate) DomainCert {
	return DomainCert{
		Domain: domain,
		Certs:  certs,
	}
}

var _ CertVerifier = (*DomainCert)(nil)

func (domainCert DomainCert) Verify(expire time.Time) error {
	verifyOptions := x509.VerifyOptions{
		CurrentTime:   expire,
		Intermediates: x509.NewCertPool(),
		Roots:         x509.NewCertPool(),
	}

	for _, cert := range domainCert.Certs {
		verifyOptions.Intermediates.AddCert(cert)
	}

	systemRoots, err := x509.SystemCertPool()
	if err == nil {
		verifyOptions.Roots = systemRoots
	}

	for _, cert := range domainCert.Certs {
		_, err = cert.Verify(verifyOptions)
		if err != nil {
			return err
		}
	}

	return nil
}

func (domainCert DomainCert) GetCert() *x509.Certificate {
	return domainCert.Certs[0]
}

func downloadCert(domain string) ([]*x509.Certificate, error) {
	host, port, err := net.SplitHostPort(domain)

	if err != nil {
		host = domain
		port = "443"
	}

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%s", host, port), &tls.Config{})
	if err != nil {
		return nil, err
	}

	defer conn.Close()
	return conn.ConnectionState().PeerCertificates, nil
}

func checkDomainCerts(domainList []string, expire time.Time) ([]ExpiringCert, error) {
	var expiredDomains []ExpiringCert
	for _, domain := range domainList {
		certs, err := downloadCert(domain)
		if err != nil {
			return nil, err
		}

		domainCert := NewDomainCert(domain, certs)

		// fmt.Println("Checking domain: ", cert.DNSNames, cert.NotAfter, expire)
		if err := domainCert.Verify(expire); err != nil {
			expiredDomains = append(expiredDomains, ExpiringCert{
				Domains: certs[0].DNSNames,
				Expire:  certs[0].NotAfter,
				Error:   err,
			})
		}

	}
	return expiredDomains, nil
}
