package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lhducc/bookmark-management/internal/service"
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

// GenPass handles a GET request to generate a random password.
// It returns a random password of length 10, using characters from the character set 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'.
// If an error occurs while generating the password, it returns an HTTP 500 error with the message 'err'.
// Otherwise, it returns an HTTP 200 OK response with the generated password as a string.
func (h *passwordHandler) GenPass(c *gin.Context) {
	pass, err := h.svc.GeneratePassword()
	if err != nil {
		c.String(http.StatusInternalServerError, "err")
		return
	}

	c.String(http.StatusOK, pass)
}
