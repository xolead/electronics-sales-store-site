package handlers

import (
	"authoriz-service/pkg/auth"
	"authoriz-service/pkg/dbwork"
	"authoriz-service/pkg/models"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

func (handler *Handler) Authorization(c *gin.Context) {
	req := models.Request{}
	if err := c.ShouldBindJSON(&req); err != nil {
		models.SendInternalServerError(c)
		log.Error().Msgf("Ошибка парсинга данных из json: %v", err)
		return
	}

	access, err := auth.CreateAccessToken(handler.db, req.ID, req.Admin)
	if err != nil {
		models.SendInternalServerError(c)
		log.Error().Msgf("Ошибка создания access: %v", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	refresh, err := auth.CreateRefreshToken(ctx, handler.db, req.ID)
	if err != nil {
		models.SendInternalServerError(c)
		log.Error().Msgf("Ошибка создания refresh: %v", err)
		return
	}

	models.SendTokens(c, access, refresh)
}

func (handler *Handler) Refresh(c *gin.Context) {
	tokens := models.Tokens{}

	if err := c.ShouldBindJSON(&tokens); err != nil {
		models.SendBadRequest(c)
		log.Error().Msgf("Ошибка получения токенов из запроса: %v", err)
		return
	}

	GUID, admin, err := auth.CheckAccessToken(tokens.Access)
	if err != nil {
		models.SendInternalServerError(c)
		log.Error().Msgf("Ошибка проверки access: %v", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	if err = handler.db.CheckActiveSession(ctx, GUID); err != nil {
		models.SendInternalServerError(c)
		log.Error().Msgf("Ошибка проверки активности сессии: %v", err)
		return
	}

	if err = handler.db.CheckRefreshToken(ctx, GUID, tokens.Refresh); err != nil {
		models.SendInternalServerError(c)
		log.Error().Msgf("Ошибка проверки refresh: %v", err)
		return
	}

	if err = handler.db.StopWorkerRefresh(ctx, GUID); err != nil {
		models.SendInternalServerError(c)
		log.Error().Msgf("Ошибка остановки refresh: %v", err)
		return
	}

	tokens.Access, err = auth.CreateAccessToken(handler.db, GUID, admin)
	if err != nil {
		models.SendInternalServerError(c)
		log.Error().Msgf("Ошибка создания access: %v", err)
		return
	}

	tokens.Refresh, err = auth.CreateRefreshToken(ctx, handler.db, GUID)
	if err != nil {
		models.SendInternalServerError(c)
		log.Error().Msgf("Ошибка создания refresh: %v", err)
		return
	}

	models.SendTokens(c, tokens.Access, tokens.Refresh)
}

func (handler *Handler) Logout(c *gin.Context) {
	tokens := models.Tokens{}
	if err := c.ShouldBindJSON(&tokens); err != nil {
		models.SendBadRequest(c)
		log.Error().Msgf("Ошибка чтения json: %v", err)
		return
	}

	GUID, _, err := auth.CheckAccessToken(tokens.Access)
	if err != nil {
		models.SendInternalServerError(c)
		log.Error().Msgf("Ошибка проверки access: %v", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = handler.db.StopSession(ctx, GUID)
	if err != nil {
		models.SendInternalServerError(c)
		log.Error().Msgf("Ошибка остановки сессии: %v", err)
		return
	}

	models.SendResponse(c, http.StatusOK, "Пользователь успешно вышел из аккаунта")
}

func (handler *Handler) Admin(c *gin.Context) {
	access := c.Param("access")

	GUID, admin, err := auth.CheckAccessToken(access)
	if err != nil {
		models.SendResponse(c, http.StatusUnauthorized, "Ошибка проверки токена")
		log.Error().Msgf("Ошибка проверки access: %v", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err = handler.db.CheckActiveSession(ctx, GUID); err != nil {
		models.SendInternalServerError(c)
		log.Error().Msgf("Ошибка проверки сессии: %v", err)
		return
	}

	if admin {
		models.SendResponse(c, http.StatusOK, "Уровень администратора подтвержден")
		return
	}

	models.SendResponse(c, http.StatusNotFound, "Уровень администратора не подтвержден")

}

func (handler *Handler) GetUUID(c *gin.Context) {
	access := c.Param("access")

	GUID, admin, err := auth.CheckAccessToken(access)
	if err != nil {
		models.SendResponse(c, http.StatusUnauthorized, "Ошибка проверки токена")
		log.Error().Msgf("Ошибка проверки access: %v", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err = handler.db.CheckActiveSession(ctx, GUID); err != nil {
		models.SendInternalServerError(c)
		log.Error().Msgf("Ошибка проверки сессии: %v", err)
		return
	}
	models.SendResponseGetUUID(c, GUID, admin)
}
