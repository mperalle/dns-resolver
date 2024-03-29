package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/miekg/dns"
)

type cachedRecord struct {
	record     dns.RR
	expiration time.Time
}

var cache map[string]cachedRecord

func dnsQuery(hostname string, ipAdress string) *dns.Msg {

	m := new(dns.Msg)
	m.SetQuestion(hostname, dns.TypeA)

	c := new(dns.Client)
	fmt.Println("\nQuerying the server at", ipAdress, "for", hostname)
	response, _, err := c.Exchange(m, ipAdress)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Ok answer received from", ipAdress)

	return response
}

func hasAnswer(response *dns.Msg) dns.RR {

	for _, record := range response.Answer {
		if record.Header().Rrtype == dns.TypeA {
			fmt.Println("Found A record from the nameserver!")
			return record
		}
	}
	fmt.Println("No A record in Answer section...")
	return nil
}

func hasExtra(response *dns.Msg) string {
	for _, record := range response.Extra {
		if record.Header().Rrtype == dns.TypeA {
			fmt.Println("Found A record in Additional section, let's query its nameserver:", record.Header().Name)
			return record.(*dns.A).A.String() + ":53"
		}
	}
	fmt.Println("No A record in Additional section...")
	return ""
}

func hasNs(response *dns.Msg) string {
	for _, record := range response.Ns {
		if record.Header().Rrtype == dns.TypeNS {
			fmt.Println("Found NS record in Authoritative section, let's resolve its domain name:", record.(*dns.NS).Ns)
			return record.(*dns.NS).Ns
		}
	}
	fmt.Println("No NS record in Authoritative section...")
	return ""
}

func resolve(hostname string) dns.RR {

	// The first server to query is one of the 13 root nameservers
	serverIp := "198.41.0.4:53"

	for {
		response := dnsQuery(hostname, serverIp)
		// Check in the Anwser section
		if record := hasAnswer(response); record != nil {
			return record
			// Check in the Additional section
		} else if ipAddress := hasExtra(response); ipAddress != "" {
			serverIp = ipAddress
			continue
			// Check in the Authoritative section
		} else if domainName := hasNs(response); domainName != "" {
			serverIp = resolve(domainName).(*dns.A).A.String() + ":53"
			continue
		}

		fmt.Println("Nothing found in the different sections... please check the format of the hostname")
		return nil

	}
}

func dnsHandler(w dns.ResponseWriter, r *dns.Msg) {
	fmt.Println("******************************************************************")
	fmt.Println("\nNew message received from", w.RemoteAddr().String()+":")
	fmt.Println(r)

	hostname := r.Question[0].Name

	var answer dns.RR

	if r, inCache := cache[hostname]; inCache && time.Now().Before(r.expiration) {
		answer = r.record
		fmt.Println("Found matching record in cache")

	} else if inCache && time.Now().After(r.expiration) {

		fmt.Println("Cached record expired! Removing it from the cache...")
		delete(cache, hostname)

		answer = resolve(hostname)

		// Add record and expiration time to the cache
		ttl := answer.Header().Ttl

		fmt.Println("The TTL is:", ttl)

		cache[hostname] = cachedRecord{
			record:     answer,
			expiration: time.Now().Add(time.Duration(ttl) * time.Second),
		}

		fmt.Println("The expiration time is:", cache[hostname].expiration)

		fmt.Println("New record added to cache", cache)

	} else {
		answer = resolve(hostname)

		// Add record and expiration time to the cache
		ttl := answer.Header().Ttl

		fmt.Println("The TTL is:", ttl)

		cache[hostname] = cachedRecord{
			record:     answer,
			expiration: time.Now().Add(time.Duration(ttl) * time.Second),
		}

		fmt.Println("The expiration time is:", cache[hostname].expiration)

		fmt.Println("New record added to cache", cache)

	}

	// Prepare response to client
	m := new(dns.Msg)
	m.Id = r.Id
	m.Answer = append(m.Answer, answer)
	w.WriteMsg(m)
	fmt.Println("\nSending back A record:")
	fmt.Println(m)
}

func main() {

	// Initialize cache instance
	cache = make(map[string]cachedRecord)

	// Start dns server
	fmt.Println("Listen and serve on localhost:3000")
	err := dns.ListenAndServe("localhost:3000", "udp", dns.HandlerFunc(dnsHandler))
	if err != nil {
		log.Fatal(err)
	}
}
