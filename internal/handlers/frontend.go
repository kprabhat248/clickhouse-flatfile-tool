// --- internal/handlers/frontend.go ---
package handlers

import (
	"net/http"
	"os"
)

func ServeFrontend(w http.ResponseWriter, r *http.Request) {
	html, err := os.ReadFile("web/index.html")
	if err != nil {
		http.Error(w, "Could not load UI", 500)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(html)
}