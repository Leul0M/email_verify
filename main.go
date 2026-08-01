package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
)

type VerificationResult struct {
	Domain      string `json:"domain"`
	HasMX       bool   `json:"has_mx"`
	HasSPF      bool   `json:"has_spf"`
	SPFRecord   string `json:"spf_record"`
	HasDMARC    bool   `json:"has_dmarc"`
	DMARCRecord string `json:"dmarc_record"`
	IsValid     bool   `json:"is_valid"`
}

func main() {
	// Serve static files from the "web" directory
	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", fs)

	// API endpoint to verify domains
	http.HandleFunc("/api/verify", verifyHandler)

	// Start the web server on port 8080
	port := ":8080"
	fmt.Printf("==================================================\n")
	fmt.Printf("  Email Domain Verifier Server Starting...       \n")
	fmt.Printf("  Open your browser at: http://localhost%s       \n")
	fmt.Printf("==================================================\n")
	log.Fatal(http.ListenAndServe(port, nil))
}

func verifyHandler(w http.ResponseWriter, r *http.Request) {
	// Enable CORS (Cross-Origin Resource Sharing)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	domain := strings.TrimSpace(r.URL.Query().Get("domain"))
	if domain == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "domain query parameter is required"})
		return
	}

	result := checkDomain(domain)
	json.NewEncoder(w).Encode(result)
}

func checkDomain(domain string) VerificationResult {
	var res VerificationResult
	res.Domain = domain

	// 1. MX Lookup
	mxrecords, _ := net.LookupMX(domain)
	if len(mxrecords) > 0 {
		res.HasMX = true
	}

	// 2. SPF Lookup
	txtrecords, _ := net.LookupTXT(domain)
	for _, record := range txtrecords {
		if strings.HasPrefix(record, "v=spf1") {
			res.HasSPF = true
			res.SPFRecord = record
			break
		}
	}

	// 3. DMARC Lookup
	dmarcRecords, _ := net.LookupTXT("_dmarc." + domain)
	for _, record := range dmarcRecords {
		if strings.HasPrefix(record, "v=DMARC1") {
			res.HasDMARC = true
			res.DMARCRecord = record
			break
		}
	}

	res.IsValid = res.HasMX && res.HasSPF && res.HasDMARC
	return res
}
