package models

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Tokens struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

type ResponseUUID struct {
	Response
	UUID string `json:"id"`
}

type ResponseTokens struct {
	Response
	Tokens
}

type ResponseAccess struct {
	Response
	Access string `json:"access"`
}
type ResponseAuth struct {
	ResponseUUID
	Admin bool `json:"admin"`
}

func SendAccess(c *gin.Context, code int, access string) {
	c.JSON(code, ResponseAccess{
		Response: Response{
			Code:    code,
			Message: "Токены устарели и были обновлены",
		},
		Access: access,
	},
	)
}

func SendTokensRefresh(c *gin.Context, access, refresh string) {
	c.JSON(http.StatusContinue, ResponseTokens{
		Response: Response{
			Code:    http.StatusContinue,
			Message: "Токены устарели и были обновлены",
		},
		Tokens: Tokens{
			Access:  access,
			Refresh: refresh,
		},
	},
	)
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
func SendResponseUUID(c *gin.Context, UUID string) {
	c.JSON(http.StatusOK, ResponseUUID{
		Response: Response{
			Code:    http.StatusOK,
			Message: "Сессия пользователя активна",
		},
		UUID: UUID,
	})
}

func SendResponse(c *gin.Context, code int, message string) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
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

func SendParseResponse(c *gin.Context)
