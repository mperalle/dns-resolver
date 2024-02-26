package main

import (
	"fmt"
	"os"

	"github.com/miekg/dns"
)

func main() {

	// How to do a simple query to the rootserver ?

	const hostname = "google.com."

	// in the list of iana:
	rootServerIp := "198.41.0.4:53"

	m := new(dns.Msg)
	fmt.Println("Creating  message...")
	m.SetQuestion(hostname, dns.TypeA)

	c := new(dns.Client)
	fmt.Println("Querying the root server at", rootServerIp, "for,", hostname)
	in, _, err := c.Exchange(m, rootServerIp)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Ok received from", rootServerIp)

	// fmt.Println("Question section: ", in.Question)
	fmt.Println("Answer section: ", in.Answer)
	fmt.Println("Authoritative section: ", in.Ns)
	//fmt.Println("RTT: ", rtt)

	// fmt.Println(in.Extra[0].Header().Name)
	// fmt.Println(in.Extra[0].Header().Class)
	// fmt.Println(in.Extra[0].Header().Rrtype)

	fmt.Println("So there is only NS records for .com domain in the authoritative section, let's check in the extra section for their IP adresses")
	fmt.Println("Extra section: ", in.Extra)

	// Careful use a pointer to dns.A because method Header has a pointer receiver
	fmt.Println("Let's extract the IP adress of the first one,", in.Extra[0].Header().Name)
	fmt.Println(in.Extra[0].(*dns.A).A)

	fmt.Println("Ok, let's query to this server for", hostname)

	tldServerIp := in.Extra[0].(*dns.A).A

	m = new(dns.Msg)
	fmt.Println("Creating message...")
	m.SetQuestion(hostname, dns.TypeA)

	c = new(dns.Client)
	fmt.Println("Querying the root server at", tldServerIp, "for,", hostname)
	in, _, err = c.Exchange(m, tldServerIp.String()+":53")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Ok received from", tldServerIp)
	fmt.Println("Answer section: ", in.Answer)
	fmt.Println("Authoritative section: ", in.Ns)

	fmt.Println("Ok same here no type A record in the authoritative section, let's check in the extra section")
	fmt.Println("Extra section: ", in.Extra)

	fmt.Println("Let's extract the IP adress of the first one of type A,", in.Extra[0].Header().Name)
	fmt.Println(in.Extra[1].(*dns.A).A)

	fmt.Println("Ok, let's query to this server for", hostname)

	nameServerIp := in.Extra[1].(*dns.A).A

	m = new(dns.Msg)
	fmt.Println("Creating message...")
	m.SetQuestion(hostname, dns.TypeA)

	c = new(dns.Client)
	fmt.Println("Querying the root server at", nameServerIp, "for,", hostname)
	in, _, err = c.Exchange(m, nameServerIp.String()+":53")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Ok received from", nameServerIp)
	fmt.Println("Answer section: ", in.Answer)
	fmt.Println("Authoritative section: ", in.Ns)
	fmt.Println("Extra section: ", in.Extra)

	fmt.Println("Ok we have now an answer in the Answer section with a record of type A")
	fmt.Println("So the IP adress for", hostname, "is:", in.Answer[0].(*dns.A).A)

}
