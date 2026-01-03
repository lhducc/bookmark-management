package service

import (
	"github.com/lhducc/bookmark-management/pkg/stringutils"
)

const (
	charset    = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	passLength = 10
)

type passwordService struct{}

// Password is an interface that defines the GeneratePassword method.
// GeneratePassword generates a random password of length passLength, using characters from the character set 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'. The generated password is returned as a string, or an error is returned if there was an issue generating the password.
// The character set used for generating the password is constant and does not change across different implementations of the interface. The length of the password is also constant and does not change across different implementations of the interface.
//
//go:generate mockery --name Password --filename pass_service.go
type Password interface {
	GeneratePassword() (string, error)
}

// NewPassword returns a new instance of the passwordService, which implements the Password interface.
// The returned passwordService is used to generate random passwords of length passLength, using characters from the character set 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'.
func NewPassword() Password {
	return &passwordService{}
}

// GeneratePassword generates a random password of length passLength, using characters from the character set 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'. The generated password is returned as a string, or an error is returned if there was an issue generating the password.
// The character set used for generating the password is constant and cannot be changed externally.
// The length of the generated password is constant and cannot be changed externally.
// If an error occurs while generating the password, the error is returned immediately and the generated password is an empty string.
func (s *passwordService) GeneratePassword() (string, error) {
	return stringutils.GenerateCode(passLength)
}
