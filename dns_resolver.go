package main

import (
	"fmt"
	"log"
	"net"
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

	// The first server to query is one of the 13 root nameservers
	serverIp := "198.41.0.4:53"

	for {
		response := dnsQuery(hostname, serverIp)
		// Check in the Anwser section
		if ipAddress := hasAnswer(response); ipAddress != "" {
			return ipAddress
			// Check in the Additional section
		} else if ipAddress := hasExtra(response); ipAddress != "" {
			serverIp = ipAddress
			continue
			// Check in the Authoritative section
		} else if domainName := hasNs(response); domainName != "" {
			serverIp = resolve(domainName) + ":53"
			continue
		}

		fmt.Println("Nothing found in the different sections... please check the format of the hostname")
		return ""

	}
}

func main() {

	c, err := net.ListenPacket("udp", "localhost:3000")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Starting udp server on localhost:3000")

	defer c.Close()

	for {
		buffer := make([]byte, 1024)
		n, addr, err := c.ReadFrom(buffer)
		if err != nil {
			fmt.Println("Error in reading:", err)
			continue
		}

		hostname := string(buffer[:n])
		fmt.Println("Query received:", hostname)
		// Fully qualified hostname
		if !strings.HasSuffix(hostname, ".") {
			hostname = hostname + "."
		}

		ipAddr := resolve(hostname)
		if ipAddr != "" {
			fmt.Println("Sending IP address back...")
			c.WriteTo([]byte(ipAddr), addr)
		}

	}

}
