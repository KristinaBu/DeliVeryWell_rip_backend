package handler

import (
	"BMSTU_IU5_53B_rip/internal/app/ds"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"os"
	"strings"
)

const prefix = "Bearer"

// RoleMiddleware проверяет роль пользователя и проверяет блеклист
func (h *Handler) RoleMiddleware(allowedRoles ...ds.User) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractTokenFromHeader(c.Request)
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		userID, ok := claims["user_id"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		if !h.Repository.IsBlackListed(userID, tokenString) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token in blacklist"})
		}

		// установим контекст, чтоб каждый раз не парсить токен
		c.Set("user_id", uint(userID))

		isAdmin, ok := claims["is_admin"].(bool)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		c.Set("is_admin", isAdmin)

		// проверка на роль
		if !isRoleAllowed(isAdmin, allowedRoles) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "forbidden",
				"msg":   "role not allowed",
			})
		}
		// мидлвара успешна, возвращаем сообщение с /protected
		c.Next()
	}
}

func extractTokenFromHeader(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	if bearerToken == "" {
		return ""
	}
	if strings.Split(bearerToken, " ")[0] != prefix {
		return ""
	}
	return strings.Split(bearerToken, " ")[1]
}

func isRoleAllowed(isAdmin bool, allowedRoles []ds.User) bool {
	for _, role := range allowedRoles {
		if isAdmin == role.IsAdmin {
			return true
		}
	}
	return false
}

//
