package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lhducc/bookmark-management/internal/service"
	"github.com/rs/zerolog/log"
	"net/http"
)

type urlShortenRequest struct {
	Url string `json:"url" binding:"required,url"`
	Exp int    `json:"exp" binding:"required,gte=604800"`
}

type urlShortenResponse struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

type UrlShortenHandler interface {
	ShortenUrl(c *gin.Context)
	GetUrl(c *gin.Context)
}

type urlShortenHandler struct {
	urlService service.ShortenUrl
}

func NewUrlShortenHandler(svc service.ShortenUrl) UrlShortenHandler {
	return &urlShortenHandler{urlService: svc}
}

// ShortenUrl shortens a given URL and returns a shortened URL code.
// @Summary Shorten URL
// @Description Shortens a given URL and returns a shortened URL code.
// @Tags URL Shortener
// @Accept json
// @Produce json
// @Param urlShortenRequest body urlShortenRequest true "URL to shorten"
// @Success 200 {object} urlShortenResponse
// @Failure 400 {object} map[string]string "Bad Request - invalid URL or validation error"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /v1/links/shorten [post]
func (h *urlShortenHandler) ShortenUrl(c *gin.Context) {
	var req urlShortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	code, err := h.urlService.ShortenUrl(c, req.Url, req.Exp)
	if err != nil {
		log.Error().Str("url", req.Url).Err(err).Msg("Service return error on ShortenUrl")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, urlShortenResponse{
		Message: "Shorten URL generated successfully!",
		Code:    code,
	})
}

// GetUrl shortens a given URL and returns a shortened URL code.
// @Summary Get URL
// @Description Get URL by code
// @Tags URL Shortener
// @Accept json
// @Produce json
// @Param code path string true "Url code" Format(string)
// @Success 302
// @Failure 400  "Bad Request - invalid URL or validation error"
// @Failure 404  "URL not found"
// @Failure 500  "Internal Server Error"
// @Router /v1/links/redirect/{code} [get]
func (h *urlShortenHandler) GetUrl(c *gin.Context) {
	code := c.Param("code")

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "wrong format"})
		return
	}

	url, err := h.urlService.GetUrl(c, code)
	if err != nil {
		if errors.Is(err, service.ErrCodeNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "url not found"})
			return
		}

		log.Error().Err(err).Msg("Service return error on GetUrl")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
	}

	c.Redirect(http.StatusFound, url)
}
