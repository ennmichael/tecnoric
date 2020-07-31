package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"log"
	"os"
	"strings"
	"tecnoric"
	"tecnoric/woocommerce"
)

func main() {
	atetFileName := flag.String("atet", "", "the atet scrape file")
	omniaFileName := flag.String("omnia", "", "the omnia scrape file")
	outputFileName := flag.String("output", "", "the output file name")
	flag.Parse()

	atetProducts := loadItems(*atetFileName)
	omniaProducts := loadItems(*omniaFileName)

	outputFile, err := os.Create(*outputFileName)
	if err != nil {
		log.Fatalf("error creating output file %s: %v", *outputFileName, err)
	}

	defer func() {
		if err := outputFile.Close(); err != nil {
			log.Fatalf("error closing output file %s: %v", *outputFileName, err)
		}
	}()

	joinedProducts := joinProducts(atetProducts, omniaProducts)
	if err := json.NewEncoder(outputFile).Encode(joinedProducts); err != nil {
		log.Fatalf("error encoding to JSON for %#v: %v", joinedProducts, err)
	}

	log.Printf("output written to %s", *outputFileName)
}

func loadItems(fileName string) []tecnoric.Product {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("error creating file %s: %v", fileName, err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Fatalf("error closing file %s: %v", fileName, err)
		}
	}()

	var result []tecnoric.Product
	if err := json.NewDecoder(file).Decode(&result); err != nil {
		log.Fatalf("error decoding JSON from file %s: %v", fileName, err)
	}

	return result
}

type productSearches = map[string][]woocommerce.Product

func searchProducts(woo *woocommerce.Woocommerce, products []tecnoric.Product) productSearches {
	result := productSearches{}

	for _, product := range products {
		log.Printf("searching for %s\n", product.Code)
		search, err := woo.SearchProducts(product.Code)
		if err != nil {
			log.Panicf("error searching for %s: %v", product.Code, err)
		}
		result[product.Code] = search
	}

	return result
}

func joinProducts(atetProducts, omniaProducts []tecnoric.Product) []tecnoric.JoinedProduct {
	var result []tecnoric.JoinedProduct
	woo := woocommerce.New()
	var allProducts []tecnoric.Product
	allProducts = append(allProducts, atetProducts...)
	allProducts = append(allProducts, omniaProducts...)
	searches := searchProducts(woo, allProducts)

omniaLoop:
	for _, omniaProduct := range omniaProducts {
		omniaSearch := searches[omniaProduct.Code]

		if len(omniaSearch) == 1 {
			log.Printf("skipped %s, already available\n", omniaProduct.Code)
			continue
		}

		if len(omniaSearch) > 1 {
			log.Printf("multiple hits for omnia product %s: %#v. Enter to continue.\n", omniaProduct.Code, omniaSearch)
			_, _, err := bufio.NewReader(os.Stdin).ReadLine()
			if err != nil {
				log.Fatalf("error while reading input")
			}

			continue
		}

		for _, atetProduct := range atetProducts {
			if originalCodesMatch(atetProduct, omniaProduct) {
				result = append(result, tecnoric.JoinedProduct{
					OmniaCode:   omniaProduct.Code,
					AtetCode:    atetProduct.Code,
					Description: omniaProduct.Description,
					ImageURL:    atetProduct.ImageURL,
				})
				log.Printf("matched %s with %s\n", omniaProduct.Code, atetProduct.Code)
				continue omniaLoop
			}
		}

		// This omnia product was unmatched.
		result = append(result, tecnoric.JoinedProduct{
			OmniaCode:   omniaProduct.Code,
			AtetCode:    "",
			Description: omniaProduct.Description,
			ImageURL:    omniaProduct.ImageURL,
		})
	}

	return addUnmatchedAtetProducts(atetProducts, result, searches)
}

func originalCodesMatch(atetProduct, omniaProduct tecnoric.Product) bool {
	for _, atetCode := range atetProduct.OriginalCodes {
		for _, omniaCode := range omniaProduct.OriginalCodes {
			if strings.ToLower(atetCode) == strings.ToLower(omniaCode) {
				return true
			}
		}
	}
	return false
}

func addUnmatchedAtetProducts(
	atetProducts []tecnoric.Product,
	result []tecnoric.JoinedProduct,
	searches productSearches) []tecnoric.JoinedProduct {
	for _, atetProduct := range atetProducts {
		productWasMatched := func() bool {
			for _, joinedProduct := range result {
				if atetProduct.Code == joinedProduct.AtetCode {
					return true
				}
			}
			return false
		}()

		if productWasMatched {
			continue
		}

		if len(searches[atetProduct.Code]) > 0 {
			log.Printf("atet product %s already available\n", atetProduct.Code)
			continue
		}

		result = append(result, tecnoric.JoinedProduct{
			OmniaCode:   "",
			AtetCode:    atetProduct.Code,
			Description: atetProduct.Description,
			ImageURL:    atetProduct.ImageURL,
		})
	}

	return result
}
