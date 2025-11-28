package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"core-service/pkg/models"
	"core-service/pkg/service"
)

func CreateProduct(rw http.ResponseWriter, r *http.Request) {
	resp := models.ResponseCreateProduct{}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(fmt.Errorf("Handler CreateProduct ошибка чтения данных: %w", err))
		resp.Error(http.StatusBadRequest, "Ошибка чтения данных")
		resp.Write(rw)
		return
	}

	req := models.ProductResponse{}
	err = json.Unmarshal(data, &req)
	if err != nil {
		log.Println(fmt.Errorf("Handlers create product ошибка декодирования json: %w", err))
		resp.Error(http.StatusBadRequest, "Ошибка чтения json")
		resp.Write(rw)
		return
	}

	resp = service.CreateProduct(req)
	log.Println("RESP CREATE ", resp, "\n", "REQ CREATE ", req)

	resp.Write(rw)
}

func ChangeCountProduct(rw http.ResponseWriter, r *http.Request) {

	resp := models.Response{}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(fmt.Errorf("Handler ChangeCountProduct: %w", err))
		resp.Error(http.StatusBadRequest, "Ошибка чтения данных")
		resp.Write(rw)
		return
	}
	req := models.RequestChangeCount{}
	if err = json.Unmarshal(data, &req); err != nil {
		resp.Error(http.StatusBadRequest, "Ошибка чтения json")
		resp.Write(rw)
		return
	}

	resp = service.ChangeCountProduct(req)
	resp.Write(rw)
}

func ReadProduct(rw http.ResponseWriter, r *http.Request) {
	resp := models.ResponseReadProduct{}
	vars := mux.Vars(r)
	strID := vars["id"]
	id, err := strconv.Atoi(strID)
	if err != nil {
		resp.Error(http.StatusBadRequest, "Не найден ID в запросе")
		resp.Write(rw)
		return
	}

	resp = service.ReadProduct(id)
	resp.Write(rw)
}

func ReadAllProduct(rw http.ResponseWriter, r *http.Request) {
	resp := service.ReadAllProduct()
	resp.Write(rw)
}

func DeleteProduct(rw http.ResponseWriter, r *http.Request) {
	resp := models.Response{}
	vars := mux.Vars(r)
	strID := vars["id"]
	id, err := strconv.Atoi(strID)
	if err != nil {
		resp.Error(http.StatusBadRequest, "Не найден ID в запросе")
		resp.Write(rw)
		return
	}

	resp = service.DeleteProduct(id)
	resp.Write(rw)
}

func HealthCheck(rw http.ResponseWriter, r *http.Request) {
	resp := models.Response{}
	resp.StatusOK()
	resp.Write(rw)
}
