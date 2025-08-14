// internal/middlewares/auth.go
package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"accessv2/internal/config"
	"accessv2/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func SignInAuthenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verificar si SECURE está habilitado
		secureEnabled, err := strconv.ParseBool(os.Getenv("SECURE"))
		if err != nil {
			// Si hay error al parsear, asumimos false por seguridad
			secureEnabled = false
		}

		// Si SECURE no está activado, continuamos sin validar
		if !secureEnabled {
			c.Next()
			return
		}

		// Validación original cuando SECURE=true
		incoming := c.GetHeader("X-Auth-Trigger")
		if incoming != config.AuthHeader {
			fmt.Println("Unauthorized access attempt.")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Invalid or missing X-Auth-Trigger",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func CheckJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 0. Verificar si SECURE está habilitado
		secureEnabled, err := strconv.ParseBool(os.Getenv("SECURE"))
		if err != nil {
			// Si hay error al parsear, asumimos false por seguridad
			secureEnabled = false
		}

		// Si SECURE no está activado, continuamos sin validar
		if !secureEnabled {
			c.Next()
			return
		}

		// 1. Extraer el token del encabezado Authorization
		tokenString := extractToken(c)
		if tokenString == "" {
			return
		}

		// 2. Parsear y validar el token
		token, err := jwt.ParseWithClaims(tokenString, &utils.Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return utils.JWTSecret, nil
		}, jwt.WithLeeway(5*time.Second))

		// 3. Manejar errores de validación
		if err != nil {
			handleJWTError(c, err)
			return
		}

		// 4. Verificar claims y almacenar en contexto
		if claims, ok := token.Claims.(*utils.Claims); ok && token.Valid {
			c.Set("jwtClaims", claims)
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Invalid token claims",
			})
		}
	}
}

func extractToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Authorization header is required",
		})
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Invalid Authorization header format (expected: Bearer <token>)",
		})
		return ""
	}

	return parts[1]
}

func handleJWTError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, jwt.ErrTokenMalformed):
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Malformed token",
		})
	case errors.Is(err, jwt.ErrTokenExpired):
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Token has expired",
		})
	case errors.Is(err, jwt.ErrTokenNotValidYet):
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Token not active yet",
		})
	default:
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Token validation failed: " + err.Error(),
		})
	}
}
