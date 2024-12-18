package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handleMetrics(w http.ResponseWriter, r *http.Request) {
	myValue := cfg.fileserverHits.Load()
	htmlData, err := ioutil.ReadFile("metrics.html")
	if err != nil {
		http.Error(w, "Unable to read metrics.html", http.StatusInternalServerError)
		return
	}

	formattedHTML := fmt.Sprintf(string(htmlData), myValue)

	w.Header().Set("Content-type", "text/html")
	fmt.Fprint(w, formattedHTML)
}
