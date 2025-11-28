package handlers

import (
	"manage-service/pkg/communication"
	"manage-service/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func init() {
	log.Logger = log.With().Str("package", "handlers").Logger()
}

func Registration(c *gin.Context) {
	user := models.User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		models.SendBadRequest(c)
		log.Error().Msgf("Ошибка парсинга данных из json: %v", err)
		return
	}

	GUID, admin, err := communication.RegistrationRequest(c, user)
	if err != nil {
		if err == communication.LoginBusy {
			models.SendResponse(c, http.StatusConflict, err.Error())
			return
		}
		models.SendInternalServerError(c)
		log.Error().Msgf("ошибка communication: %v", err)
		return
	}

	tokens, err := communication.AuthorizationRequest(GUID, admin)
	if err != nil {
		models.SendInternalServerError(c)
		log.Error().Msgf("ошибка communication: %v", err)
		return
	}

	c.SetCookie("refreshToken", tokens.Refresh, 60*60*24*3, "/", "localhost", true, true)

	models.SendAccess(c, tokens.Access)

}

func Login(c *gin.Context) {
	user := models.User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		models.SendBadRequest(c)
		log.Error().Msgf("Ошибка парсинга данных из json: %v", err)
		return
	}

	GUID, admin, err := communication.LoginRequest(user)
	if err != nil {
		if err == communication.LoginBusy {
			models.SendResponse(c, http.StatusConflict, err.Error())
			return
		}
		models.SendInternalServerError(c)
		log.Error().Msgf("ошибка communication: %v", err)
		return
	}

	tokens, err := communication.AuthorizationRequest(GUID, admin)
	if err != nil {
		models.SendInternalServerError(c)
		log.Error().Msgf("ошибка communication: %v", err)
		return
	}

	c.SetCookie("refreshToken",
		tokens.Refresh,
		60*60*24*3,
		"/",
		"localhost",
		true,
		true)

	models.SendAccess(c, tokens.Access)

}

func Logout(c *gin.Context) {
	access, ok := c.Get("access")
	if !ok {
		models.SendBadRequest(c)
		return
	}
	c.SetCookie(
		"refresh_token",
		"",
		-1, // удаляем куку
		"/",
		"localhost",
		true,
		true,
	)

	if strAccess, ok := access.(string); ok {
		err := communication.LogoutRequest(strAccess)
		if err != nil {
			models.SendInternalServerError(c)
			return
		}

	}

	models.SendInternalServerError(c)
}

func GetAllProduct(c *gin.Context) {

}
