package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
	"tecnoric"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func main() {
	inputFileName := flag.String("input", "", "input file name")
	outputFileName := flag.String("output", "", "output file name")
	flag.Parse()

	finalProducts := loadProductsFromExcelFile(*inputFileName)
	outputFinalProducts(finalProducts, *outputFileName)
}

func loadProductsFromExcelFile(fileName string) []tecnoric.FinalProduct {
	file, err := excelize.OpenFile(fileName)
	if err != nil {
		log.Fatalf("Error opening excel file %s: %s\n", fileName, err)
	}

	result := []tecnoric.FinalProduct{}
	for i := 2; ; i++ {
		index := strconv.Itoa(i)
		omniaCode := file.GetCellValue("Sheet1", "A"+index)
		atetCode := file.GetCellValue("Sheet1", "B"+index)
		sku := file.GetCellValue("Sheet1", "C"+index)
		price := file.GetCellValue("Sheet1", "D"+index)
		name := file.GetCellValue("Sheet1", "E"+index)
		imageURL := file.GetCellValue("Sheet1", "F"+index)
		description := file.GetCellValue("Sheet1", "G"+index)

		if name == "" {
			return result
		}

		price = strings.ToLower(price)
		if strings.Contains(price, "rsd") {
			price = strings.ReplaceAll(price, "rsd", "")
			price = strings.TrimSpace(price)
		}

		result = append(result, tecnoric.FinalProduct{
			OmniaCode:   omniaCode,
			AtetCode:    atetCode,
			Name:        name,
			SKU:         sku,
			Price:       price,
			ImageURL:    imageURL,
			Description: description,
		})
	}
}

func outputFinalProducts(finalProducts []tecnoric.FinalProduct, outputFileName string) {
	f, err := os.Create(outputFileName)
	if err != nil {
		log.Fatalf("Error opening %s: %s\n", outputFileName, err)
	}

	encoder := json.NewEncoder(f)
	if err = encoder.Encode(finalProducts); err != nil {
		log.Fatalf("Error encoding JSON: %s", err)
	}
}
