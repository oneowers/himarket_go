package scrapper

import (
	"fmt"
	"net/http"
	"strings"
	"github.com/PuerkitoBio/goquery"
	"encoding/json"
	// "github.com/gorilla/mux"
	// "strconv"
	"errors"
	"regexp"
)

type Product struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	VariantTitle string `json:"variant-title"`
	Price        string `json:"price"`
	Image        string `json:"image"`
	Image1       string `json:"image1"`
	Link         string `json:"link"`
}

var Products []Product

func ScrapeBrostore() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://brostore.uz/collections/noutbuki", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed to retrieve data. Status Code: %d\n", resp.StatusCode)
		return
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return
	}

	doc.Find(".product-card.text-center.product-card--content-spacing-false.product-card--border-false.has-shadow--false").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find(".product-card-title").Text())
		price := strings.TrimSpace(s.Find(".amount").Text())
		imageURL, _ := s.Find(".product-primary-image").Attr("data-srcset")
		imageURL = "https:" + strings.Fields(imageURL)[0]

		imageURL1, _ := s.Find(".product-secondary-image").Attr("data-srcset")
		imageURL1 = "https:" + strings.Fields(imageURL1)[0]

		parent := s.Find(".product-card-title").Parent()
		linkSelection := parent.Find("a").First()
		link, exists := linkSelection.Attr("href")

		if exists {
			link = "https://brostore.uz" + link
			product := Product{
				ID:    i + 1,
				Title: title,
				Price: price,
				Image: imageURL,
				Image1: imageURL1,
				Link:  link,
			}
			Products = append(Products, product)
		} else {
			fmt.Println("Link does not exist")
		}
	})
}

func scrapeProductInfo(doc *goquery.Document) (title, variantTitle, price string, err error) {
	title = strings.TrimSpace(doc.Find(".product-title").Text())
	variantTitle = strings.TrimSpace(doc.Find(".product-variant-title").Text())
	price = doc.Find(".amount").Text()
	re := regexp.MustCompile(`(\d{2}\s\d{3}\s\d{3}\sсум)`)
	matches := re.FindStringSubmatch(price)
	if len(matches) > 1 {
		price = matches[1]
	}
	return title, variantTitle, price, nil
}

func cleanUpPrice(rawPrice string) (string, error) {
	priceParts := strings.Fields(rawPrice)
	if len(priceParts) > 0 {
		return priceParts[0], nil
	}
	return "", errors.New("No valid price found")
}

func ScrapeBrostoreDetail(id int, link string) (*Product, error) {
	fullURL := link
	fmt.Printf(link)

	client := &http.Client{}
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %v", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to retrieve data. Status Code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error parsing HTML: %v", err)
	}

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

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(Products)
	if err != nil {
		http.Error(w, "Error converting to JSON", http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}
