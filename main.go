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
	fmt.Println("==================================================")
	fmt.Println("          Email Domain Verifier (CLI)             ")
	fmt.Println("==================================================")
	fmt.Println("Enter domain names (one per line, press Ctrl+C to exit):")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		domain := strings.TrimSpace(scanner.Text())
		if domain == "" {
			continue
		}
		checkDomain(domain)
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error: could not read from input: %v\n", err)
	}
}

func checkDomain(domain string) {
	var hasMx, hasSPF, hasDMARC bool
	var spfRecord, dmarcRecord string
	var mxErr, spfErr, dmarcErr error

	// 1. MX Lookup
	mxrecords, err := net.LookupMX(domain)
	if err != nil {
		mxErr = err
	}
	if len(mxrecords) > 0 {
		hasMx = true
	}

	// 2. SPF Lookup
	txtrecords, err := net.LookupTXT(domain)
	if err != nil {
		spfErr = err
	}
	for _, record := range txtrecords {
		if strings.HasPrefix(record, "v=spf1") {
			hasSPF = true
			spfRecord = record
			break
		}
	}

	// 3. DMARC Lookup
	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		dmarcErr = err
	}
	for _, record := range dmarcRecords {
		if strings.HasPrefix(record, "v=DMARC1") {
			hasDMARC = true
			dmarcRecord = record
			break
		}
	}

	// Format output block
	fmt.Println("--------------------------------------------------")
	fmt.Printf("Domain: %s\n", domain)

	// MX Status
	if hasMx {
		fmt.Println("  MX Record:    ✔ YES")
	} else {
		if mxErr != nil {
			fmt.Printf("  MX Record:    ❌ NO (%v)\n", cleanDNSErr(mxErr))
		} else {
			fmt.Println("  MX Record:    ❌ NO")
		}
	}

	// SPF Status
	if hasSPF {
		fmt.Printf("  SPF Record:   ✔ YES (%s)\n", spfRecord)
	} else {
		if spfErr != nil {
			fmt.Printf("  SPF Record:   ❌ NO (%v)\n", cleanDNSErr(spfErr))
		} else {
			fmt.Println("  SPF Record:   ❌ NO")
		}
	}

	// DMARC Status
	if hasDMARC {
		fmt.Printf("  DMARC Record: ✔ YES (%s)\n", dmarcRecord)
	} else {
		if dmarcErr != nil {
			fmt.Printf("  DMARC Record: ❌ NO (%v)\n", cleanDNSErr(dmarcErr))
		} else {
			fmt.Println("  DMARC Record: ❌ NO")
		}
	}

	// Overall Verification
	if hasMx && hasSPF && hasDMARC {
		fmt.Println("  Verification: 🎉 PASS (Domain is fully configured for secure email)")
	} else {
		var missing []string
		if !hasMx {
			missing = append(missing, "MX")
		}
		if !hasSPF {
			missing = append(missing, "SPF")
		}
		if !hasDMARC {
			missing = append(missing, "DMARC")
		}
		fmt.Printf("  Verification: ⚠️  FAIL (Missing %s)\n", strings.Join(missing, ", "))
	}
	fmt.Println("--------------------------------------------------")
	fmt.Println()
}

// cleanDNSErr extracts a shorter, friendlier message from standard DNS errors
func cleanDNSErr(err error) string {
	if err == nil {
		return ""
	}
	errStr := err.Error()
	if strings.Contains(errStr, "no such host") {
		return "host not found"
	}
	if strings.Contains(errStr, "i/o timeout") {
		return "timeout"
	}
	return errStr
}
