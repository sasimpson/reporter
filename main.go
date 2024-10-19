package main

import (
	_ "embed"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"os"
	"text/template"
)

var (
	//go:embed templates/index.html.tmpl
	indexTemplate string
	logger        *log.Logger
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

func main() {

	logger = slog.NewLogLogger(slog.NewJSONHandler(os.Stdout, nil), slog.LevelInfo)

	mux := http.NewServeMux()
	mux.HandleFunc("/", ViolationMakerHandler)
	mux.HandleFunc("POST /cspro", CSPViolationHandler)

	server := &http.Server{
		Addr:     ":8080",
		Handler:  mux,
		ErrorLog: logger,
	}

	logger.Print("starting service on port :8080")
	err := server.ListenAndServe()
	if err != nil {
		return
	}
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

func ViolationMakerHandler(w http.ResponseWriter, _ *http.Request) {
	// If you wanted to hard block anything not adhering to the CSP, you can replace the
	// Content-Security-Policy-Report-Only with Content-Security-Policy.
	w.Header().Set("Content-Security-Policy-Report-Only", "default-src 'self'; script-src 'sha256-oFpGexFYa81iRs0wnRObU36W0bkPCqdNLJg7Vggphvk='; report-uri cspro;")
	// locations to report any violations to, will add the "name" to the url
	w.Header().Set("Reporting-To", `cspro="http://localhost:8080`)

	t, err := template.New("index").Parse(indexTemplate)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, nil)
	if err != nil {
		return
	}
}
