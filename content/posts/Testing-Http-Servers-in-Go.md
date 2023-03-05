---
layout: post
title: Testing HTTP Servers in Go
summary: Let's dive into end to to end testing for http servers in Go. 
date: 2023-03-05
tags: [go, http, testing]
---

# Introduction

Go is an excellent programming language for building HTTP servers, thanks to its `net/http` package in the standard library, which makes it easy to attach HTTP handlers to any Go program. The standard library also includes packages that facilitate testing HTTP servers, making it just as effortless to test them as it is to build them. 

Nowadays, test coverage is widely accepted as an essential and valuable part of software development. Developers invest time in testing their code to get quick feedback when making changes, and a good test suite becomes an invaluable component of the software project when combined with continuous integration and delivery methodologies.

Given the importance of a good test suite, what approach should developers using Go take when testing their HTTP servers? This article provides everything you need to know to test your Go HTTP servers thoroughly.

## Http Server For Conversion of Roman Numerals

We will have a web server which gives the roman numeral of the given number. We will only have 1 endpoint.

- Show the roman numeral of the number  `GET /roman`

### Example Request and Response

**Request**

```bash
curl --location --request GET 'http://localhost:8080/roman?query=1'
```

**Response**

```json
{
    "output": "I"
}
```

### Code and Explanation

```go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

var (
	nums    = []int{1, 4, 5, 9, 10, 40, 50, 90, 100, 400, 500, 900, 1000}
	symbols = []string{"I", "IV", "V", "IX", "X", "XL", "L", "XC", "C", "CD", "D", "CM", "M"}
)

func convertIntegerToRoman(input int) string {
	var (
		i      = len(nums) - 1
		result string
	)

	for input > 0 {
		division := input / nums[i]
		input = input % nums[i]

		for division > 0 {
			result += symbols[i]
			division = division - 1
		}

		i = i - 1
	}

	return result
}

type romanHandler struct{}

func (h romanHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		http.Error(w, "unsupported method", http.StatusMethodNotAllowed)
		return
	}

	input := r.URL.Query().Get("query")
	inputInt, err := strconv.Atoi(input)
	if err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	output := convertIntegerToRoman(inputInt)

	response := map[string]interface{}{
		"output": output,
	}

	if err := json.NewEncoder(w).Encode(&response); err != nil {
		return
	}
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/roman", romanHandler{})

	log.Fatal(http.ListenAndServe(":8080", mux))
}
```

Points
- The function `convertIntegerToRoman` takes an integer and return the roman numeral conversion of the number. Please have a look on [Convert Number Into Roman Numeral](https://www.geeksforgeeks.org/converting-decimal-number-lying-between-1-to-3999-to-roman-numerals/)

- We accept a single query parameter named `query` in the URL which should have the number which will be converted. 

The struct implements the `http.Handler` interface by implementing the method of
- `ServeHTTP(ResponseWriter, *Request)`

## Testing Of the Server

The whole purpose of this blog was to learn how to test http servers in Go, so let's find out.

As we mentioned in the beginning Go has all of the tools we need to both create `net/http`  and test `net/http/httptest`. All of the tools are included in the `net` module. 

Let's create a file named `main_test.go` which has all of the tests for the HTTP Server.

### Tests

```go
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

```

To test the handler, we use the common table-driven approach and provide three cases:

1. the http method is not correct
2. http method is correct, but the query param is invalid
3. both http method and query param is valid.


For each case, we run a subtest that creates a new request and a response recorder. We use the `httptest.NewRequest` function to create an `http.Request` struct, which represents an incoming request to the handler. This allows us to simulate a real request without relying on an actual HTTP server.

However, this function only handles the request half of the testing. To handle the response half, we use `httptest.ResponseRecorder`, which records the mutations of the `http.ResponseWriter` and enables us to make assertions on it later in the test.

By using this duo of `httptest.ResponseRecorder` and `http.Request`, we can successfully test any HTTP handler in Go. Running the test will produce the following output.

```go
=== RUN   TestRomanHandler
=== RUN   TestRomanHandler/unsupported_method
=== RUN   TestRomanHandler/invalid_input
=== RUN   TestRomanHandler/correct_query_param
--- PASS: TestRomanHandler (0.00s)
    --- PASS: TestRomanHandler/unsupported_method (0.00s)
    --- PASS: TestRomanHandler/invalid_input (0.00s)
    --- PASS: TestRomanHandler/correct_query_param (0.00s)
PASS
```

### REFERENCES
- [net/http](https://pkg.go.dev/net/http)
- [Testing HTTP Servers By Ieftimov](https://ieftimov.com/posts/testing-in-go-testing-http-servers/)
- [Converting Decimal To Roman](https://www.geeksforgeeks.org/converting-decimal-number-lying-between-1-to-3999-to-roman-numerals/)