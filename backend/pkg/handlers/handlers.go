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
		resp.ErrorResponse(0, "")
		resp.Write(rw)
		return
	}

	req := service.Product{}
	err = json.Unmarshal(data, &req)
	json.NewEncoder(rw).Encode(req)
	if err != nil {
		log.Println(fmt.Errorf("Handlers create product ошибка декодирования json: %w", err))
		resp.ErrorResponse(http.StatusBadRequest, "Ой что-то сломалось")
		resp.Write(rw)
		return
	}

	resp = service.CreateProduct(req)
	resp.Write(rw)
}

func ReadProduct(rw http.ResponseWriter, r *http.Request) {
	resp := service.ResponseReadProduct{}
	log.Println("REad product handlers")
	vars := mux.Vars(r)
	srtID := vars["id"]
	id, err := strconv.Atoi(srtID)
	if err != nil {
		resp.ErrorResponse(http.StatusBadRequest, "Не найден ID")
		resp.Write(rw)
		return
	}

	resp = service.ReadProduct(service.RequestReadProduct{ID: id})
	resp.Write(rw)
}

func ReadAllProduct(rw http.ResponseWriter, r *http.Request) {
	resp := service.ReadAllProduct()
	resp.Write(rw)
}
