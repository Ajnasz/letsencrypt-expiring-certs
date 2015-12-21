package main

import (
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

var pemName string
var certsRoot string
var expireTime string

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

func getCertDirectoryNames(dir string) []string {
	f, err := os.Open(dir)

	if err != nil {
		log.Fatal(err)
	}

	names, err := f.Readdirnames(0)

	if err != nil {
		log.Fatal(err)
	}

	return names
}

func readPem(dir string) []byte {
	f, err := ioutil.ReadFile(path.Join(dir, pemName))

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

func getDefaultExpireTime() time.Time {
	return time.Now().Add(time.Hour * 24 * 7 * 2)
}

func getUserDefinedExpireTime(expireTime string) time.Time {
	expire, err := time.Parse(time.UnixDate, expireTime)

	if err != nil {
		log.Fatal(err)
	}

	return expire
}

func getExpireTime(expireTime string) time.Time {
	if expireTime == "" {
		return getDefaultExpireTime()
	} else {
		return getUserDefinedExpireTime(expireTime)
	}
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

func collectDomains(expiringCerts []*x509.Certificate) [][]string {
	domains := make([][]string, 0, len(expiringCerts))

	for _, cert := range expiringCerts {
		domains = append(domains, cert.DNSNames)
	}

	return domains
}

func printDomains(domains [][]string) {
	for _, domain := range domains {
		fmt.Println(strings.Join(domain, " "))
	}
}

func init() {
	flag.StringVar(&pemName, "pem-name", "fullchain.pem", "The name of the pem file, usually fullchain.pem")
	flag.StringVar(&certsRoot, "certs-path", "/etc/letsencrypt/live", "The path to the directory which stores the certificates")
	flag.StringVar(&expireTime, "expire", "", "Expire time of the certificates in unix date format (run date command \"$(date --date='15/03/2016')\"), eg.: Mon Dec 14 13:36:37 CET 2015")

	flag.Parse()
}

func main() {

	expire := getExpireTime(expireTime)

	dirs := getCertDirectoryNames(certsRoot)

	certificates := getCertificates(dirs)

	expiringCerts := filterExpiringCerts(certificates, expire)

	domains := collectDomains(expiringCerts)

	printDomains(domains)
}
