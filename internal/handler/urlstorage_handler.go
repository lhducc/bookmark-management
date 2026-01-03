package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lhducc/bookmark-management/internal/service"
	"net/http"
)

type urlShortenRequest struct {
	Url string `json:"url" biding:"required,url"`
	Exp int    `json:"exp" binding:"gte=0"`
}

type urlShortenResponse struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

type UrlShortenHandler interface {
	ShortenUrl(c *gin.Context)
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
// @Router /shorten-url [post]
func (h *urlShortenHandler) ShortenUrl(c *gin.Context) {
	var req urlShortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	code, err := h.urlService.ShortenUrl(c, req.Url, req.Exp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, urlShortenResponse{
		Message: "Shorten URL generated successfully!",
		Code:    code,
	})
}
