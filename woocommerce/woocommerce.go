package woocommerce

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const apiUrl = "https://tecnoricambi.rs/wp-json/wc/v3/"

type Woocommerce struct {
	key    string
	secret string
}

func New() *Woocommerce {
	return &Woocommerce{
		key:    "ck_39d1148983c7c666dcd450085af59d343970b38d",
		secret: "cs_1b05b868693a739fd54786b2667f3c1283af03e5",
	}
}

func (w *Woocommerce) request(
	method, endpoint string,
	values url.Values,
	body io.Reader,
) (*http.Request, error) {
	query := url.Values{
		"consumer_key":    []string{w.key},
		"consumer_secret": []string{w.secret},
	}

	for k, v := range values {
		query[k] = v
	}

	req, err := http.NewRequest(method, apiUrl+endpoint+"?"+query.Encode(), body)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (w *Woocommerce) SearchProducts(query string) (result []Product, err error) {
	err = w.search("products", query, result)
	return
}

func (w *Woocommerce) SearchCategories(query string) (result []Category, err error) {
	err = w.search("categories", query, result)
	return
}

func (w *Woocommerce) search(endpoint, query string, result interface{}) error {
	req, err := w.request("GET", endpoint, url.Values{"search": []string{query}}, nil)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return err
	}

	return nil
}

func (w *Woocommerce) CreateProduct(product Product) error {
	body := &bytes.Buffer{}
	if err := json.NewEncoder(body).Encode(product); err != nil {
		return err
	}

	log.Printf("%s\n", body)
	req, err := w.request("POST", "products", url.Values{}, body)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 201 {
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("error reading response body: %s", err)
		}

		return fmt.Errorf("bad response: %d - %s - %s", res.StatusCode, res.Status, b)
	}

	return nil
}

type Product struct {
	Name        string     `json:"name"`
	Description string     `json:"short_description"`
	SKU         string     `json:"sku"`
	Price       string     `json:"price"`
	Categories  []Category `json:"categories"`
	Images      []Image    `json:"images"`
}

type Category struct {
	ID int `json:"id"`
}

type Image struct {
	Source string `json:"src"`
}
