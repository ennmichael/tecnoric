package main

import (
	"encoding/json"
	"flag"
	"github.com/360EntSecGroup-Skylar/excelize"
	"os"
	"strconv"
	"tecnoric"
)

func main() {
	inputFileName := flag.String("input", "", "input file name")
	outputFileName := flag.String("output", "", "output file name")
	flag.Parse()

	products := loadProducts(*inputFileName)
	excelFile := createExcelFile(products)

	if err := excelFile.SaveAs(*outputFileName); err != nil {
		panic(err)
	}
}

func loadProducts(fileName string) []tecnoric.JoinedProduct {
	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	var products []tecnoric.JoinedProduct
	if err := json.NewDecoder(f).Decode(&products); err != nil {
		panic(err)
	}

	return products
}

func createExcelFile(products []tecnoric.JoinedProduct) *excelize.File {
	f := excelize.NewFile()
	f.SetCellValue("Sheet1", "A1", "omnia")
	f.SetCellValue("Sheet1", "B1", "atet")
	f.SetCellValue("Sheet1", "C1", "šifra")
	f.SetCellValue("Sheet1", "D1", "cena")
	f.SetCellValue("Sheet1", "E1", "ime")
	f.SetCellValue("Sheet1", "F1", "slika")

	for k, product := range products {
		index := strconv.Itoa(k + 2)
		f.SetCellValue("Sheet1", "A"+index, product.OmniaCode)
		f.SetCellValue("Sheet1", "B"+index, product.AtetCode)
		f.SetCellValue("Sheet1", "E"+index, product.Description)
		f.SetCellValue("Sheet1", "F"+index, product.ImageURL)
	}

	return f
}
