package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	cloudstorage "electronic/pkg/cloud_storage"
	"electronic/pkg/dbwork"
	"electronic/pkg/models"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Product struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Parameters  string   `json:"parameters"`
	Count       int      `json:"count"`
	Images      []string `json:"images"`
}

type ResponseCreateProduct struct {
	Response
	URLs []string `json:"urls"`
}

type ResponseReadProduct struct {
	Response
	Product
}

type ResponseReadAllProduct struct {
	Response
	Products []Product
}

type RequestChangeCount struct {
	ID    int
	Count int
}

func (resp *Response) Write(rw http.ResponseWriter) {
	json.NewEncoder(rw).Encode(resp)
}

func (resp *ResponseReadAllProduct) Write(rw http.ResponseWriter) {
	json.NewEncoder(rw).Encode(resp)
}

func (resp *ResponseCreateProduct) Write(rw http.ResponseWriter) {
	json.NewEncoder(rw).Encode(resp)
}

func (resp *ResponseReadProduct) Write(rw http.ResponseWriter) {
	json.NewEncoder(rw).Encode(resp)
}

func (resp *Response) Error(code int, message string) {
	resp.Code = code
	resp.Message = message
}

func (resp *Response) InternalError() {
	resp.Code = http.StatusInternalServerError
	resp.Message = "Внутренняя ошибка сервера"
}

func (resp *Response) StatusOK() {
	resp.Code = http.StatusOK
	resp.Message = "Выполнение прошло успешно"
}

func (resp *Response) StatusCreated() {
	resp.Code = http.StatusCreated
	resp.Message = "Объект успешно создан"
}

func CreateProduct(product Product) ResponseCreateProduct {
	resp := ResponseCreateProduct{}

	psql, err := dbwork.NewPostgreSQL(dbwork.LoadPSQLConfig())
	if err != nil {
		log.Println(err)
		resp.InternalError()
		return resp
	}

	s3s, err := cloudstorage.NewS3S(cloudstorage.LoadCfg())
	if err != nil {
		log.Println(err)
		resp.InternalError()
		return resp
	}

	productDB := models.Product{
		Description: product.Description,
		Name:        product.Name,
		Count:       product.Count,
		Parameters:  product.Parameters,
		Images:      make([]models.ProductImage, len(product.Images)),
	}
	urls := make([]string, len(product.Images))
	for i, name := range product.Images {
		image, err := s3s.UploadURL(name)
		if err != nil {
			log.Println(err)
			resp.InternalError()
			return resp
		}
		urls = append(urls, image.URL)
		productDB.Images[i].Key = image.FileID
		productDB.Images[i].Name = name
	}

	_, err = psql.CreateProduct(productDB)
	if err != nil {
		log.Println(err)
		resp.InternalError()
		return resp
	}

	resp.URLs = urls
	resp.StatusCreated()
	return resp
}

func ReadProduct(id int) ResponseReadProduct {

	resp := ResponseReadProduct{}

	product := models.Product{}

	s3s, err := cloudstorage.NewS3S(cloudstorage.LoadCfg())
	if err != nil {
		log.Println(err)
		resp.InternalError()
		return resp
	}

	psql, err := dbwork.NewPostgreSQL(dbwork.LoadPSQLConfig())
	if err != nil {
		log.Println(err)
		resp.InternalError()
		return resp
	}

	product, err = psql.ReadProduct(id)
	if err != nil {
		log.Println(err)
		resp.InternalError()
		return resp
	}

	resp.Product, err = createDownloadURLs(product, s3s)
	if err != nil {
		log.Println(err)
		resp.InternalError()
		return resp
	}
	resp.StatusOK()
	return resp

}

func createDownloadURLs(productDB models.Product, s3s cloudstorage.CloudStorage) (Product, error) {

	product := Product{
		ID:          productDB.ID,
		Name:        productDB.Name,
		Description: productDB.Description,
		Parameters:  productDB.Parameters,
		Count:       productDB.Count,
	}

	URLs := make([]string, 0, len(productDB.Images))
	for _, image := range productDB.Images {

		image, err := s3s.DownloadURL(image.Key)
		if err != nil {
			return product, fmt.Errorf("createDownloadURLs ошибка создания url: %w", err)
		}

		URLs = append(URLs, image.FileID)

	}
	product.Images = URLs
	return product, nil
}

func ReadAllProduct() ResponseReadAllProduct {
	resp := ResponseReadAllProduct{}

	s3s, err := cloudstorage.NewS3S(cloudstorage.LoadCfg())
	if err != nil {
		log.Println(err)
		resp.InternalError()
		return resp
	}

	psql, err := dbwork.NewPostgreSQL(dbwork.LoadPSQLConfig())
	if err != nil {
		log.Println(err)
		resp.InternalError()
		return resp
	}

	products, err := psql.ReadListProduct()
	if err != nil {
		log.Println(err)
		resp.InternalError()
		return resp
	}

	for _, product := range products {
		tempProduct, err := createDownloadURLs(product, s3s)
		if err != nil {
			log.Println(err)
			resp.InternalError()
			return resp
		}
		resp.Products = append(resp.Products, tempProduct)
	}

	resp.StatusOK()
	return resp
}

func ChangeCountProduct(req RequestChangeCount) Response {
	resp := Response{}

	psql, err := dbwork.NewPostgreSQL(dbwork.LoadPSQLConfig())
	if err != nil {
		log.Println(err)
		resp.InternalError()
		return resp
	}

	err = psql.ChangeCountProduct(req.ID, req.Count)
	if err != nil {
		log.Println(err)
		resp.InternalError()
		return resp
	}

	resp.StatusOK()
	return resp
}

func DeleteProduct(id int) Response {
	resp := Response{}

	psql, err := dbwork.NewPostgreSQL(dbwork.LoadPSQLConfig())
	if err != nil {
		log.Println(err)
		resp.InternalError()
		return resp
	}

	s3s, err := cloudstorage.NewS3S(cloudstorage.LoadCfg())
	if err != nil {
		log.Println(err)
		resp.InternalError()
		return resp
	}

	keys, err := psql.DeleteProduct(id)
	for _, key := range keys {
		image, err := s3s.DeleteURL(key)
		if err != nil {
			log.Println(err)
			resp.InternalError()
			return resp
		}

		r, err := http.NewRequest("DELETE", image.URL, nil)

		client := http.Client{}

		re, err := client.Do(r)
		if err != nil {
			log.Println(err)
			resp.InternalError()
			return resp
		}

		defer re.Body.Close()

		if re.StatusCode != http.StatusNoContent {
			log.Println(err)
			resp.InternalError()
			return resp
		}

	}

	resp.StatusOK()
	return resp

}
