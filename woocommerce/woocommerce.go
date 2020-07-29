package woocommerce

import (
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

func (w *Woocommerce) request(method, endpoint string, values url.Values) (*http.Request, error) {
	query := url.Values{
		"consumer_key":    []string{w.key},
		"consumer_secret": []string{w.secret},
	}

	for k, v := range values {
		query[k] = v
	}

	req, err := http.NewRequest(method, apiUrl+endpoint+"?"+query.Encode(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (w *Woocommerce) Search(query string) ([]Product, error) {
	req, err := w.request("GET", "products", url.Values{"search": []string{query}})
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var products []Product
	if err := json.NewDecoder(res.Body).Decode(&products); err != nil {
		return nil, err
	}
	return products, nil
}

type Product struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
