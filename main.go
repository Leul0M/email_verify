package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {

	Scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("domain, hasMx, hasSPF, spfRecords,hasDMARC, dmarcRecords,validation\n")
	for Scanner.Scan() {
		checkDomain(Scanner.Text())
	}
	if err := Scanner.Err(); err != nil {
		log.Fatalf("Error: could not read from the input: %v\n", err)
	}
}

func checkDomain(domain string) {
	var hasMx, hasSPF, hasDMARC bool
	var spfRecord, dmarcRecord string

	mxrecords, err := net.LookupMX(domain)
	if err != nil {
		log.Printf("MX Lookup Error: %v\n", err)
	}
	if len(mxrecords) > 0 {
		hasMx = true
	}
	txtrecords, err := net.LookupTXT(domain)
	if err != nil {
		log.Printf("TXT Lookup Error: %v\n", err)
	}
	for _, record := range txtrecords {
		if strings.HasPrefix(record, "v=spf1") {
			hasSPF = true
			spfRecord = record
			break
		}
	}
	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		log.Printf("DMARC Lookup Error: %v\n", err)
	}
	for _, record := range dmarcRecords {
		if strings.HasPrefix(record, "v=DMARC1") {
			hasDMARC = true
			dmarcRecord = record
			break
		}
	}
	fmt.Printf("%s, %t, %t, %s, %t, %s, %t\n", domain, hasMx, hasSPF, spfRecord, hasDMARC, dmarcRecord, hasMx && hasSPF && hasDMARC)
}
