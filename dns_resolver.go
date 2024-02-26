package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/miekg/dns"
)

func dnsQuery(hostname string, ipAdress string) *dns.Msg {

	m := new(dns.Msg)
	fmt.Println("Creating  message...")
	m.SetQuestion(hostname, dns.TypeA)

	c := new(dns.Client)
	fmt.Println("Querying the server at", ipAdress, "for", hostname)
	response, _, err := c.Exchange(m, ipAdress)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Ok received from", ipAdress)
	fmt.Println("Answer section: ", response.Answer)
	fmt.Println("Authoritative section: ", response.Ns)
	fmt.Println("Extra section: ", response.Extra)

	return response
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

	var serverIp string
	// We start with the first root server of the list of iana:
	serverIp = "198.41.0.4:53"

	// Query the root server
	response := dnsQuery(hostname, serverIp)

	for _, resourceRecord := range response.Extra {
		if resourceRecord.Header().Rrtype == dns.TypeA {
			serverIp = resourceRecord.(*dns.A).A.String() + ":53"
		}
	}

	// Query the TLD nameserver
	response = dnsQuery(hostname, serverIp)

	for _, resourceRecord := range response.Extra {
		if resourceRecord.Header().Rrtype == dns.TypeA {
			serverIp = resourceRecord.(*dns.A).A.String() + ":53"
		}
	}

	// Query the final nameserver
	response = dnsQuery(hostname, serverIp)

	for _, resourceRecord := range response.Answer {
		if resourceRecord.Header().Rrtype == dns.TypeA {
			serverIp = resourceRecord.(*dns.A).A.String()
		}
	}

	fmt.Println("The IP adress for", hostname, "is:", serverIp)

}
