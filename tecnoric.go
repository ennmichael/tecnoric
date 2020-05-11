package tecnoric

type Product struct {
	Code          string   `json:"code"`
	OriginalCodes []string `json:"original_codes"`
	Description   string   `json:"description"`
	ImageURL      string   `json:"image_url"`
}
