package main

import (
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/Ajnasz/letsencrypt-expiring-certs/expiration_time"
)

var pemName string
var certsRoot string
var expireTime string
var printDate bool
var domains string
var quit bool

type CertVerifier interface {
	Verify(time.Time) error
	GetCert() *x509.Certificate
}

func NewCertPem(pem []byte) CertPem {
	certPem := CertPem{
		Pem: pem,
	}

	certPem.ParseCert()

	return certPem
}

type CertPem struct {
	Cert *x509.Certificate
	Pem  []byte
}

var _ CertVerifier = (*CertPem)(nil)

func (certPem CertPem) getIntermediateCertPool() (*x509.CertPool, error) {
	pool := x509.NewCertPool()

	ok := pool.AppendCertsFromPEM(certPem.Pem)

	if !ok {
		return nil, fmt.Errorf("Failed to append certificate to pool")
	}

	return pool, nil
}

func (certPem CertPem) Verify(expire time.Time) error {
	pool, err := certPem.getIntermediateCertPool()
	if err != nil {
		return err
	}
	options := x509.VerifyOptions{
		CurrentTime:   expire,
		Intermediates: pool,
	}

	_, err = certPem.Cert.Verify(options)

	return err
}

func (certPem *CertPem) ParseCert() {
	block, _ := pem.Decode(certPem.Pem)

	cert, err := x509.ParseCertificate(block.Bytes)

	if err != nil {
		log.Fatal(err)
	}

	certPem.Cert = cert
}

func (certPem CertPem) GetCert() *x509.Certificate {
	return certPem.Cert
}

func isDir(name string) bool {
	f, err := os.Stat(name)

	if err != nil {
		log.Fatal(err)
	}

	return f.IsDir()
}

func getCertDirectoryNames(dir string) []string {
	f, err := os.Open(dir)

	if err != nil {
		log.Fatal(err)
	}

	names, err := f.Readdirnames(0)

	if err != nil {
		log.Fatal(err)
	}

	var dirs []string

	for index, name := range names {
		if isDir(path.Join(dir, name)) {
			dirs = append(dirs, names[index])
		}
	}

	return dirs
}

func readPem(dir string) []byte {
	f, err := os.ReadFile(path.Join(dir, pemName))

	if err != nil {
		log.Fatal(err)
	}

	return f
}

func filterExpiringCerts(certs []CertPem, expire time.Time) ([]*x509.Certificate, error) {
	output := make([]*x509.Certificate, 0, len(certs))

	for _, cert := range certs {
		if err := cert.Verify(expire); err != nil {
			output = append(output, cert.GetCert())
		}
	}

	return output, nil
}

func getCertificates(dirs []string) []CertPem {
	certificates := make([]CertPem, len(dirs))

	for index, dir := range dirs {
		certPath := path.Join(certsRoot, dir)
		pem := readPem(certPath)

		var certPem CertPem = NewCertPem(pem)

		certificates[index] = certPem
	}

	return certificates
}

type ExpiringCert struct {
	Domains []string
	Expire  time.Time
	Error   error
}

func collectExpirations(expiringCerts []*x509.Certificate) []ExpiringCert {
	expires := make([]ExpiringCert, 0, len(expiringCerts))

	for _, cert := range expiringCerts {
		expiringCert := ExpiringCert{
			Domains: cert.DNSNames,
			Expire:  cert.NotAfter,
		}
		expires = append(expires, expiringCert)
	}

	return expires
}

func printDomains(domains [][]string) {
	for _, domain := range domains {
		fmt.Println(strings.Join(domain, " "))
	}
}

func printExpiringCerts(expiringCerts []ExpiringCert, printDate bool) {
	for _, cert := range expiringCerts {
		for _, domain := range cert.Domains {
			if printDate {
				fmt.Printf("%s\t%s\t%s", domain, cert.Expire.Format(time.RFC3339), cert.Error)
			} else {
				fmt.Printf("%s", domain)
			}
			fmt.Println()
		}
	}
}

func checkFileCerts(certsRoot string, expire time.Time) ([]*x509.Certificate, error) {

	dirs := getCertDirectoryNames(certsRoot)

	certificates := getCertificates(dirs)

	return filterExpiringCerts(certificates, expire)
}

func init() {
	flag.StringVar(&pemName, "pem-name", "fullchain.pem", "The name of the pem file, usually fullchain.pem")
	flag.StringVar(&certsRoot, "certs-path", "/etc/letsencrypt/live", "The path to the directory which stores the certificates")
	flag.StringVar(&expireTime, "expire", "", "Expire time of the certificates (run date command \"$(date -Im --date='03/15/2016')\"), eg.: 2016-03-15T00+01:00. If empty, 2 weeks from now will be used")
	flag.BoolVar(&printDate, "print-date", false, "Print the expiration date of the certificates")
	flag.StringVar(&domains, "domains", "", "Comma separated list of domains to check")
	flag.BoolVar(&quit, "quit", false, "Quit the program with a non 0 exit code after printing the expiring certificates, if there are any")

	flag.Parse()
}

func getExpired(expire time.Time) ([]ExpiringCert, error) {
	if domains == "" {
		expiringCerts, err := checkFileCerts(certsRoot, expire)
		if err != nil {
			return nil, err
		}

		return collectExpirations(expiringCerts), nil
	} else {
		domainList := strings.Split(domains, ",")

		return checkDomainCerts(domainList, expire)
	}
}

func main() {
	expire, err := expiration_time.GetExpireTime(expireTime)

	if err != nil {
		log.Fatal(err)
	}

	expiredDomains, err := getExpired(expire)

	if err != nil {
		log.Fatal(err)
	}
	printExpiringCerts(expiredDomains, printDate)

	if quit && len(expiredDomains) > 0 {
		os.Exit(1)
	}
}
