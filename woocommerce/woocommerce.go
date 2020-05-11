package woocommerce

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
)

const apiUrl = "https://tecnoricambi.rs/wp-json/wc/v3/"

type Woocommerce struct {
	key    string
	secret string
}

func New(key, secret string) *Woocommerce {
	return &Woocommerce{
		key:    key,
		secret: secret,
	}
}

func (w *Woocommerce) request(method, endpoint string, data interface{}) (*http.Request, error) {
	body := bytes.NewBuffer(nil)
	if err := json.NewEncoder(body).Encode(data); err != nil {
		return nil, err
	}

	query := url.Values{
		"consumer_key":    []string{w.key},
		"consumer_secret": []string{w.secret},
	}

	req, err := http.NewRequest(method, apiUrl+endpoint+"?"+query.Encode(), body)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (w *Woocommerce) Search(query string) ([]Product, error) {
	type productsRequest struct {
		Search string `json:"search"`
	}

	type productsResponse = []Product

	req, err := w.request("GET", "products", productsRequest{Search: query})
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var products productsResponse
	if err := json.NewDecoder(res.Body).Decode(&products); err != nil {
		return nil, err
	}
	return products, nil
}

type Product struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
