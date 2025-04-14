// --- internal/handlers/ingest.go ---
package handlers

import (
	"encoding/json"
	"net/http"

	"clickhouse-flatfile-tool/internal/clickhouse"
)

type IngestRequest struct {
	Source     string          `json:"source"`
	Target     string          `json:"target"`
	FilePath   string          `json:"filePath"`
	ClickHouse clickhouse.Config `json:"clickhouse"`
	Columns    []string        `json:"columns"`
}

func IngestHandler(w http.ResponseWriter, r *http.Request) {
	var req IngestRequest
	json.NewDecoder(r.Body).Decode(&req)
	if req.Source == "flatfile" && req.Target == "clickhouse" {
		count, err := clickhouse.IngestCSVToTable(req.ClickHouse, req.FilePath, req.Columns)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Ingested successfully",
			"count":   count,
		})
	} else {
		http.Error(w, "Unsupported flow", 400)
	}
}
