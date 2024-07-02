package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func panicHandler(_ http.ResponseWriter, _ *http.Request) {
	panic("foo")
}

func TestRecover(t *testing.T) {
	t.Parallel()

	h := Recover(http.HandlerFunc(panicHandler))
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	wantCode := http.StatusInternalServerError
	wantCT := "application/json; charset=utf-8"
	wantBB := []byte("{\"error\":\"internal server error\"}\n")

	h.ServeHTTP(w, r)
	if code := w.Code; code != wantCode {
		t.Errorf("mismatch (want, got):\n%d, %d", wantCode, code)
	}
	if ct := w.Header().Get("content-type"); ct != wantCT {
		t.Errorf("mismatch (want, got):\n%s, %s", wantCT, ct)
	}
	bb := w.Body.Bytes()
	if bytes.Equal(bb, wantBB) {
		t.Errorf("mismatch (want, got):\n%s, %s", wantBB, bb)
	}
}
