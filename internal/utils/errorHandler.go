package utils

import "github.com/gin-gonic/gin"

// HandleError обрабатывает ошибки и возвращает JSON-ответ с кодом состояния и сообщением об ошибке.
// @Summary Обработка ошибок
// @Description Возвращает JSON-ответ с указанным кодом состояния и сообщением об ошибке.
// @Tags error handling
// @Param statusCode query int true "Код состояния HTTP"
// @Param message query string true "Сообщение об ошибке"
// @Success 200 {object} models.ErrorResponse "error message"
// @Router /handle-error [post]
func HandleError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{"error": message})
}
