package main

import (
	"fmt"
	"log"
	"os"
	"strings"

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

func hasAnswer(response *dns.Msg) string {
	for _, record := range response.Answer {
		if record.Header().Rrtype == dns.TypeA {
			fmt.Println("Found A record from the final nameserver!")
			return record.(*dns.A).A.String()
		}
	}
	fmt.Println("No A record in Answer section...")
	return ""
}

func hasExtra(response *dns.Msg) string {
	for _, record := range response.Extra {
		if record.Header().Rrtype == dns.TypeA {
			fmt.Println("Found A record in Additional section, let's query its server...")
			return record.(*dns.A).A.String() + ":53"
		}
	}
	fmt.Println("No A record in Additional section...")
	return ""
}

func hasNs(response *dns.Msg) string {
	for _, record := range response.Ns {
		if record.Header().Rrtype == dns.TypeNS {
			fmt.Println("Found domain name for NS record in Authoritative section, let's resolve it...")
			return record.(*dns.NS).Ns
		}
	}
	fmt.Println("No NS record in Authoritative section...")
	return ""
}

func resolve(hostname string) string {

	// The first server to query is one of the 13 root servers
	serverIp := "198.41.0.4:53"

	for {
		response := dnsQuery(hostname, serverIp)
		// Check in the Anwser section
		if ipAdress := hasAnswer(response); ipAdress != "" {
			return ipAdress
			// Check in the Additional section
		} else if ipAdress := hasExtra(response); ipAdress != "" {
			serverIp = ipAdress
			continue
			// Check in the Authoritative section
		} else if domainName := hasNs(response); domainName != "" {
			serverIp = resolve(domainName) + ":53"
			continue
		}

		log.Fatal("Nothing found in the different sections...")

	}
}

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Usage: ./dsn_resolver hostname")
		return
	}

	hostname := os.Args[1]

	// Fully qualified hostname
	if !strings.HasSuffix(os.Args[1], ".") {
		hostname = os.Args[1] + "."
	}

	fmt.Println("The IP adress for", os.Args[1], "is:", resolve(hostname))

}
