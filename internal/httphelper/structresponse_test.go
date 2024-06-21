package httphelper

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_StructResponse_ok(t *testing.T) {
	t.Parallel()

	type something struct {
		A    string `json:"a"`
		AaBb int    `json:"aa_bb"`
	}

	// given
	expBody := `
	{
		"a": "hello",
		"aa_bb": 2
	}
	`

	w := httptest.NewRecorder()
	rt := httptest.NewRequest(http.MethodGet, "/asdf", http.NoBody)
	handler := StructResponse(func(r *http.Request) (something, int, *HTTPError) {
		if r != rt {
			t.Errorf("mismatched request. Got %v, want %v", r, rt)
		}
		return something{A: "hello", AaBb: 2}, http.StatusNoContent, nil
	})
	// when
	handler.ServeHTTP(w, rt)
	// then
	if got := w.Result().StatusCode; got != http.StatusNoContent {
		t.Errorf("status code. Got %d, want %d", got, http.StatusNoContent)
	}
	body, err := io.ReadAll(w.Result().Body)
	if err != nil {
		t.Fatal(err)
	}
	ensureJSONEqual(t, body, []byte(expBody))
}

func ensureJSONEqual(t *testing.T, a, b []byte) {
	t.Helper()
	var aU, bU any
	errA := json.Unmarshal(a, &aU)
	errB := json.Unmarshal(b, &bU)
	if combinedErr := errors.Join(errA, errB); combinedErr != nil {
		t.Fatal(combinedErr)
	}
	if diff := cmp.Diff(aU, bU); diff != "" {
		t.Errorf("diff(-a+b): %s\n", diff)
	}
}

func Test_StructResponse_err(t *testing.T) {
	t.Parallel()
	// given

	w := httptest.NewRecorder()
	rt := httptest.NewRequest(http.MethodGet, "/asdf", http.NoBody)
	handler := StructResponse(func(r *http.Request) (int, int, *HTTPError) {
		if r != rt {
			t.Errorf("mismatched request. Got %v, want %v", r, rt)
		}
		return 0, 0, NewError("What doing?", http.StatusExpectationFailed, errors.New("internal info"))
	})
	// when
	handler.ServeHTTP(w, rt)
	// then
	if got := w.Result().StatusCode; got != http.StatusExpectationFailed {
		t.Errorf("status code. Got %d, want %d", got, http.StatusExpectationFailed)
	}
	body, err := io.ReadAll(w.Result().Body)
	if err != nil {
		t.Fatal(err)
	}

	if strings.Contains(string(body), "internal info") {
		t.Errorf("body %q contains internal info", string(body))
	}
	if !strings.Contains(string(body), "What doing?") {
		t.Errorf("body %q does not contain message What doing?", string(body))
	}
}
