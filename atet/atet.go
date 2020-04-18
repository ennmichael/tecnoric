package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gocolly/colly/v2"
	"log"
	"os"
	"strings"
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

	items, err := ScrapeATET(targetURL)
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

type item struct {
	Code          string   `json:"code"`
	OriginalCodes []string `json:"original_codes"`
	Description   string   `json:"description"`
	ImageURL      string   `json:"image_url"`
}

func ScrapeATET(targetURL string) ([]item, error) {
	var result []item

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
		item := item{
			Code:          ExtractValue(e.ChildText(".inner h5")),
			OriginalCodes: ScrapeOriginalCodes(e),
			Description:   ExtractValue(e.ChildText(".inner > p.ProductNotes.ProductTitle")),
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

func ExtractValue(text string) string {
	split := SplitAndTrim(text, ":")
	if len(split) < 2 {
		return ""
	}
	return split[1]
}

func ScrapeOriginalCodes(e *colly.HTMLElement) []string {
	var result []string
	for _, childText := range e.ChildTexts(".inner div .ProductNotes") {
		codes := ExtractValue(childText)
		if codes == "" {
			continue
		}

		for _, code := range SplitAndTrim(codes, "-") {
			result = append(result, code)
		}
	}
	return result
}

func SplitAndTrim(s, sep string) []string {
	var result []string
	for _, ss := range strings.Split(s, sep) {
		result = append(result, strings.TrimSpace(ss))
	}
	return result
}
