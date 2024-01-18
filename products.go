package main

import (
	"fmt"
	// "html/template"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"encoding/json"

)

// Product struct to represent the data structure
type Product struct {
	Title string `json:"title"`
	Price string `json:"price"`
	Image string `json:"image"`
}

var products []Product

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

		// Save information to the products slice
		product := Product{
			Title: title,
			Price: price,
			Image: imageURL,
		}
		products = append(products, product)
	})
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	http.ListenAndServe(":8080", nil)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем заголовок Content-Type для ответа в формате JSON
	w.Header().Set("Content-Type", "application/json")

	// Преобразуем данные в формат JSON
	jsonData, err := json.Marshal(products)
	if err != nil {
		http.Error(w, "Error converting to JSON", http.StatusInternalServerError)
		return
	}

	// Отправляем данные в ответе
	w.Write(jsonData)
}

func main() {
	scrapeBrostore()
	handleRequests()
}
