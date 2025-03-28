package main

import (
	"crypto/x509"
	"testing"
)

func TestCollectDomains(t *testing.T) {
	var data []*x509.Certificate

	data = append(data, &x509.Certificate{
		DNSNames: []string{"foo.bar", "foo.baz"},
	})
	data = append(data, &x509.Certificate{
		DNSNames: []string{"foo.qux", "foo.norf"},
	})

	collectedDomains := collectExpirations(data)

	actual := len(collectedDomains)
	expected := len(data)
	if actual != expected {
		t.Fatal("Domains length is different")
	}

	for index, expiringCert := range collectedDomains {
		expected := len(data[index].DNSNames)
		actual := len(expiringCert.Domains)

		if actual != expected {
			t.Fatal("Collected domains length is different than cert.DNSNames")
		}
	}
}
