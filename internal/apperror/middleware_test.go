package apperror

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddleware(t *testing.T) {
	testCases := []struct {
		name          string
		handler       appHandler
		expectedCode  int
		expectedError string
	}{
		{
			name: "valid request",
			handler: func(w http.ResponseWriter, r *http.Request) error {
				w.WriteHeader(http.StatusOK)
				return nil
			},
			expectedCode:  http.StatusOK,
			expectedError: "",
		},
		{
			name: "error response with AppError",
			handler: func(w http.ResponseWriter, r *http.Request) error {
				return NewAppError(nil, "error message", "developer message", "code")
			},
			expectedCode:  http.StatusBadRequest,
			expectedError: `{"message":"error message","developerMessage":"developer message","code":"code"}`,
		},
		{
			name: "error response with ErrNotFound",
			handler: func(w http.ResponseWriter, r *http.Request) error {
				return ErrNotFound
			},
			expectedCode:  http.StatusNotFound,
			expectedError: `{"message":"not found","code":"US-0000003"}`,
		},
		{
			name: "error response with system error",
			handler: func(w http.ResponseWriter, r *http.Request) error {
				return errors.New("system error")
			},
			expectedCode:  http.StatusTeapot,
			expectedError: `{"message":"internal system error","developerMessage":"system error","code":"US-000000"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			Middleware(tc.handler)(recorder, req)

			if recorder.Code != tc.expectedCode {
				t.Errorf("expected status code %d but got %d", tc.expectedCode, recorder.Code)
			}

			if tc.expectedError != "" && recorder.Body.String() != tc.expectedError {
				t.Errorf("expected body %s but got %s", tc.expectedError, recorder.Body.String())
			}
		})
	}
}
func MiddlewareAlt(h appHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err != nil {
			w.Header().Set("Content-Type", "application/type")
			switch e := err.(type) {
			case *AppError:
				if e == ErrNotFound {
					w.WriteHeader(http.StatusNotFound)
					w.Write(ErrNotFound.Marshal())
					return
				}
				w.WriteHeader(http.StatusBadRequest)
				w.Write(e.Marshal())
			default:
				w.WriteHeader(http.StatusTeapot)
				w.Write(SystemError(err).Marshal())
			}
		}
	}
}

func BenchmarkMiddleware(b *testing.B) {
	handler := Middleware(func(w http.ResponseWriter, r *http.Request) error {
		return nil
	})

	req, _ := http.NewRequest("GET", "http://localhost:8080/", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler(w, req)
	}
}
func BenchmarkMiddlewareAlt(b *testing.B) {
	handler := MiddlewareAlt(func(w http.ResponseWriter, r *http.Request) error {
		return nil
	})

	req, _ := http.NewRequest("GET", "http://localhost:8080/", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler(w, req)
	}
}
