package httphelper

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_MethodPlexer_UnsupportedMethod(t *testing.T) {
	t.Parallel()
	// given
	p := MethodPlexMiddleware(MethodPlexer{})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPatch, "/asdf", http.NoBody)
	// when
	p.ServeHTTP(w, r)
	// then
	if got := w.Result().StatusCode; got != http.StatusMethodNotAllowed {
		t.Errorf("status code. Got %d, want %d", got, http.StatusMethodNotAllowed)
	}
}

func Test_MethodPlexer_SupportedMethods(t *testing.T) {
	t.Parallel()
	var okResponder = func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	tcs := map[string]struct {
		plexer MethodPlexer
		method string
	}{
		"get": {
			plexer: MethodPlexer{
				Get: okResponder,
			},
			method: http.MethodGet,
		},
		"post": {
			plexer: MethodPlexer{
				Post: okResponder,
			},
			method: http.MethodPost,
		},
		"put": {
			plexer: MethodPlexer{
				Put: okResponder,
			},
			method: http.MethodPut,
		},
	}
	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			// given
			p := MethodPlexMiddleware(tc.plexer)
			w := httptest.NewRecorder()
			r := httptest.NewRequest(tc.method, "/asdf", http.NoBody)
			// when
			p.ServeHTTP(w, r)
			// then
			if got := w.Result().StatusCode; got != http.StatusOK {
				t.Errorf("status code. Got %d, want %d", got, http.StatusOK)
			}
		})
	}
}

func Test_MethodPlexer_SupportedMethod_HandlerNil(t *testing.T) {
	t.Parallel()
	for _, method := range []string{http.MethodGet, http.MethodPost, http.MethodPut} {
		t.Run(method, func(t *testing.T) {
			t.Parallel()
			// given
			p := MethodPlexMiddleware(MethodPlexer{})
			w := httptest.NewRecorder()
			r := httptest.NewRequest(method, "/asdf", http.NoBody)
			// when
			p.ServeHTTP(w, r)
			// then
			if got := w.Result().StatusCode; got != http.StatusMethodNotAllowed {
				t.Errorf("status code. Got %d, want %d", got, http.StatusMethodNotAllowed)
			}
		})
	}
}
