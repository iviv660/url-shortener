package handler

import (
	"app/internal/usecase"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	uc *usecase.UseCase
}

func NewHandler(uc *usecase.UseCase) (*Handler, *gin.Engine) {
	h := &Handler{uc: uc}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.POST("/shorten", h.handleShorten)
	r.GET("/:short", h.handleRedirect)

	// Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return h, r
}

// handleShorten godoc
// @Summary      Создать короткую ссылку
// @Description  Принимает URL и возвращает короткий код/ссылку
// @Accept       json
// @Produce      json
// @Param        request body ShortenRequest true "URL для сокращения"
// @Success      201 {object} ShortenResponse
// @Failure      400 {object} ErrorResponse "Неверный запрос"
// @Failure      429 {object} ErrorResponse "Превышен лимит"
// @Router       /shorten [post]
func (h *Handler) handleShorten(c *gin.Context) {
	var req ShortenRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.URL == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	short, err := h.uc.CreateShortURL(c.Request.Context(), req.URL)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	base := os.Getenv("BASE_URL")
	if base == "" {
		base = "http://localhost:3000"
	}
	shortURL := base + "/" + short

	// 201 + Location
	c.Header("Location", shortURL)
	c.JSON(http.StatusCreated, ShortenResponse{
		ShortCode: short,
		ShortURL:  shortURL,
	})
}

// handleRedirect godoc
// @Summary      Редирект по короткому коду
// @Description  Делает 302 Redirect на оригинальный URL
// @Param        short path string true "короткий код"
// @Success      302
// @Failure      404 {object} ErrorResponse "Не найдено"
// @Router       /{short} [get]
func (h *Handler) handleRedirect(c *gin.Context) {
	short := c.Param("short")
	if short == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "empty short code"})
		return
	}

	long, err := h.uc.OpenUrl(c.Request.Context(), short)
	if err != nil || long == "" {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "not found"})
		return
	}

	c.Redirect(http.StatusFound, long)
}
