package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"electronic/pkg/service"
)

func CreateProduct(rw http.ResponseWriter, r *http.Request) {
	resp := service.ResponseCreateProduct{}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		resp.ErrorResponse(0, "")
		resp.Write(rw)
		return
	}

	req := service.RequestCreateProduct{}
	err = json.Unmarshal(data, &req)

	if err != nil {
		resp.ErrorResponse(http.StatusBadRequest, "Ой что-то сломалось")
		resp.Write(rw)
		return
	}

	resp = service.CreateProduct(req)
	resp.Write(rw)
}

func ReadProduct(rw http.ResponseWriter, r *http.Request) {
	resp := service.ResponseReadProduct{}
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
