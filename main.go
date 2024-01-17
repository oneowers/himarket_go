package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Product struct {
	Title string `json:"title"`
	Price string `json:"price"`
	Image string `json:"image"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	file, err := ioutil.ReadFile("go_data.json")
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		http.Error(w, "Error reading JSON file", http.StatusInternalServerError)
		return
	}

	var products []Product
	err = json.Unmarshal(file, &products)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		http.Error(w, "Error unmarshalling JSON", http.StatusInternalServerError)
		return
	}

	// Set Content-Type header
	w.Header().Set("Content-Type", "application/json")

	// Write JSON response
	json.NewEncoder(w).Encode(products)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
