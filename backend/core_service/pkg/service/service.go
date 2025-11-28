package service

import (
	"fmt"
	"log"
	"net/http"

	cloudstorage "core-service/pkg/cloud_storage"
	"core-service/pkg/connectionpool"
	"core-service/pkg/dbwork"
	"core-service/pkg/models"
)

func CreateProduct(product models.ProductResponse) models.ResponseCreateProduct {
	resp := models.ResponseCreateProduct{}

	pool := connectionpool.NewConnectionPool()
	s3s := pool.GetS3Storage()
	dataBase := pool.GetDataBase()

	productDB := models.Product{
		Description: product.Description,
		Name:        product.Name,
		Count:       product.Count,
		Parameters:  product.Parameters,
		Price:       product.Price,
		Images:      make([]models.ProductImage, len(product.Images)),
	}
	urls := make([]string, 0, len(product.Images))
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

	_, err := dataBase.CreateProduct(productDB)
	if err != nil {
		log.Println(err)
		resp.InternalError()
		return resp
	}

	resp.URLs = urls
	resp.StatusCreated()
	return resp
}

func ReadProduct(id int) models.ResponseReadProduct {

	resp := models.ResponseReadProduct{}

	product := models.Product{}

	pool := connectionpool.NewConnectionPool()
	dataBase := pool.GetDataBase()

	product, err := dataBase.ReadProduct(id)
	if err != nil {
		log.Println(err)
		resp.InternalError()
		return resp
	}

	s3s := pool.GetS3Storage()

	resp.ProductResponse, err = createDownloadURLs(product, s3s)
	if err != nil {
		log.Println(err)
		resp.InternalError()
		return resp
	}
	resp.StatusOK()
	return resp

}

func createDownloadURLs(
	productDB models.Product,
	s3s cloudstorage.CloudStorage,
) (models.ProductResponse, error) {

	product := models.ProductResponse{
		ID:          productDB.ID,
		Name:        productDB.Name,
		Description: productDB.Description,
		Parameters:  productDB.Parameters,
		Count:       productDB.Count,
		Price:       productDB.Price,
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

func ReadAllProduct() models.ResponseReadAllProduct {
	resp := models.ResponseReadAllProduct{}

	pool := connectionpool.NewConnectionPool()
	s3s := pool.GetS3Storage()
	dataBase := pool.GetDataBase()

	products, err := dataBase.ReadListProduct()
	if err != nil {
		log.Println(err)
		resp.InternalError()
		return resp
	}

	resp.Products = make([]models.ProductResponse, 0, len(products))

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

func ChangeCountProduct(req models.RequestChangeCount) models.Response {
	resp := models.Response{}

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

func DeleteProduct(id int) models.Response {
	resp := models.Response{}

	pool := connectionpool.NewConnectionPool()
	dataBase := pool.GetDataBase()
	s3s := pool.GetS3Storage()

	keys, err := dataBase.DeleteProduct(id)
	if err != nil {
		log.Println(err)
		resp.InternalError()
		return resp
	}

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
