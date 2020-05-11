package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gocolly/colly/v2"
	"log"
	"os"
	"tecnoric"
	"tecnoric/utils"
	"time"
)

func main() {
	var targetURL string
	var outputFileName string

	flag.StringVar(&targetURL, "url", "", "the URL to begin scraping from")
	flag.StringVar(&outputFileName, "output", "atet.json", "output file name")
	flag.Parse()

	if targetURL == "" {
		log.Fatalln("error: no target URL")
	}

	outputFile, err := os.Create(outputFileName)
	if err != nil {
		log.Fatalf("error creating or opening file %v: %v\n", outputFileName, err)
	}

	defer func() {
		if err := outputFile.Close(); err != nil {
			log.Fatalf("error closing the output file: %v", err)
		}
	}()

	items, err := scrapeATET(targetURL)
	if err != nil {
		log.Fatalln(err)
	}

	jsonEncoder := json.NewEncoder(outputFile)
	err = jsonEncoder.Encode(items)
	if err != nil {
		log.Fatalf("error encoding json: %v\n", err)
	}

	log.Printf("output written to %v\n", outputFileName)
}

func scrapeATET(targetURL string) ([]tecnoric.Product, error) {
	var result []tecnoric.Product

	const domain = "atet-ricambi.it"

	c := colly.NewCollector(
		colly.AllowedDomains(domain))

	err := c.Limit(&colly.LimitRule{
		DomainRegexp: domain,
		RandomDelay:  time.Second * 2,
		Parallelism:  1,
	})

	if err != nil {
		log.Fatalf("error setting up collector limits: %v", err)
	}

	c.OnRequest(func(request *colly.Request) {
		log.Printf("scraping %v.\n", request.URL)
	})

	c.OnHTML(".ProductExt", func(e *colly.HTMLElement) {
		item := tecnoric.Product{
			Code:          extractValue(e.ChildText(".inner h5")),
			OriginalCodes: scrapeOriginalCodes(e),
			Description:   extractValue(e.ChildText(".inner > p.ProductNotes.ProductTitle")),
			ImageURL:      e.Request.AbsoluteURL(e.ChildAttr("img", "src")),
		}
		result = append(result, item)

		log.Printf("scraped item %#v.\n", item)
	})

	c.OnHTML("ul.pagination > li:last-child > a", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		if err := e.Request.Visit(href); err != nil {
			log.Printf("already visited %v, skipping\n", href)
		}
	})

	if err := c.Visit(targetURL); err != nil {
		return nil, fmt.Errorf("error visiting target URL: %v", err)
	}

	c.Wait()

	return result, nil
}

func extractValue(text string) string {
	split := utils.SplitAndTrim(text, ":")
	if len(split) < 2 {
		return ""
	}
	return split[1]
}

func scrapeOriginalCodes(e *colly.HTMLElement) []string {
	var result []string
	for _, childText := range e.ChildTexts(".inner div .ProductNotes") {
		codes := extractValue(childText)
		if codes == "" {
			continue
		}

		for _, code := range utils.SplitAndTrim(codes, "-") {
			result = append(result, code)
		}
	}
	return result
}
