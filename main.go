// main.go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
    "github.com/Rokki-Khazratov/golang-scrapper/api/data-scrappers/scrapper"

)

func handleRequests() {
	r := mux.NewRouter()
	c := cors.AllowAll() // Move the CORS handler definition here
	r.HandleFunc("/api/laptops/", getAllLaptops).Methods("GET")
	r.HandleFunc("/api/product/{id}/", getProductDetail).Methods("GET")

	http.Handle("/", c.Handler(r))
	http.ListenAndServe(":8080", nil)
}


func main() {
	scrapper.ScrapeBrostore()
	handleRequests()
}

func getAllLaptops(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(scrapper.Laptops)
	if err != nil {
		http.Error(w, "Error converting to JSON", http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}

func getProductDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get the product ID from the URL parameters
	params := mux.Vars(r)
	id := params["id"]

	// Convert the ID to an integer
	productID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Ensure that the product ID is within the valid range
	if productID < 1 || productID > len(scrapper.Laptops) {
		http.Error(w, "Product ID out of range", http.StatusNotFound)
		return
	}

	// Get detailed information for the specified product ID and link
	link := scrapper.Laptops[productID-1].Link
	detailedProduct, err := scrapper.ScrapeBrostoreDetail(productID, link)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting product detail: %v", err), http.StatusInternalServerError)
		return
	}

	// Convert detailed product information to JSON
	jsonData, err := json.Marshal(detailedProduct)
	if err != nil {
		http.Error(w, "Error converting to JSON", http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}
