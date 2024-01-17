import os
import requests
from bs4 import BeautifulSoup

def scrape_brostore():
    url = "https://brostore.uz/collections/noutbuki"
    headers = {
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"
    }

    response = requests.get(url, headers=headers)

    if response.status_code == 200:
        soup = BeautifulSoup(response.text, 'lxml')

        # Find all product-card elements
        product_cards = soup.find_all('product-card', class_='product-card text-center product-card--content-spacing-false product-card--border-false has-shadow--false')

        for i, card in enumerate(product_cards):
            # Extract information from each product-card
            title = card.select_one('.product-card-title').text.strip()
            price = card.select_one('.amount').text.strip()
            image_url = "https:" + card.select_one('.product-secondary-image')['data-srcset'].split(' ')[0]

            # Print the information
            print(f"Title: {title}")
            print(f"Price: {price}")
            print(f"Image URL: {image_url}")
            print('-' * 30)

            # Save information to data.json
            data = {
                "title": title,
                "price": price,
                "image": image_url
            }
            save_to_json(data)

    else:
        print(f"Failed to retrieve data. Status Code: {response.status_code}")

def save_to_json(data):
    import json

    # Load existing data from data.json if it exists
    existing_data = []
    if os.path.exists('data.json') and os.path.getsize('data.json') > 0:
        with open('data.json', 'r') as json_file:
            try:
                existing_data = json.load(json_file)
            except json.JSONDecodeError:
                pass

    # Append the new data
    existing_data.append(data)

    # Save the combined data back to data.json
    with open('data.json', 'w') as json_file:
        json.dump(existing_data, json_file, indent=2, ensure_ascii=False)

if __name__ == "__main__":
    scrape_brostore()
