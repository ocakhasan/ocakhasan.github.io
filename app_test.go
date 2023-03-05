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
		name         string
		httpMethod   string
		query        string
		responseBody string
		statusCode   int
	}{
		{
			name:         "unsupported httpMethod",
			httpMethod:   http.MethodPost,
			query:        "1",
			responseBody: "unsupported httpMethod",
			statusCode:   http.StatusMethodNotAllowed,
		},
		{
			name:         "invalid input",
			httpMethod:   http.MethodGet,
			query:        "asd",
			responseBody: `invalid input`,
			statusCode:   http.StatusBadRequest,
		},
		{
			name:         "correct query param",
			httpMethod:   http.MethodGet,
			query:        "1",
			responseBody: `{"output":"I"}`,
			statusCode:   http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			path := fmt.Sprintf("/roman?query=%s", tc.query)
			request := httptest.NewRequest(tc.httpMethod, path, nil)
			responseRecorder := httptest.NewRecorder()

			romanHandler{}.ServeHTTP(responseRecorder, request)

			if responseRecorder.Code != tc.statusCode {
				t.Errorf("Want status '%d', got '%d'", tc.statusCode, responseRecorder.Code)
			}

			if strings.TrimSpace(responseRecorder.Body.String()) != tc.responseBody {
				t.Errorf("Want '%s', got '%s'", tc.responseBody, responseRecorder.Body)
			}
		})
	}
}
