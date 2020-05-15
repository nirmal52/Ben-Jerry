package models

type Product struct {
	Name                  string   `json:"name"`
	ImageClosed           string   `json:"image_closed"`
	ImageOpen             string   `json:"image_open"`
	Description           string   `json:"description"`
	Story                 string   `json:"story"`
	SourcingValues        []string `json:"sourcing_values"`
	Ingredients           []string `json:"ingredients"`
	AllergyInfo           string   `json:"allergy_info"`
	DietaryCertifications string   `json:"dietary_certifications"`
	ProductId             string   `json:"productId"`
}
