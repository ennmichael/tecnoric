package main

import (
	"flag"
	"log"
	"tecnoric"
	"tecnoric/woocommerce"
)

func main() {
	inputFileName := flag.String("input", "", "input JSON file with final products")
	flag.Parse()

	products := tecnoric.LoadFinalProducts(*inputFileName)
	insertProducts(products)
}

func insertProducts(products []tecnoric.FinalProduct) {
	w := woocommerce.New()
	err := w.CreateProduct(woocommerce.Product{
		Name:        "primer",
		Description: "opis",
		SKU:         "123",
		Price:       "250,00",
		Categories: []woocommerce.Category{
			{ID: 2},
		},
		Images: []woocommerce.Image{
			{Source: "https://www.tecnoricambi.rs/scraped-images/167AC01.jpeg"},
		},
	})

	if err != nil {
		log.Panic(err)
	}
}
