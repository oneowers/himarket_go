package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Product struct to represent the data structure
type Product struct {
	Title string `json:"title"`
	Price string `json:"price"`
	Image string `json:"image"`
}

func scrapeBrostore() {
	url := "https://brostore.uz/collections/noutbuki"

	// Set a user-agent to mimic a browser request
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	// Check if the request was successful (status code 200)
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed to retrieve data. Status Code: %d\n", resp.StatusCode)
		return
	}

	// Use goquery to parse the HTML response
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return
	}

	// Find all product-card elements
	doc.Find(".product-card.text-center.product-card--content-spacing-false.product-card--border-false.has-shadow--false").Each(func(i int, s *goquery.Selection) {
		// Extract information from each product-card
		title := strings.TrimSpace(s.Find(".product-card-title").Text())
		price := strings.TrimSpace(s.Find(".amount").Text())
		imageURL, _ := s.Find(".product-secondary-image").Attr("data-srcset")
		imageURL = "https:" + strings.Fields(imageURL)[0]

		// Print the information
		fmt.Printf("Title: %s\n", title)
		fmt.Printf("Price: %s\n", price)
		fmt.Printf("Image URL: %s\n", imageURL)
		fmt.Println(strings.Repeat("-", 30))

		// Save information to go_data.json
		product := Product{
			Title: title,
			Price: price,
			Image: imageURL,
		}
		saveToJSON(product)
	})
}

func saveToJSON(product Product) {
	// Load existing data from go_data.json if it exists
	var existingData []Product
	data, err := ioutil.ReadFile("go_data.json")
	if err == nil {
		if err := json.Unmarshal(data, &existingData); err != nil {
			fmt.Println("Error decoding JSON:", err)
			return
		}
	}

	// Append the new data
	existingData = append(existingData, product)

	// Save the combined data back to go_data.json
	jsonData, err := json.MarshalIndent(existingData, "", "  ")
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	err = ioutil.WriteFile("go_data.json", jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing to go_data.json:", err)
	}
}

func main() {
	scrapeBrostore()
}
