package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/PuerkitoBio/goquery"
	"github.com/rs/cors"
	"encoding/json"
	"strconv"
	"errors"
	"regexp"
)

type Product struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	VariantTitle string `json:"variant-title"`
	Price        string `json:"price"`
	Image        string `json:"image"`
	Image1       string `json:"image1"` // Add the Image1 field
	Link         string `json:"link"`
}
var products []Product

var r = mux.NewRouter() // Initialize the router outside the handleRequests function

func handleRequests() {
	// Endpoint for the list of products
	r.HandleFunc("/api/products/", getAllProducts).Methods("GET")

	// Endpoint for detailed information about a specific product
	r.HandleFunc("/api/product/{id}/", getProductDetail).Methods("GET")

	// New endpoint for parsing a specific URL
	r.HandleFunc("/api/parse/", parseURL).Methods("POST")

	// Enable CORS middleware
	c := cors.AllowAll()

	// Wrap the existing router with CORS middleware
	http.Handle("/", c.Handler(r))

	// Start the server
	http.ListenAndServe(":8080", nil)
}

func parseURL(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    var requestBody struct {
        URL string `json:"url"`
    }

    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&requestBody); err != nil {
        http.Error(w, "Error decoding request body", http.StatusBadRequest)
        return
    }

    // Reset the products slice before scraping new data
    products = nil

    // Print the URL from the request body
    fmt.Println("[POST]: url: " + requestBody.URL)

    // Scrape data from the provided URL
    scrapeBrostore(w, requestBody.URL)

    // Return the parsed data as JSON
    jsonData, err := json.Marshal(products)
    if err != nil {
        http.Error(w, "Error converting to JSON", http.StatusInternalServerError)
        return
    }
    w.Write(jsonData)
}

func scrapeBrostore(w http.ResponseWriter, url string) {
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
		imageURL, _ := s.Find(".product-primary-image").Attr("data-srcset")
		imageURL = "https:" + strings.Fields(imageURL)[0]

		imageURL1, _ := s.Find(".product-secondary-image").Attr("data-srcset")
		imageURL1 = "https:" + strings.Fields(imageURL1)[0]  // Fix this line

		// Find the parent container that contains the link
		parent := s.Find(".product-card-title").Parent()

		// Find the link within the parent container
		linkSelection := parent.Find("a").First()
		link, exists := linkSelection.Attr("href")

		// Check if the link exists
		if exists {
			link = "https://brostore.uz" + link

			// Save information to the products slice
			product := Product{
				ID:    i + 1,
				Title: title,
				Price: price,
				Image: imageURL,
				Image1: imageURL1,
				Link:  link,
			}
			products = append(products, product)
		} else {
			fmt.Println("Link does not exist")
		}
	})

}

func scrapeProductInfo(doc *goquery.Document) (title, variantTitle, price string, err error) {
	// Extract title
	title = strings.TrimSpace(doc.Find(".product-title").Text())

	// Extract variant title
	variantTitle = strings.TrimSpace(doc.Find(".product-variant-title").Text())

	// Extract price
	price = doc.Find(".amount").Text()

	// Use a regular expression to capture the desired part of the price
	re := regexp.MustCompile(`(\d{2}\s\d{3}\s\d{3}\sсум)`)
	matches := re.FindStringSubmatch(price)
	if len(matches) > 1 {
		price = matches[1]
	}

	return title, variantTitle, price, nil
}

func cleanUpPrice(rawPrice string) (string, error) {
	// Split the raw price by whitespace
	priceParts := strings.Fields(rawPrice)

	// Take the first part as the cleaned-up price
	if len(priceParts) > 0 {
		return priceParts[0], nil
	}

	return "", errors.New("No valid price found")
}

func scrapeBrostoreDetail(id int, link string) (*Product, error) {
    // Create the full URL for the product detail page
    fullURL := link
    fmt.Printf(link)

    // Set a user-agent to mimic a browser request
    client := &http.Client{}
    req, err := http.NewRequest("GET", fullURL, nil)
    if err != nil {
        return nil, fmt.Errorf("Error creating request: %v", err)
    }
    req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")

    // Perform the request
    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("Error making request: %v", err)
    }
    defer resp.Body.Close()

    // Check if the request was successful (status code 200)
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("Failed to retrieve data. Status Code: %d", resp.StatusCode)
    }

    // Use goquery to parse the HTML response
    doc, err := goquery.NewDocumentFromReader(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("Error parsing HTML: %v", err)
    }

    // Extract detailed information from the product detail page using the scrapeProductInfo function
    title, variantTitle, price, err := scrapeProductInfo(doc)
    if err != nil {
        return nil, fmt.Errorf("Error extracting detailed information: %v", err)
    }

    imageURL, _ := doc.Find(".product-secondary-image").Attr("data-srcset")
    detailedInfo := struct {
        Title        string `json:"title"`
        VariantTitle string `json:"variant-title"`
        Price        string `json:"price"`
        Image        string `json:"image"`
    }{
        Title:        title,
        VariantTitle: variantTitle,
        Price:        price,
        Image:        "https:" + strings.Fields(strings.TrimSpace(imageURL))[0],
    }

    // Create a new Product instance with the detailed information and the provided ID
    detailedProduct := &Product{
        ID:           id,
        Title:        detailedInfo.Title,
        VariantTitle: detailedInfo.VariantTitle,
        Price:        detailedInfo.Price,
        Image:        detailedInfo.Image,
        Link:         link,
    }

    return detailedProduct, nil
}

func getAllProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(products)
	if err != nil {
		http.Error(w, "Error converting to JSON", http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}

func getProductDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id := params["id"]

	productID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	if productID < 1 || productID > len(products) {
		http.Error(w, "Product ID out of range", http.StatusNotFound)
		return
	}

	link := products[productID-1].Link
	detailedProduct, err := scrapeBrostoreDetail(productID, link)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting product detail: %v", err), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(detailedProduct)
	if err != nil {
		http.Error(w, "Error converting to JSON", http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}

func main() {
	// Initialize the products slice
	products = make([]Product, 0)

	// Start the server
	handleRequests()
}

