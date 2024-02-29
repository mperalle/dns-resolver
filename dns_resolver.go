package main

import (
	"fmt"
	"log"
	"os"

	"github.com/miekg/dns"
)

func dnsQuery(hostname string, ipAdress string) *dns.Msg {

	m := new(dns.Msg)
	//fmt.Println("Creating  message...")
	m.SetQuestion(hostname, dns.TypeA)

	c := new(dns.Client)
	fmt.Println("Querying the server at", ipAdress, "for", hostname)
	response, _, err := c.Exchange(m, ipAdress)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Ok answer received from", ipAdress)
	// fmt.Println("Answer section: ", response.Answer)
	// fmt.Println("Authoritative section: ", response.Ns)
	// fmt.Println("Extra section: ", response.Extra)

	return response
}

func hasAnswer(response *dns.Msg) dns.RR {

	if len(response.Answer) != 0 {
		for _, record := range response.Answer {
			if record.Header().Rrtype == dns.TypeA {
				fmt.Println("Found A record from the final nameserver!")
				return record
			}
		}
	}

	fmt.Println("No A record in Answer section...")
	return nil
}

func hasExtra(response *dns.Msg) string {

	if len(response.Extra) != 0 {
		for _, record := range response.Extra {
			if record.Header().Rrtype == dns.TypeA {
				fmt.Println("Found A record in Additional section, let's query its server...")
				return record.(*dns.A).A.String() + ":53"
			}
		}
	}

	fmt.Println("No A record in Additional section...")
	return ""
}

func hasNs(response *dns.Msg) string {

	if len(response.Ns) != 0 {
		for _, record := range response.Ns {
			if record.Header().Rrtype == dns.TypeNS {
				fmt.Println("Found domain name for NS record in Authoritative section, let's resolve it...")
				return record.(*dns.NS).Ns
			}
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
			serverIp = resolve(domainName).Header().Name + ":53"
			continue
		}

		fmt.Println("Nothing found in the different sections... please check the format of the hostname")
		return nil

	}
}

func dnsHandler(w dns.ResponseWriter, r *dns.Msg) {
	fmt.Println("Message received:", r)
	fmt.Println(r.MsgHdr)
	m := new(dns.Msg)
	m.Id = r.Id
	anwser := resolve(r.Question[0].Name)
	m.Answer = append(m.Answer, anwser)
	w.WriteMsg(m)
	fmt.Println("Answering back...")
	fmt.Println(m)
}

func main() {

	fmt.Println("Listen and serve on localhost:3000")
	err := dns.ListenAndServe("localhost:3000", "udp", dns.HandlerFunc(dnsHandler))
	if err != nil {
		log.Fatal(err)
	}
}
