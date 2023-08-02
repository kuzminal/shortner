package short

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestShortening(t *testing.T) {
	t.Parallel()
	body, err := json.Marshal(map[string]any{
		"url": "https://go.dev",
		"key": "go",
	})
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, shorteningRoute, bytes.NewReader(body))
	srv := NewServer()
	srv.ServeHTTP(w, r)

	if w.Code != http.StatusCreated {
		t.Errorf("got status code = %d, want %d", w.Code, http.StatusCreated)
	}
	if !strings.Contains(w.Body.String(), `"go"`) {
		t.Errorf("got body = %s\twant contains %s", w.Body.String(), `"go"`)
	}
}
