package handlers

import (
	"auth-service/pkg/dbwork"
	"auth-service/pkg/models"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func init() {
	log.Logger = log.With().Str("package", "handlers").Logger()
}

type Handler struct {
	db *dbwork.DataBase
}

func NewHandler(db *dbwork.DataBase) *Handler {
	return &Handler{db: db}
}

func (handler *Handler) Registration(c *gin.Context) {
	user := models.RUser{}
	if err := c.ShouldBindJSON(&user); err != nil {
		models.SendBadRequest(c)
		log.Error().Msgf("Ошибка получения данных из json: %v", err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	id, err := handler.db.CreateUser(ctx, user.Login, user.Password)
	if err == dbwork.LoginBusy {
		models.SendResponse(c, http.StatusConflict, err.Error(), uuid.Nil, false)
		return
	}
	if err != nil {
		models.SendInternalServerError(c)
		log.Error().Msgf("Ошибка создания нового пользователя: %v", err)
		return
	}

	models.SendResponse(c, http.StatusCreated, "Пользователь успешно зарегистрирован", id, false)
}

func (handler *Handler) Login(c *gin.Context) {
	user := models.RUser{}
	if err := c.ShouldBindJSON(&user); err != nil {
		models.SendBadRequest(c)
		log.Error().Msgf("Ошибка получения данных из json: %v", err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	id, admin, err := handler.db.VerifyPassword(ctx, user.Login, user.Password)
	if err == dbwork.LoginNotFound {
		models.SendResponse(c, http.StatusNotFound, err.Error(), uuid.Nil, false)
		return
	}
	if err == dbwork.PasswordIsNotCorrect {
		models.SendResponse(c, http.StatusUnauthorized, err.Error(), uuid.Nil, false)
		log.Error().Msgf("Пользователь ввёл неправильный пароль")
		return
	}
	if err != nil {
		models.SendInternalServerError(c)
		log.Error().Msgf("Ошибка проверки пароль: %v", err)
		return
	}

	models.SendResponse(c, http.StatusOK, "Пользователь успешно вошёл", id, admin)
}
