package controllers

import (
	"Projectmugen/internal/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// Login обрабатывает входящие запросы на аутентификацию пользователя.
// @Summary Аутентификация пользователя
// @Description Проверяет учетные данные пользователя и возвращает токен при успешной аутентификации.
// @Accept json
// @Produce json
// @Param creds body services.Credentials true "Учетные данные пользователя"
// @Success 200 {object} models.TokenResponse "токен доступа"
// @Failure 400 {object} models.ErrorResponse "invalid request"
// @Failure 401 {object} models.ErrorResponse "unauthorized"
// @Failure 500 {object} models.ErrorResponse "could not create token"
// @Router /login [post]
func Login(c *gin.Context) {
	var creds services.Credentials
	if err := c.BindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	storedPassword, ok := Users[creds.Username]
	if !ok || storedPassword != creds.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	role, roleExists := Roles[creds.Username]
	if !roleExists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "role not assigned"})
		return
	}

	token, err := services.GenerateToken(creds.Username, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// AuthMiddleware обеспечивает защиту маршрутов, проверяя наличие и валидность токена авторизации.
// @Summary Проверка токена авторизации
// @Description Middleware для проверки JWT токена в заголовке Authorization.
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} models.MessageResponse "success"
// @Failure 401 {object} models.ErrorResponse "unauthorized"
// @Router /protected-route [get]
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		claims := &services.Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return services.JwtKey, nil
		})

		if err != nil || !token.Valid {
			if err == jwt.ErrSignatureInvalid {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
				c.Abort()
				return
			}

			if ve, ok := err.(*jwt.ValidationError); ok && ve.Errors == jwt.ValidationErrorExpired {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "token expired"})
				c.Abort()
				return
			}

			c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			c.Abort()
			return
		}

		c.Next()
	}
}

var Users = map[string]string{
	"admin":    "admin123",
	"user":     "password",
	"elektrik": "2003",
}

// Register обрабатывает регистрацию нового пользователя.
// @Summary Регистрация нового пользователя
// @Description Обрабатывает запрос на регистрацию, проверяет данные и сохраняет пользователя.
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body services.Credentials true "Данные для регистрации пользователя"
// @Success 201 {object} models.MessageResponse "user registered successfully"
// @Failure 400 {object} models.ErrorResponse "invalid request"
// @Failure 409 {object} models.ErrorResponse "user already exists"
// @Router /register [post]
func Register(c *gin.Context) {
	var creds services.Credentials
	if err := c.BindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	if _, exists := Users[creds.Username]; exists {
		c.JSON(http.StatusConflict, gin.H{"message": "user already exists"})
		return
	}

	role := "user"

	if creds.Role != "" {
		role = creds.Role
	}

	Users[creds.Username] = creds.Password
	Roles[creds.Username] = role

	c.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
}

var Roles = map[string]string{
	"admin":    "admin",
	"user":     "user",
	"elektrik": "user",
}

// RoleMiddleware проверяет, имеет ли пользователь необходимую роль для доступа к маршруту.
// @Summary Проверка роли пользователя
// @Description Middleware для проверки роли пользователя на основе JWT токена.
// @Tags auth
// @Accept json
// @Produce json
// @Param requiredRole path string true "Необходимая роль для доступа"
// @Success 200 {object} models.MessageResponse "success"
// @Failure 401 {object} models.ErrorResponse "unauthorized"
// @Failure 403 {object} models.ErrorResponse "forbidden"
// @Router /protected-route/{requiredRole} [get]
func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		claims := &services.Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return services.JwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			c.Abort()
			return
		}

		if claims.Role != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"message": "forbidden"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Refresh обрабатывает запрос на обновление токена авторизации.
// @Summary Обновление токена авторизации
// @Description Middleware для обновления JWT токена, если он истек.
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} models.TokenResponse "новый токен"
// @Failure 400 {object} models.ErrorResponse "token not expired enough"
// @Failure 401 {object} models.ErrorResponse "unauthorized"
// @Failure 500 {object} models.ErrorResponse "could not create token"
// @Router /refresh [post]
func Refresh(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	claims := &services.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return services.JwtKey, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		c.JSON(http.StatusBadRequest, gin.H{"message": "token not expired enough"})
		return
	}

	newToken, err := services.GenerateToken(claims.Username, claims.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": newToken})
}
