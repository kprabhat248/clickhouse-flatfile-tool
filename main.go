// --- main.go ---
package main

import (
	"log"
	"net/http"

	"clickhouse-flatfile-tool/internal/handlers"
)

func main() {
	http.HandleFunc("/upload", handlers.UploadFileHandler)
	http.HandleFunc("/columns/flatfile", handlers.FlatFileColumnsHandler)
	http.HandleFunc("/columns/clickhouse", handlers.ClickHouseColumnsHandler)
	http.HandleFunc("/ingest", handlers.IngestHandler)
	http.HandleFunc("/", handlers.ServeFrontend)

	log.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}