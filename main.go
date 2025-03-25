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

func getDefaultExpireTime() time.Time {
	return time.Now().Add(time.Hour * 24 * 7 * 2)
}

var dateFormats = []string{
	time.UnixDate,
	time.RFC3339,
	time.RFC3339Nano,
	time.RFC1123,
	time.RFC1123Z,
	time.RFC822,
	time.RFC822Z,
	time.RFC850,
	time.RFC1123Z,
}

func getUserDefinedExpireTime(expireTime string) (time.Time, error) {
	for _, format := range dateFormats {
		expire, err := time.Parse(format, expireTime)
		if err == nil {
			return expire, nil
		}
	}

	return time.Time{}, nil
}

func getExpireTime(expireTime string) (time.Time, error) {
	if expireTime == "" {
		return getDefaultExpireTime(), nil
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

func printExpiringCerts(expiringCerts []ExpiringCert) {
	for _, cert := range expiringCerts {
		for _, domain := range cert.Domains {
			fmt.Printf("%s\t%s", domain, cert.Expire.Format(time.RFC3339))
			fmt.Println()
		}
	}
}

func init() {
	flag.StringVar(&pemName, "pem-name", "fullchain.pem", "The name of the pem file, usually fullchain.pem")
	flag.StringVar(&certsRoot, "certs-path", "/etc/letsencrypt/live", "The path to the directory which stores the certificates")
	flag.StringVar(&expireTime, "expire", getDefaultExpireTime().Format(time.RFC3339), "Expire time of the certificates (run date command \"$(date -Ih --date='03/15/2016')\"), eg.: 2016-03-15T00+01:00")

	flag.Parse()
}

func main() {

	expire, err := getExpireTime(expireTime)

	if err != nil {
		log.Fatal(err)
	}

	dirs := getCertDirectoryNames(certsRoot)

	certificates := getCertificates(dirs)

	expiringCerts := filterExpiringCerts(certificates, expire)

	domains := collectExpirations(expiringCerts)

	printExpiringCerts(domains)
}
