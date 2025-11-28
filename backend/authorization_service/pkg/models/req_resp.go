package models

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Request struct {
	ID    string `json:"id"`
	Admin bool   `json:"admin"`
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Tokens struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

type ResponseUUID struct {
	Response
	ID string `json:"id"`
}

type ResponseTokens struct {
	Response
	Tokens
}

type ResponseUUIDAdmin struct {
	Response
	ID    string `json:"id"`
	Admin bool   `json:"admin"`
}

func SendResponseGetUUID(c *gin.Context, id string, admin bool) {
	c.JSON(http.StatusOK, ResponseUUIDAdmin{
		Response: Response{
			Code:    http.StatusOK,
			Message: "Сессия пользователя активна",
		},
		ID:    id,
		Admin: admin,
	})
}

func SendResponseUUID(c *gin.Context, id string) {
	c.JSON(http.StatusOK, ResponseUUID{
		Response: Response{
			Code:    http.StatusOK,
			Message: "Сессия пользователя активна",
		},
		ID: id,
	})
}

func SendResponse(c *gin.Context, code int, message string) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
	})
}

func SendTokens(c *gin.Context, access, refresh string) {
	c.JSON(http.StatusOK, ResponseTokens{
		Response: Response{
			Code:    http.StatusOK,
			Message: "Токены успешно созданы",
		},
		Tokens: Tokens{
			Access:  access,
			Refresh: refresh,
		},
	})
}

func SendBadRequest(c *gin.Context) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    http.StatusBadRequest,
		Message: "Ошибка запроса, невозможно прочитать данные",
	})
}

func SendInternalServerError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, Response{
		Code:    http.StatusInternalServerError,
		Message: "Неизвестная ошибка сервера",
	})
}
