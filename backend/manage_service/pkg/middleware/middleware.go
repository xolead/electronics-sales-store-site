package middleware

import (
	"manage-service/pkg/communication"
	"manage-service/pkg/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.Response{
				Code:    http.StatusUnauthorized,
				Message: "Не найден заголовок Authorization",
			})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, models.Response{
				Code:    http.StatusUnauthorized,
				Message: "Не найден токен в заголовке Authorization Bearer",
			})
			c.Abort()
			return
		}
		access := parts[1]

		refresh, err := c.Cookie("refreshToken")
		if err != nil {
			models.SendInternalServerError(c)
			c.Abort()
			return
		}

		GUID, admin, err := communication.GUIDRequest(access)
		if err != nil {
			tokens, err := communication.RefreshRequest(access, refresh)
			if err != nil {
				models.SendResponse(c, http.StatusUnauthorized, "Токены не действительный")
			}
			c.SetCookie("refreshToken",
				tokens.Refresh,
				60*60*24*3,
				"/",
				"localhost",
				true,
				true)
			models.SendAccess(c, http.StatusContinue, tokens.Access)
			c.Abort()
			return
		}

		c.Set("GUID", GUID)
		c.Set("admin", admin)
		c.Set("access", access)
		c.Set("refresh", refresh)
		c.Next()
	}
}
