package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"electronic/pkg/service"
)

func CreateProduct(rw http.ResponseWriter, r *http.Request) {
	resp := service.ResponseCreateProduct{}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(fmt.Errorf("Handler CreateProduct ошибка чтения данных: %w", err))
		resp.ErrorResponse(http.StatusBadRequest, "Ошибка чтения данных")
		resp.Write(rw)
		return
	}

	req := service.Product{}
	err = json.Unmarshal(data, &req)
	if err != nil {
		log.Println(fmt.Errorf("Handlers create product ошибка декодирования json: %w", err))
		resp.ErrorResponse(http.StatusBadRequest, "Ой что-то сломалось")
		resp.Write(rw)
		return
	}

	resp = service.CreateProduct(req)
	resp.Write(rw)
}

func ChangeCountProduct(rw http.ResponseWriter, r *http.Request) {

	resp := service.Response{}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(fmt.Errorf("Handler ChangeCountProduct: %w", err))
		resp.ErrorResponse(http.StatusBadRequest, "Ошибка чтения данных")
		resp.Write(rw)
		return
	}
	req := service.RequestChangeCount{}
	if err = json.Unmarshal(data, &req); err != nil {
		resp.ErrorResponse(http.StatusBadRequest, "Ошибка чтения json")
		resp.Write(rw)
		return
	}

	resp = service.ChangeCountProduct(req)
	resp.Write(rw)
}

func ReadProduct(rw http.ResponseWriter, r *http.Request) {
	resp := service.ResponseReadProduct{}
	vars := mux.Vars(r)
	strID := vars["id"]
	id, err := strconv.Atoi(strID)
	if err != nil {
		resp.ErrorResponse(http.StatusBadRequest, "Не найден ID в запросе")
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
	resp := service.Response{}
	vars := mux.Vars(r)
	strID := vars["id"]
	id, err := strconv.Atoi(strID)
	if err != nil {
		resp.ErrorResponse(http.StatusBadRequest, "Не найден ID в запросе")
		resp.Write(rw)
		return
	}

	resp = service.DeleteProduct(id)
	resp.Write(rw)
}
