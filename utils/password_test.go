package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPass(t *testing.T) {
	password := RandomString(6)

	hash, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hash)

	err = CheckPasswordHash(password, hash)
	require.NoError(t, err)

	wrongPassword := RandomString(6)
	err = CheckPasswordHash(wrongPassword, hash)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

}
