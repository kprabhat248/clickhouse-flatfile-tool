// --- internal/handlers/columns.go ---
package handlers

import (
	"context"
	"encoding/csv"
	"encoding/json"
	
	"net/http"
	"os"
	"path/filepath"

	"clickhouse-flatfile-tool/internal/clickhouse"
)

func FlatFileColumnsHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("file")
	f, err := os.Open(filepath.Join("data", path))
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	defer f.Close()

	rdr := csv.NewReader(f)
	headers, err := rdr.Read()
	if err != nil {
		http.Error(w, "Could not read headers", 500)
		return
	}
	json.NewEncoder(w).Encode(headers)
}

func ClickHouseColumnsHandler(w http.ResponseWriter, r *http.Request) {
	var cfg clickhouse.Config
	json.NewDecoder(r.Body).Decode(&cfg)
	conn, err := clickhouse.NewClient(cfg)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	columns, err := clickhouse.FetchTableColumns(context.Background(), conn, cfg.Table)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(columns)
}
