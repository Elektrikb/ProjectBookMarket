package controllers

import (
	"Projectmugen/internal/services"
	"Projectmugen/internal/utils"
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetBooks обрабатывает запрос на получение списка книг с поддержкой фильтрации, сортировки и пагинации.
// @Summary Получение списка книг
// @Description Возвращает список книг с возможностью фильтрации по заголовку, сортировки и пагинации.
// @Tags books
// @Accept json
// @Produce json
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Количество книг на странице" default(10)
// @Param sort query string false "Поле для сортировки" default(id)
// @Param order query string false "Порядок сортировки (asc или desc)" default(asc)
// @Param title query string false "Фильтр по заголовку книги"
// @Success 200 {object} services.Book "total": int64, "page": int, "limit": int
// @Router /books [get]
func GetBooks(c *gin.Context) {
	var books []services.Book
	var total int64

	// Получаем параметры фильтров, сортировки и пагинации
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	sort := c.DefaultQuery("sort", "id")
	order := c.DefaultQuery("order", "asc")
	title := c.Query("title")

	// Преобразуем строковые параметры в int
	pageInt, _ := strconv.Atoi(page)
	limitInt, _ := strconv.Atoi(limit)
	offset := (pageInt - 1) * limitInt

	query := services.Db.Model(&services.Book{})

	// Применяем фильтры
	if title != "" {
		query = query.Where("title ILIKE ?", "%"+title+"%")
	}

	query.Count(&total)

	// Применяем сортировку
	if order != "asc" && order != "desc" {
		order = "asc" // По умолчанию ascending
	}
	query = query.Order(sort + " " + order).Limit(limitInt).Offset(offset)

	// Загружаем продукты и считаем общее количество
	query.Find(&books)

	// Возращаем результат
	c.JSON(http.StatusOK, gin.H{
		"data":  books,
		"total": total,
		"page":  pageInt,
		"limit": limitInt,
	})
}

// GetBookByID обрабатывает запрос на получение книги по ее идентификатору.
// @Summary Получение книги по ID
// @Description Возвращает книгу с указанным идентификатором.
// @Tags books
// @Accept json
// @Produce json
// @Param id path string true "Идентификатор книги"
// @Success 200 {object} services.Book
// @Failure 404 {object} models.ErrorResponse "Book not found"
// @Router /books/{id} [get]
func GetBookByID(c *gin.Context) {
	id := c.Param("id")
	var book services.Book
	if err := services.Db.First(&book, id).Error; err != nil {
		utils.HandleError(c, http.StatusNotFound, "Book not found")
		return
	}
	c.JSON(http.StatusOK, book)

}

// CreateBook обрабатывает запрос на создание новой книги.
// @Summary Создание новой книги
// @Description Создает новую книгу на основе переданных данных.
// @Tags books
// @Accept json
// @Produce json
// @Param book body services.Book true "Данные книги"
// @Success 201 {object} services.Book
// @Failure 400 {object} models.ErrorResponse "Invalid request"
// @Router /books [post]
func CreateBook(c *gin.Context) {
	var newBook services.Book

	if err := c.BindJSON(&newBook); err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Invalid request")
		return
	}

	services.Db.Create(&newBook)
	c.JSON(http.StatusCreated, newBook)

}

// UpdateBook обрабатывает запрос на обновление существующей книги по ее идентификатору.
// @Summary Обновление книги по ID
// @Description Обновляет данные книги с указанным идентификатором на основе переданных данных.
// @Tags books
// @Accept json
// @Produce json
// @Param id path string true "Идентификатор книги"
// @Param book body services.Book true "Обновленные данные книги"
// @Success 200 {object} services.Book
// @Failure 400 {object} models.ErrorResponse "Invalid request"
// @Failure 404 {object} models.ErrorResponse "Book not found"
// @Router /books/{id} [put]
func UpdateBook(c *gin.Context) {
	id := c.Param("id")
	var updatedBook services.Book

	if err := c.BindJSON(&updatedBook); err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Invalid request")
		return
	}

	if err := services.Db.Model(&services.Book{}).Where("id = ?", id).Updates(updatedBook).Error; err != nil {
		utils.HandleError(c, http.StatusNotFound, "Book not found")
		return
	}

	c.JSON(http.StatusOK, updatedBook)
}

// DeleteBook обрабатывает запрос на удаление книги по ее идентификатору.
// @Summary Удаление книги по ID
// @Description Удаляет книгу с указанным идентификатором.
// @Tags books
// @Accept json
// @Produce json
// @Param id path string true "Идентификатор книги"
// @Success 200 {object} models.ErrorResponse "Book deleted"
// @Failure 404 {object} models.ErrorResponse "Book not found"
// @Router /books/{id} [delete]
func DeleteBook(c *gin.Context) {
	id := c.Param("id")

	if err := services.Db.Delete(&services.Book{}, id).Error; err != nil {
		utils.HandleError(c, http.StatusNotFound, "Book not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book deleted"})

}

// GetBooksByYearRange обрабатывает запрос на получение книг в указанном диапазоне лет.
// @Summary Получение книг по диапазону лет
// @Description Возвращает список книг, выпущенных в заданном диапазоне лет.
// @Tags books
// @Accept json
// @Produce json
// @Param startYear query string true "Начальный год"
// @Param endYear query string true "Конечный год"
// @Success 200 {array} services.Book
// @Failure 500 {object} models.ErrorResponse "Error fetching books"
// @Router /books [get]
func GetBooksByYearRange(c *gin.Context) {
	startYear := c.Query("startYear")
	endYear := c.Query("endYear")

	var books []services.Book
	if err := services.Db.Where("year BETWEEN ? AND ?", startYear, endYear).Find(&books).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching books"})
		return
	}
	c.JSON(http.StatusOK, books)
}

// UpdateBooksPublisher обрабатывает запрос на обновление издателя для всех книг.
// @Summary Обновление издателя для всех книг
// @Description Обновляет издателя для всех книг в базе данных.
// @Tags books
// @Accept json
// @Produce json
// @Param publisher query string true "Имя издателя"
// @Success 200 {object} models.MessageResponse "Publisher updated successfully"
// @Failure 500 {object} models.ErrorResponse "Error updating publisher"
// @Router /books/publisher [put]
func UpdateBooksPublisher(c *gin.Context) {
	publisher := c.Query("publisher")
	tx := services.Db.Begin()

	if err := tx.Model(&services.Book{}).Update("publisher", publisher).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating publisher"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Publisher updated successfully"})
}

// CountBooksByAuthor обрабатывает запрос на подсчет количества книг по каждому автору.
// @Summary Подсчет книг по авторам
// @Description Возвращает количество книг для каждого автора в базе данных.
// @Tags books
// @Accept json
// @Produce json
// @Success 200 {array} models.Order "Количество книг по каждому автору"
// @Router /books/authors/count [get]
func CountBooksByAuthor(c *gin.Context) {
	var result []struct {
		Author string
		Count  int
	}

	services.Db.Model(&services.Book{}).Select("author, COUNT(*) as count").Group("author").Scan(&result)
	c.JSON(http.StatusOK, result)
}

// GetBooksWithTimeout обрабатывает запрос на получение книг с учетом тайм-аута.
// @Summary Получение книг с тайм-аутом
// @Description Возвращает список книг с фильтрацией, сортировкой и пагинацией, с тайм-аутом на выполнение запроса.
// @Tags books
// @Accept json
// @Produce json
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Количество книг на странице" default(10)
// @Param sort query string false "Поле для сортировки" default(id)
// @Param order query string false "Порядок сортировки" default(asc)
// @Param title query string false "Название книги"
// @Success 200 {object} services.Book "total": int64, "page": int, "limit": int}
// @Failure 500 {object} models.ErrorResponse "Failed to fetch books"}
// @Failure 408 {object} models.ErrorResponse "Request timed out"}
// @Router /books [get]
func GetBooksWithTimeout(c *gin.Context) {
	// Создаем контекст с тайм-аутом 2 секунды
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	var books []services.Book
	var total int64

	// Получаем параметры фильтров, сортировки и пагинации
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	sort := c.DefaultQuery("sort", "id")
	order := c.DefaultQuery("order", "asc")
	title := c.Query("title")

	// Преобразуем строковые параметры в int
	pageInt, _ := strconv.Atoi(page)
	limitInt, _ := strconv.Atoi(limit)
	offset := (pageInt - 1) * limitInt

	query := services.Db.Model(&services.Book{})

	// Применяем фильтры
	if title != "" {
		query = query.Where("name ILIKE ?", "%"+title+"%")
	}

	query.Count(&total)

	// Применяем сортировку
	if order != "asc" && order != "desc" {
		order = "asc" // По умолчанию ascending
	}
	query = query.Order(sort + " " + order).Limit(limitInt).Offset(offset)

	// Загружаем продукты с использованием контекста
	if err := query.WithContext(ctx).Find(&books).Error; err != nil {
		if err == context.DeadlineExceeded {
			utils.HandleError(c, http.StatusRequestTimeout, "Request timed out")
		} else {
			utils.HandleError(c, http.StatusInternalServerError, "Failed to fetch books")
		}
		return
	}

	// Возвращаем результат
	c.JSON(http.StatusOK, gin.H{
		"data":  books,
		"total": total,
		"page":  pageInt,
		"limit": limitInt,
	})
}
