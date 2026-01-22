package stringutils

import (
	"bytes"
	"crypto/rand"
	"math/big"
)

const (
	charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
)

type KeyGen interface {
	GenerateCode(length int) (string, error)
}

type keyGen struct {
}

//go:generate mockery --name KeyGen --filename keygen.go
func NewKeyGen() KeyGen {
	return &keyGen{}
}

// GenerateCode generates a random string of the given length, using characters from the character set 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'.
// The generated code is returned as a string, or an error is returned if there was an issue generating the code.
// The character set used for generating the code is constant and does not change across different implementations of the interface.
// The length of the generated code is constant and does not change across different implementations of the interface.
// If an error occurs while generating the code, the error is returned immediately and the generated code is an empty string.
func (k *keyGen) GenerateCode(length int) (string, error) {
	return GenerateCode(length)
}

// GenerateCode generates a random string of the given length, using characters from the character set 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'.
// The generated code is returned as a string, or an error is returned if there was an issue generating the code.
// The character set used for generating the code is constant and does not change across different implementations of the interface.
// The length of the generated code is constant and does not change across different implementations of the interface.
// If an error occurs while generating the code, the error is returned immediately and the generated code is an empty string.
func GenerateCode(length int) (string, error) {
	var strBuilder bytes.Buffer

	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}

		strBuilder.WriteByte(charset[randomIndex.Int64()])
	}

	return strBuilder.String(), nil
}
