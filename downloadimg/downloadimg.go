package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"tecnoric"
	"time"
)

func main() {
	inputFileName := flag.String("input", "", "input file name")
	imagesDir := flag.String("images-dir", "", "directory for downloaded images")
	outputFileName := flag.String("output", "", "output file name")
	flag.Parse()

	finalProducts := loadFinalProducts(*inputFileName)
	downloadImages(*imagesDir, finalProducts)
	saveFinalProducts(*outputFileName, finalProducts)
}

func loadFinalProducts(fileName string) (result []tecnoric.FinalProduct) {
	f, err := os.Open(fileName)
	if err != nil {
		log.Panicf("Error opening file %s: %s", fileName, err)
	}

	defer f.Close()

	if err = json.NewDecoder(f).Decode(&result); err != nil {
		log.Panicf("Error decoding JSON: %s\n", err)
	}

	return
}

func downloadImages(outputDir string, finalProducts []tecnoric.FinalProduct) {
	year := time.Now().Year()
	month := int(time.Now().Month())
	for _, product := range finalProducts {
		time.Sleep(200 * time.Millisecond)

		u, err := url.Parse(product.ImageURL)
		if err != nil {
			log.Panicf("Error parsing image URL %s: %s", product.ImageURL, err)
		}

		if product.ImageURL != "" {
			split := strings.Split(u.Path, "/")
			localName := path.Join(outputDir, split[len(split)-1])
			log.Printf("Downloading %s to %s", product.ImageURL, localName)
			download(product.ImageURL, localName)
			product.ImageURL = fmt.Sprintf(
				"https://www.tecnoricambi.rs/wp-content/uploads/%d/%02d/%s", year, month, localName)
		}
	}
}

func download(url string, outputFileName string) {
	res, err := http.DefaultClient.Get(url)
	f, err := os.Create(outputFileName)
	if err != nil {
		log.Panicf("Error opening %s: %s\n", outputFileName, err)
	}

	defer res.Body.Close()
	defer f.Close()

	io.Copy(f, res.Body)
}

func saveFinalProducts(outputFileName string, products []tecnoric.FinalProduct) {
	f, err := os.Create(outputFileName)
	if err != nil {
		log.Panicf("Error creating file %s: %s", outputFileName, err)
	}

	if err := json.NewEncoder(f).Encode(products); err != nil {
		log.Panicf("Error encoding JSON: %s", err)
	}
}
