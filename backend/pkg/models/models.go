package models

type Product struct {
	ID          int            `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  string         `json:"parameters"`
	Count       int            `json:"count"`
	Images      []ProductImage `json:"images"`
}

type ProductImage struct {
	ID        int    `json:"id"`
	ProductID int    `json:"product_id"`
	Name      string `json:"name"`
	Key       string `json:"key"`
}

type S3SImage struct {
	URL    string `json:"url"`
	FileID string `json:"file_id"`
}

type CreateProductResponse struct {
	code    int
	message string
	URLs    []string
}
