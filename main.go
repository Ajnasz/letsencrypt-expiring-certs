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
	"time"
)

var pemName string
var certsRoot string
var expireTime string

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

func getPem(dir string) *pem.Block {
	f, err := ioutil.ReadFile(path.Join(dir, pemName))

	if err != nil {
		log.Fatal(err)
	}

	block, _ := pem.Decode(f)

	if block == nil {
		log.Fatal("Failed parse cert pem")
	}

	return block
}

func parseCert(bytes []byte) *x509.Certificate {
	cert, err := x509.ParseCertificate(bytes)

	if err != nil {
		log.Fatal(err)
	}

	return cert
}

func filterExpiringCerts(certs []*x509.Certificate, expire time.Time) []*x509.Certificate {
	var output []*x509.Certificate

	for _, cert := range certs {
		if cert.NotAfter.Before(expire) {
			output = append(output, cert)
		}
	}

	return output
}

func init() {
	flag.StringVar(&pemName, "pem-name", "fullchain.pem", "The name of the pem file, usually fullchain.pem")
	flag.StringVar(&certsRoot, "certs-path", "/etc/letsencrypt/live", "The path to the directory which stores the certificates")
	flag.StringVar(&expireTime, "expire", "", "Expire time of the certificates in unix date format (run date command \"$(date --date='15/03/2016')\"), eg.: Mon Dec 14 13:36:37 CET 2015")

	flag.Parse()
}

func main() {
	var expire time.Time

	if expireTime == "" {
		expire = time.Now().Add(time.Hour * 24 * 7 * 2)
	} else {
		var err error
		expire, err = time.Parse(time.UnixDate, expireTime)

		if err != nil {
			log.Fatal(err)
		}
	}

	dirs := getCertDirectoryNames(certsRoot)

	var certikek []*x509.Certificate

	for _, dir := range dirs {
		certPath := path.Join(certsRoot, dir)
		cert := parseCert(getPem(certPath).Bytes)

		certikek = append(certikek, cert)
	}

	expiringCerts := filterExpiringCerts(certikek, expire)

	for _, cert := range expiringCerts {
		for _, domain := range cert.DNSNames {
			fmt.Println(domain)
		}
	}
}
