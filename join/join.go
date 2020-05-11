package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"tecnoric"
	"tecnoric/woocommerce"
)

func main() {
	atetFileName := flag.String("atet", "", "the atet scrape file")
	omniaFileName := flag.String("omnia", "", "the omnia scrape file")
	outputFileName := flag.String("output", "",
		"the output file name, also a file with the same name and"+
			" a '-partial' suffix will be written for articles with only one of the codes")
	flag.Parse()

	atetItems := loadItems(*atetFileName)
	omniaItems := loadItems(*omniaFileName)

	woo := woocommerce.New(
		"ck_39d1148983c7c666dcd450085af59d343970b38d",
		"cs_1b05b868693a739fd54786b2667f3c1283af03e5")

	print(outputFileName, atetItems, omniaItems, woo)
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
