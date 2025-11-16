package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/jmoiron/sqlx"
)

var (
	TestURL    string = "http://localhost:6060"
	AdminToken string = "admin"
	UserToken  string = "user"
	DbURL      string
	httpClient = &http.Client{}
)

func RequireStatusCode(t *testing.T, resp *http.Response, expected int) {
	t.Helper()
	if resp.StatusCode != expected {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected status %d, got %d. Body: %s", expected, resp.StatusCode, string(body))
	}
}

func CleanDB(db *sqlx.DB) error {
	rows, err := db.Queryx(`
		SELECT schemaname, tablename
		FROM pg_tables
		WHERE schemaname NOT IN ('pg_catalog','information_schema')
		  AND schemaname <> 'pg_toast'
	`)
	if err != nil {
		return err
	}
	defer func(rows *sqlx.Rows) {
		_ = rows.Close()
	}(rows)

	var parts []string
	for rows.Next() {
		var schemaname, tablename string
		if err := rows.Scan(&schemaname, &tablename); err != nil {
			return err
		}
		parts = append(parts, fmt.Sprintf(`"%s"."%s"`, schemaname, tablename))
	}

	if len(parts) == 0 {
		return nil
	}

	q := "TRUNCATE TABLE " + strings.Join(parts, ", ") + " RESTART IDENTITY CASCADE"
	_, err = db.Exec(q)
	return err
}

func ReadBody(t *testing.T, resp *http.Response) []byte {
	t.Helper()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	return body
}

func PostJSON(t *testing.T, path string, body interface{}, token string) *http.Response {
	t.Helper()
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal body: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, TestURL+path, bytes.NewReader(b))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	return resp
}

func GetJSON(t *testing.T, path string, body interface{}, token string) *http.Response {
	t.Helper()

	jsonBody, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("wrong json format: %v", err)
	}
	req, err := http.NewRequest(http.MethodGet, TestURL+path, bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	return resp
}

func PatchJSON(t *testing.T, path string, body interface{}, token string) *http.Response {
	t.Helper()
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal body: %v", err)
	}

	req, err := http.NewRequest(http.MethodPatch, TestURL+path, bytes.NewReader(b))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	return resp
}
