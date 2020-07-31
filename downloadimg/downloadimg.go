package main

import (
	"encoding/json"
	"flag"
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

	finalProducts := tecnoric.LoadFinalProducts(*inputFileName)
	downloadImages(*imagesDir, finalProducts)
	saveFinalProducts(*outputFileName, finalProducts)
}

func downloadImages(outputDir string, finalProducts []tecnoric.FinalProduct) {
	for k, product := range finalProducts {
		time.Sleep(200 * time.Millisecond)

		u, err := url.Parse(product.ImageURL)
		if err != nil {
			log.Panicf("Error parsing image URL %s: %s", product.ImageURL, err)
		}

		if product.ImageURL != "" {
			split := strings.Split(u.Path, "/")
			localName := split[len(split)-1]
			localPath := path.Join(outputDir, localName)
			log.Printf("Downloading %s to %s", product.ImageURL, localPath)
			if download(product.ImageURL, localPath) {
				product.ImageURL = "https://www.tecnoricambi.rs/scraped-images/" + localName
			} else {
				product.ImageURL = ""
			}

			finalProducts[k] = product
		}
	}
}

func download(url string, outputFileName string) bool {
	res, err := http.DefaultClient.Get(url)
	if err != nil {
		log.Panicf("Error GETting %s: %s", url, err)
	}

	contentType := strings.Split(res.Header.Get("Content-Type"), "/")
	if contentType[0] != "image" {
		return false
	}

	if path.Ext(outputFileName) == "" {
		outputFileName += "." + strings.ToLower(contentType[1])
	}

	f, err := os.Create(outputFileName)
	if err != nil {
		log.Panicf("Error opening %s: %s\n", outputFileName, err)
	}

	defer res.Body.Close()
	defer f.Close()

	io.Copy(f, res.Body)

	return true
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
