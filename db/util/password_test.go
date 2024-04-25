package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassw(t *testing.T) {
	password := RandomString(6)

	hashPassword, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashPassword)
	err = CheckPassword(password, hashPassword)
	require.NoError(t, err)
	wrongpassword := RandomString(6)
	err = CheckPassword(wrongpassword, hashPassword)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

}
