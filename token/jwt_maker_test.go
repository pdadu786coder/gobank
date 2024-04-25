package token

import (
	"fmt"
	"testing"
	"time"

	"github.com/pdadu/learn/db/util"
	"github.com/stretchr/testify/require"
)

func TestJWTma(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(40))
	require.NoError(t, err)
	username := util.RandomOwner()
	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := time.Now().Add(duration)
	token, err := maker.CreateToken(username, duration)
	fmt.Println(token)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	payload, err := maker.VerifyJWTToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.UserName)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredToken(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(40))
	require.NoError(t, err)
	token, err := maker.CreateToken(util.RandomOwner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	payload, err := maker.VerifyJWTToken(token)
	require.EqualError(t, err, ErrExiredToken.Error())
	require.Nil(t, payload)
}
