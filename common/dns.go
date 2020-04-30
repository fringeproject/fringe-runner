package common

import (
	"errors"
	"sort"

	"github.com/miekg/dns"
)

// LookupName returns IPv4 addresses from A records or error.
func LookupName(fqdn, serverAddr string) ([]string, error) {
	ips := []string{}

	m4 := &dns.Msg{}
	m4.SetQuestion(dns.Fqdn(fqdn), dns.TypeA)

	r, err := dns.Exchange(m4, serverAddr+":53")
	if err != nil {
		return ips, err
	}

	if len(r.Answer) == 0 {
		return ips, nil
	}

	for _, answer := range r.Answer {
		if a, ok := answer.(*dns.A); ok {
			ip := a.A.String()
			ips = append(ips, ip)
		}
	}

	// There is no A record in the Answer, maybe AAA ?
	if len(ips) == 0 {
		return ips, nil
	}
	sort.Strings(ips)

	return ips, err
}

// LookupName returns IPv4 addresses from A records or error.
func LookupCName(fqdn, serverAddr string) ([]string, error) {
	cnames := []string{}

	m := &dns.Msg{}
	m.SetQuestion(dns.Fqdn(fqdn), dns.TypeCNAME)

	r, err := dns.Exchange(m, serverAddr+":53")
	if err != nil {
		return cnames, err
	}

	if len(r.Answer) == 0 {
		return cnames, nil
	}

	for _, answer := range r.Answer {
		if cname, ok := answer.(*dns.CNAME); ok {
			cname := cname.Target
			cnames = append(cnames, cname[:len(cname)-1])
		}
	}

	if len(cnames) == 0 {
		return cnames, nil
	}
	sort.Strings(cnames)

	return cnames, err
}

// LookupNS returns the names servers for a domain.
func LookupNS(domain, serverAddr string) ([]string, error) {
	servers := []string{}
	m := &dns.Msg{}
	m.SetQuestion(dns.Fqdn(domain), dns.TypeNS)
	r, err := dns.Exchange(m, serverAddr+":53")
	if err != nil {
		return servers, err
	}

	if len(r.Answer) < 1 {
		return servers, errors.New("No DNS Answer")
	}

	for _, a := range r.Answer {
		if ns, ok := a.(*dns.NS); ok {
			servers = append(servers, ns.Ns)
		}
	}

	return servers, nil
}
