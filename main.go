package main

import (
	"crypto/tls"
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

type CertPem struct {
	Cert *x509.Certificate
	Pem  []byte
}

func (certPem *CertPem) GetCertPool() *x509.CertPool {
	pool := x509.NewCertPool()

	ok := pool.AppendCertsFromPEM(certPem.Pem)

	if !ok {
		log.Fatal("Couldn't add pem")
	}

	return pool
}

func NewCertPem(pem []byte) CertPem {
	certPem := CertPem{
		Pem: pem,
	}

	certPem.ParseCert()

	return certPem
}

func (certPem *CertPem) ParseCert() {
	block, _ := pem.Decode(certPem.Pem)

	cert, err := x509.ParseCertificate(block.Bytes)

	if err != nil {
		log.Fatal(err)
	}

	certPem.Cert = cert
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

func filterExpiringCerts(certs []CertPem, expire time.Time) []*x509.Certificate {
	output := make([]*x509.Certificate, 0, len(certs))

	for _, cert := range certs {

		verifyOptions := x509.VerifyOptions{
			CurrentTime:   expire,
			Intermediates: cert.GetCertPool(),
		}

		if _, err := cert.Cert.Verify(verifyOptions); err != nil {
			output = append(output, cert.Cert)
		}
	}

	return output
}

func getCertificates(dirs []string) []CertPem {
	certificates := make([]CertPem, len(dirs))

	for index, dir := range dirs {
		certPath := path.Join(certsRoot, dir)
		pem := readPem(certPath)

		certPem := NewCertPem(pem)

		certificates[index] = certPem
	}

	return certificates
}

type ExpiringCert struct {
	Domains []string
	Expire  time.Time
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
				fmt.Printf("%s\t%s", domain, cert.Expire.Format(time.RFC3339))
			} else {
				fmt.Printf("%s", domain)
			}
			fmt.Println()
		}
	}
}

func downloadCert(domain string) ([]*x509.Certificate, error) {

	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", domain, 443), conf)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return conn.ConnectionState().PeerCertificates, nil
}

func checkFileCerts(certsRoot string, expire time.Time) []*x509.Certificate {

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

	flag.Parse()
}

func checkDomainCerts(domainList []string, expire time.Time) []ExpiringCert {
	var expiredDomains []ExpiringCert
	for _, domain := range domainList {
		certs, err := downloadCert(domain)
		if err != nil {
			log.Fatal(err)
		}

		for _, cert := range certs {
			// fmt.Println("Checking domain: ", cert.DNSNames, cert.NotAfter, expire)
			if cert.NotAfter.Before(expire) {
				expiredDomains = append(expiredDomains, ExpiringCert{
					Domains: cert.DNSNames,
					Expire:  cert.NotAfter,
				})
			}

		}
	}
	return expiredDomains
}

func getExpiredDomains(expire time.Time) []ExpiringCert {
	if domains == "" {
		expiringCerts := checkFileCerts(certsRoot, expire)
		return collectExpirations(expiringCerts)
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

	expiredDomains := getExpiredDomains(expire)
	printExpiringCerts(expiredDomains, printDate)
}
