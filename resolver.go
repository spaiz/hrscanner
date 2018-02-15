package main

import (
	"bufio"
	"errors"
	"github.com/miekg/dns"
	"math/rand"
	"net"
	"os"
	"time"
)

var errEmptyResults = errors.New("empty results returned from domain lookup")

// DNSResolver defines interface to able to use
// different DNS resolver implementations
type DNSResolver interface {
	Resolve(host string) ([]net.IP, error)
}

// NewDNSResolver returns MyDNSResolver instance
func NewDNSResolver() *MyDNSResolver {
	return &MyDNSResolver{
		random:    rand.New(rand.NewSource(time.Now().UnixNano())),
		ips:       make([]string, 0),
		dnsClient: &dns.Client{},
	}
}

// MyDNSResolver loads DNS servers from the file
// and return resolved IP by randomly selecting DNS server from the list
type MyDNSResolver struct {
	ips       []string
	random    *rand.Rand
	dnsClient *dns.Client
}

// Load loads ips from the file
func (r *MyDNSResolver) Load(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}

	defer f.Close()
	sc := bufio.NewScanner(f)

	for sc.Scan() {
		r.ips = append(r.ips, net.JoinHostPort(sc.Text(), "53"))
	}

	return sc.Err()
}

// Resolve returns resolved IP for specific host
func (r *MyDNSResolver) Resolve(host string) ([]net.IP, error) {
	var result []net.IP

	m1 := &dns.Msg{}
	m1.SetQuestion(dns.Fqdn(host), dns.TypeA)

	msg, _, err := r.dnsClient.Exchange(m1, r.getServer())
	if err != nil {
		return result, err
	}

	if msg != nil && msg.Rcode != dns.RcodeSuccess {
		return result, errors.New(dns.RcodeToString[msg.Rcode])
	}

	for _, record := range msg.Answer {
		if t, ok := record.(*dns.A); ok {
			result = append(result, t.A)
		}
	}

	if len(result) == 0 {
		return result, errEmptyResults
	}

	return result, nil
}

// return random DNS server
func (r *MyDNSResolver) getServer() string {
	return r.ips[r.random.Intn(len(r.ips))]
}
