package models

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ProductResponse struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Parameters  string   `json:"parameters"`
	Count       int      `json:"count"`
	Price       int      `json:"price"`
	Images      []string `json:"images"`
}

type ResponseCreateProduct struct {
	Response
	URLs []string `json:"urls"`
}

type ResponseReadProduct struct {
	Response
	ProductResponse
}

type ResponseReadAllProduct struct {
	Response
	Products []ProductResponse
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
