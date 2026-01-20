package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lhducc/bookmark-management/internal/service"
	"github.com/rs/zerolog/log"
	"net/http"
)

type passwordHandler struct {
	svc service.Password
}

type Password interface {
	GenPass(c *gin.Context)
}

// NewPassword returns a new instance of the passwordHandler, which implements the Password interface.
// The returned passwordHandler is used to generate random passwords of length passLength, using characters from the character set 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'.
// The service parameter is used to generate the password and is required to implement the service.Password interface.
func NewPassword(svc service.Password) Password {

	return &passwordHandler{
		svc: svc,
	}
}

// GenPass generates a new password
// @Summary Generate a new password
// @Description Generate a new password
// @Tags Password
// @Success 200 {object} string "12345678"
// @Router /gen-pass [get]
func (h *passwordHandler) GenPass(c *gin.Context) {
	pass, err := h.svc.GeneratePassword()
	if err != nil {
		log.Error().Err(err).Msg("Service return error on GenPass")
		c.String(http.StatusInternalServerError, "err")
		return
	}

	c.String(http.StatusOK, pass)
}
