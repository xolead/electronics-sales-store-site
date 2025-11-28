package models

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Request struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Response struct {
	Code    int       `json:"code"`
	Message string    `json:"message"`
	ID      uuid.UUID `json:"id"`
	Admin   bool      `json:"admin"`
}

type RUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func SendResponse(c *gin.Context, code int, message string, id uuid.UUID, admin bool) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
		ID:      id,
		Admin:   admin,
	})
}

func SendBadRequest(c *gin.Context) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    http.StatusBadRequest,
		Message: "Ошибка запроса, невозможно прочитать данные",
		ID:      uuid.Nil,
	})
}

func SendInternalServerError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, Response{
		Code:    http.StatusInternalServerError,
		Message: "Неизвестная ошибка сервера",
		ID:      uuid.Nil,
	})
}
