package services

import (
	"time"

	"github.com/golang-jwt/jwt"
)

var JwtKey = []byte("my_secret_key")

type Credentials struct {
	Username string
	Password string
	Role     string
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
	Role string `json:"role"`
}

// GenerateToken создает JWT-токен для указанного пользователя с заданной ролью.
// @Summary Генерация JWT-токена
// @Description Создает JWT-токен с именем пользователя и ролью, срок действия токена составляет 5 минут.
// @Tags authentication
// @Param username query string true "Имя пользователя"
// @Param role query string true "Роль пользователя"
// @Success 200 {string} string "JWT token"
// @Failure 500 {string} string "Failed to generate token"
// @Router /generate-token [post]
func GenerateToken(username string, role string) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: username,
		Role:     role, // Включаем роль в токен
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtKey)
}
