package main

import (
	"fmt"
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

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Usage: ./lookup hostname")
		return
	}

	hostname := os.Args[1]
	// Fully qualified hostname
	if !strings.HasSuffix(hostname, ".") {
		hostname = hostname + "."
	}

	response := dnsQuery(hostname, "localhost:3000")

	fmt.Println("The IP address of", strings.TrimSuffix(hostname, "."), "is:", response.Answer[0].(*dns.A).A.String())

}
