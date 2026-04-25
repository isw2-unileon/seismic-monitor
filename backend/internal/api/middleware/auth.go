package middleware

import (
	"net/http"
	"strings"

	"seismic-monitor/backend/internal/auth"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware verifica el token JWT en el header Authorization
func AuthMiddleware(jwtService *auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "se requiere token de autenticación"})
			c.Abort()
			return
		}

		// El header suele ser "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "formato de token inválido"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token inválido o expirado"})
			c.Abort()
			return
		}

		// Guardar el UserID en el contexto para usarlo en los handlers
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
