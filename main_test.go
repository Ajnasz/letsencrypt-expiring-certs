package main

import (
	"crypto/x509"
	"testing"
	"time"
)

func TestCollectDomains(t *testing.T) {
	var data []*x509.Certificate

	data = append(data, &x509.Certificate{
		DNSNames: []string{"foo.bar", "foo.baz"},
	})
	data = append(data, &x509.Certificate{
		DNSNames: []string{"foo.qux", "foo.norf"},
	})

	collectedDomains := collectDomains(data)

	actual := len(collectedDomains)
	expected := len(data)
	if actual != expected {
		t.Fatal("Domains length is different")
	}

	for index, domainList := range collectedDomains {
		expected := len(data[index].DNSNames)
		actual := len(domainList)

		if actual != expected {
			t.Fatal("Collected domains length is different than cert.DNSNames")
		}
	}
}

func TestGetDefaultExpireTime(t *testing.T) {
	actual := getDefaultExpireTime()

	now := time.Now()

	actual = actual.Truncate(time.Hour)

	expected := now.Truncate(time.Hour).AddDate(0, 0, 14)

	if actual != expected {
		t.Fatal("Default expire time should be 2 weeks", actual, expected)
	}
}
