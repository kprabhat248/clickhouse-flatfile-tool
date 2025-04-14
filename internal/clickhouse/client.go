// --- internal/clickhouse/client.go ---
package clickhouse

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ClickHouse/clickhouse-go/v2"
	"encoding/csv"
)

type Config struct {
	Host     string   `json:"host"`
	Port     string   `json:"port"`
	User     string   `json:"user"`
	Database string   `json:"database"`
	JWT      string   `json:"jwt"`
	Table    string   `json:"table"`
}

func NewClient(cfg Config) (clickhouse.Conn, error) {
	return clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)},
		Auth: clickhouse.Auth{
			Database: cfg.Database,
			Username: cfg.User,
			Password: cfg.JWT,
		},
		TLS: &tls.Config{InsecureSkipVerify: true},
	})
}

func FetchTableColumns(ctx context.Context, conn clickhouse.Conn, table string) ([]string, error) {
	rows, err := conn.Query(ctx, fmt.Sprintf("DESCRIBE TABLE %s", table))
	if err != nil {
		return nil, err
	}
	var cols []string
	for rows.Next() {
		var name, typ, def, codec, comment, ttl string
		if err := rows.Scan(&name, &typ, &def, &codec, &comment, &ttl); err != nil {
			return nil, err
		}
		cols = append(cols, name)
	}
	return cols, nil
}


func IngestCSVToTable(cfg Config, file string, columns []string) (int, error) {
	conn, err := NewClient(cfg)
	if err != nil {
		return 0, err
	}
	ctx := context.Background()
	f, err := os.Open(filepath.Join("data", file))
	if err != nil {
		return 0, err
	}
	defer f.Close()
	rdr := csv.NewReader(f)
	headers, _ := rdr.Read()
	colIndex := make([]int, len(columns))
	for i, name := range columns {
		for j, h := range headers {
			if h == name {
				colIndex[i] = j
			}
		}
	}
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES", cfg.Table, strings.Join(columns, ", "))
	batch, err := conn.PrepareBatch(ctx, query)
	if err != nil {
		return 0, err
	}
	count := 0
	for {
		rec, err := rdr.Read()
		if err != nil {
			break
		}
		row := make([]interface{}, len(colIndex))
		for i, idx := range colIndex {
			row[i] = rec[idx]
		}
		batch.Append(row...)
		count++
	}
	batch.Send()
	return count, nil
}