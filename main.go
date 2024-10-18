package main

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"os"
	"text/template"
)

type Report struct {
	CSPReport struct {
		DocumentURI        string `json:"document-uri"`
		Referrer           string `json:"referrer"`
		ViolatedDirective  string `json:"violated-directive"`
		EffectiveDirective string `json:"effective-directive"`
		OriginalPolicy     string `json:"original-policy"`
		Disposition        string `json:"disposition"`
		BlockedURI         string `json:"blocked-uri"`
		StatusCode         int    `json:"status-code"`
		ScriptSample       string `json:"script-sample"`
	} `json:"csp-report"`
}

type RequestInfo struct {
	Host string
	URL  string
}

func main() {

	logger := slog.NewLogLogger(slog.NewJSONHandler(os.Stdout, nil), slog.LevelInfo)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /test", ViolationMakerHandler)
	mux.HandleFunc("POST /cspro", CSPViolationHandler)
	mux.HandleFunc("POST /", CSPViolationHandler)

	server := &http.Server{
		Addr:     ":8080",
		Handler:  mux,
		ErrorLog: logger,
	}

	server.ListenAndServe()
}

func CSPViolationHandler(w http.ResponseWriter, r *http.Request) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	var violation Report
	err := json.NewDecoder(r.Body).Decode(&violation)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	logger.Info("csp-violation", "violation", violation.CSPReport)
	w.WriteHeader(http.StatusOK)
}

func ViolationMakerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Reporting-To", `cspro="http://localhost:8080`)
	// w.Header().Set("Content-Security-Policy", "default-src '*'; script-src 'sha256-oFpGexFYa81iRs0wnRObU36W0bkPCqdNLJg7Vggphvk='; report-uri csp-endpoint;")
	w.Header().Set("Content-Security-Policy-Report-Only", "default-src 'self'; script-src 'sha256-oFpGexFYa81iRs0wnRObU36W0bkPCqdNLJg7Vggphvk='; report-uri cspro;")

	t, err := template.ParseFiles("./templates/index.html.tmpl")
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t.Execute(w, nil)
}
