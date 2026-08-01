# Email Domain Verifier

A clean, simple, and efficient web application built with Go that checks the configuration of email domains. It queries DNS records to verify if a domain is properly set up to handle emails securely and avoid spam filters.

## Key Features

The tool performs the following checks for any given domain:
1. **MX Records Check**: Verifies if the domain has Mail Exchange (MX) records configured (required to receive emails).
2. **SPF Record Check**: Looks up Sender Policy Framework (SPF) records and retrieves the policy details (used to prevent sender address spoofing).
3. **DMARC Record Check**: Looks up Domain-based Message Authentication, Reporting, and Conformance (DMARC) records (used to define how the receiver should handle emails that fail SPF/DKIM).
4. **Overall Validity**: Reports if the domain has all three essential configurations in place.

---

## How to Run

### Prerequisites
- Go installed on your machine (version 1.20+ recommended).

### Running the Application
1. Open your terminal in the project directory.
2. Run the application:
   ```bash
   go run main.go
   ```
3. Open your web browser and navigate to:
   ```
   http://localhost:8080
   ```

---

## API Reference

The server also exposes a public JSON endpoint for automated checks.

### Verify Domain
- **Endpoint**: `/api/verify`
- **Method**: `GET`
- **Query Parameter**: `domain` (e.g., `api/verify?domain=google.com`)

#### Sample Response:
```json
{
  "domain": "google.com",
  "has_mx": true,
  "has_spf": true,
  "spf_record": "v=spf1 include:_spf.google.com ~all",
  "has_dmarc": true,
  "dmarc_record": "v=DMARC1; p=reject; rua=mailto:mailauth-reports@google.com",
  "is_valid": true
}
```

---

## Project Structure
- `main.go`: The Go backend server that handles HTTP requests, serves static files, and performs DNS queries.
- `web/`: The frontend client directory.
  - `index.html`: Main UI template.
  - `style.css`: Stylesheet.
  - `app.js`: Client-side logic for making API calls and rendering results.
