package tecnoric

type Product struct {
	Code          string   `json:"code"`
	OriginalCodes []string `json:"original_codes"`
	Description   string   `json:"description"`
	ImageURL      string   `json:"image_url"`
}

type JoinedProduct struct {
	OmniaCode   string `json:"omnia_code"`
	AtetCode    string `json:"atet_code"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
}

type PartialProduct struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	AtetCode  *string `json:"atet_code"`
	OmniaCode *string `json:"omnia_code"`
}

type FinalProduct struct {
	Name        string `json:"name"`
	AtetCode    string `json:"atet_code"`
	OmniaCode   string `json:"omnia_code"`
	SKU         string `json:"sku"`
	Price       string `json:"price"`
	ImageURL    string `json:"image_url"`
	Description string `json:"description"`
}
