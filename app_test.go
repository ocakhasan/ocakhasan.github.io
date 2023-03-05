package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRomanHandler(t *testing.T) {
	tt := []struct {
		name       string
		method     string
		query      string
		want       string
		statusCode int
	}{
		{
			name:       "unsupported method",
			method:     http.MethodPost,
			query:      "1",
			want:       "unsupported method",
			statusCode: http.StatusMethodNotAllowed,
		},
		{
			name:       "invalid input",
			method:     http.MethodGet,
			query:      "asd",
			want:       `invalid input`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "correct query param",
			method:     http.MethodGet,
			query:      "1",
			want:       `{"output":"I"}`,
			statusCode: http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			path := fmt.Sprintf("/roman?query=%s", tc.query)
			request := httptest.NewRequest(tc.method, path, nil)
			responseRecorder := httptest.NewRecorder()

			romanHandler{}.ServeHTTP(responseRecorder, request)

			if responseRecorder.Code != tc.statusCode {
				t.Errorf("Want status '%d', got '%d'", tc.statusCode, responseRecorder.Code)
			}

			if strings.TrimSpace(responseRecorder.Body.String()) != tc.want {
				t.Errorf("Want '%s', got '%s'", tc.want, responseRecorder.Body)
			}
		})
	}
}
