package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/mahditd/zarrine-baft-backend/internal/utils"
)

func AuthMiddleware(secret string) gin.HandlerFunc {

	return func(ctx *gin.Context) {

		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header missing",
			})
			ctx.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {

			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization format",
			})
			ctx.Abort()
			return
		}

		tokenString := parts[1]

		token, err := jwt.ParseWithClaims(
			tokenString,
			&utils.JWTClaims{},
			func(token *jwt.Token) (interface{}, error) {

				return []byte(secret), nil
			},
		)

		if err != nil || !token.Valid {

			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			ctx.Abort()
			return
		}

		claims, ok := token.Claims.(*utils.JWTClaims)

		if !ok {

			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token claims",
			})
			ctx.Abort()
			return
		}

		ctx.Set("user_id", claims.UserID)
		ctx.Set("role", claims.Role)

		ctx.Next()
	}
}

func RequireRole(role string) gin.HandlerFunc {

	return func(ctx *gin.Context) {

		userRole, exists := ctx.Get("role")

		if !exists {
			ctx.JSON(401, gin.H{
				"error": "unauthorized",
			})
			ctx.Abort()
			return
		}

		if userRole != role {
			ctx.JSON(403, gin.H{
				"error": "forbidden",
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
