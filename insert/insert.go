package main

import (
	"flag"
	"fmt"
	"log"
	"tecnoric"
	"tecnoric/woocommerce"
)

func main() {
	inputFileName := flag.String("input", "", "input JSON file with final products")
	categoryName := flag.String("category", "", "the name of the category")
	flag.Parse()

	products := tecnoric.LoadFinalProducts(*inputFileName)
	insertProducts(products, *categoryName)
}

func insertProducts(products []tecnoric.FinalProduct, category string) {
	w := woocommerce.New()

	categories, err := w.SearchCategories(category)
	if err != nil {
		log.Panicf("Error getting categories: %s", err)
	}

	log.Printf("Categories: %#v\n", categories)

	for _, product := range products {
		if product.SKU == "" || product.SKU == "NE" {
			log.Printf("Skipped %s because no SKU was provided", product.Name)
			continue
		}

		var images []woocommerce.Image
		if product.ImageURL != "" {
			images = []woocommerce.Image{
				{Source: product.ImageURL},
			}
		}

		err := w.CreateProduct(woocommerce.Product{
			Name:        fmt.Sprintf("%s %s %s %s", product.Name, product.AtetCode, product.OmniaCode, product.SKU),
			Description: product.Description,
			Price:       product.Price,
			Categories:  categories,
			Images:      images,
		})

		if err != nil {
			log.Printf("Error adding %s with image %s: %s\n", product.Name, product.ImageURL, err)
		} else {
			log.Printf("Added %s\n", product.Name)
		}
	}
}
