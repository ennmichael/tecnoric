package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strconv"
	"tecnoric"
	"tecnoric/utils"
	"time"
)

const omniaBaseURL = "https://b2b.omniacomponents.com/"

func main() {
	categoryID := flag.String("category", "", "the category to scrape")
	outputFileName := flag.String("output", "omnia.json", "output file name")
	flag.Parse()

	categoryIDInt, err := strconv.Atoi(*categoryID)
	if err != nil {
		log.Fatalf("error converting %s to integer: %v", *categoryID, err)
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("error creating a new CookieJar: %v", err)
	}

	client := &http.Client{
		Jar: jar,
	}

	outputFile, err := os.Create(*outputFileName)
	if err != nil {
		log.Fatalf("error creating the output file: %v", err)
	}

	defer func() {
		if err := outputFile.Close(); err != nil {
			log.Fatalf("error closing the output file: %v", err)
		}
	}()

	logIntoOmnia(client)
	items := getItems(client, categoryIDInt)

	jsonEncoder := json.NewEncoder(outputFile)
	err = jsonEncoder.Encode(items)
	if err != nil {
		log.Fatalf("error encoding json: %v", err)
	}

	log.Printf("output written to %v\n", *outputFileName)
}

func logIntoOmnia(c *http.Client) {
	type loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	type loginResponse struct {
		User *struct{} `json:"user"`
	}

	res, err := c.Do(omniaJSONRequest("POST", "login", loginRequest{
		Username: "Generic Customer",
		Password: "gen_cust_2019",
	}))

	if err != nil {
		log.Fatalf("error logging into the website: %v", err)
	}

	var response loginResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		log.Fatalf("error decoding the login response: %v", err)
	}

	if response.User == nil {
		log.Fatalf("login failed, the credentials are likely incorrect")
	}

	log.Println("logged in successfully")
}

func omniaJSONRequest(method, endpoint string, data interface{}) *http.Request {
	body := bytes.NewBuffer(nil)
	if err := json.NewEncoder(body).Encode(data); err != nil {
		log.Fatalf("error encoding JSON request %v: %v", data, err)
	}

	r, err := http.NewRequest(method, omniaBaseURL+endpoint, body)
	if err != nil {
		log.Fatalf("error sending JSON request %#v: %v", body.String(), err)
	}

	addHeaders(r)
	return r
}

func addHeaders(r *http.Request) {
	r.Header.Add("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:76.0) Gecko/20100101 Firefox/76.0")
	r.Header.Add("Content-Type", "application/json")
}

func getItems(c *http.Client, categoryID int) []tecnoric.Product {
	var result []tecnoric.Product

	for i := 1; ; i++ {
		productList := getProductList(c, categoryID, i)

		if len(productList.Products) == 0 {
			return result
		}

		for _, product := range productList.Products {
			randomDelay()

			technicalDetails := getTechnicalDetails(c, product.ID)
			var originalCodes []string
			technicalDescription := ""
			if technicalDetails != nil && technicalDetails.OriginalCodes != nil {
				originalCodes = utils.SplitAndTrim(*technicalDetails.OriginalCodes, ",")
			}
			if technicalDetails != nil && technicalDetails.TechnicalDescription != nil {
				technicalDescription = *technicalDetails.TechnicalDescription
			}

			item := tecnoric.Product{
				Code:          product.Code,
				Description:   product.Name + technicalDescription,
				ImageURL:      product.Image,
				OriginalCodes: originalCodes,
			}
			result = append(result, item)

			log.Printf("scraped item %#v.\n", item)
		}
	}
}

func randomDelay() {
	n := time.Duration(rand.Int63n(1000) + 100)
	time.Sleep(n * time.Millisecond)
}

type productListRequest struct {
	CategoryID     int       `json:"category_id"`
	DivisionID     string    `json:"division_id"`
	OnlyAvailable  *struct{} `json:"onlyAvailable"` // Always nil
	OrderBy        string    `json:"orderBy"`
	PageIndex      int       `json:"page_index"`
	PageSize       int       `json:"page_size"`
	SelectedFacets string    `json:"selected_facets"`
	UserSearch     string    `json:"user_search"`
}

type productListResponse struct {
	Products []product `json:"products"`
}

type product struct {
	ID    int    `json:"id"`
	Code  string `json:"code"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

func getProductList(c *http.Client, categoryID int, index int) *productListResponse {
	res, err := c.Do(omniaJSONRequest("POST", "api/v1/public/get_productlist", &productListRequest{
		CategoryID:     categoryID,
		DivisionID:     "1",
		OnlyAvailable:  nil,
		OrderBy:        "price asc",
		PageIndex:      index,
		PageSize:       20,
		SelectedFacets: "",
		UserSearch:     "",
	}))

	if err != nil {
		log.Fatalf("error getting the products list: %v", err)
	}

	productListResponse := &productListResponse{}
	if err := json.NewDecoder(res.Body).Decode(&productListResponse); err != nil {
		log.Fatalf("error decoding the products list: %v", err)
	}
	return productListResponse
}

type techsheetRequest struct {
	ProductID string     `json:"product_id"`
	Filter    []struct{} `json:"filter"`
}

type technicalDetails struct {
	OriginalCodes        *string `json:"cross_reference_customer"`
	TechnicalDescription *string `json:"technical_description"`
}

type techsheetData struct {
	General []technicalDetails `json:"dati_generali"`
}

type techsheetResponse struct {
	Data techsheetData `json:"data"`
}

func getTechnicalDetails(c *http.Client, productID int) *technicalDetails {
	res, err := c.Do(omniaJSONRequest("POST", "api/v1/public/get_techsheet_data", techsheetRequest{
		ProductID: strconv.Itoa(productID),
		Filter:    []struct{}{},
	}))
	if err != nil {
		log.Fatalf("error getting the technical details for %d: %v", productID, err)
	}
	var techsheetRes []techsheetResponse
	if err := json.NewDecoder(res.Body).Decode(&techsheetRes); err != nil {
		log.Fatalf("error decoding the technical details: %v", err)
	}

	if len(techsheetRes) == 0 || len(techsheetRes[0].Data.General) == 0 {
		return nil
	}

	return &techsheetRes[0].Data.General[0]
}
