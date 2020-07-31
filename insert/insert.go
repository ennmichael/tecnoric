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

	for _, product := range products {
		err := w.CreateProduct(woocommerce.Product{
			Name:        fmt.Sprintf("%s %s %s", product.Name, product.AtetCode, product.OmniaCode),
			Description: product.Description,
			SKU:         product.SKU,
			Price:       product.Price,
			Categories:  categories,
			Images: []woocommerce.Image{
				{Source: product.ImageURL},
			},
		})

		if err != nil {
			log.Println(err)
		}
	}
}
