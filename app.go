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
		http.Error(w, "unsupported httpMethod", http.StatusMethodNotAllowed)
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
