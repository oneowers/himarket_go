// api/api.go
package api

import (
	"encoding/json"
	"net/http"
	"os"
)

type DataResponse struct {
	Data []map[string]interface{} `json:"data"`
}

func GetDataHandler(w http.ResponseWriter, r *http.Request) {
	// Read the go_data.json file
	file, err := os.Open("../go_data.json") // Adjust the file path accordingly
	if err != nil {
		http.Error(w, "Error reading data", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Parse the JSON data
	var data []map[string]interface{}
	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		http.Error(w, "Error decoding data", http.StatusInternalServerError)
		return
	}

	// Create a response object
	response := DataResponse{
		Data: data,
	}

	// Convert the response to JSON and send it
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
