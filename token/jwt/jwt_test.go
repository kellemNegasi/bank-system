package token

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kellemNegasi/bank-system/util"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	maker, err := NewJWT(util.RandString(32))
	require.NoError(t, err, fmt.Sprintf("Expected no error but %v was encountered", err))

	userName := util.RandomOwner()
	duration := time.Minute
	issuedAt := time.Now()
	expiresAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(userName, duration)
	require.NoError(t, err, fmt.Sprintf("Expected no error but %v was encountered", err))
	require.NotEmpty(t, token)

	pl, err := maker.Verify(token)
	require.NoError(t, err, fmt.Sprintf("Expected no error but %v was encountered", err))
	require.NotEmpty(t, pl)

	require.NotZero(t, pl.ID)
	require.Equal(t, userName, pl.Username)
	require.WithinDuration(t, issuedAt, pl.IssuedAt.Time, time.Second)
	require.WithinDuration(t, expiresAt, pl.ExpiresAt.Time, time.Second)

}

func TestExpiredJWTToken(t *testing.T) {
	maker, err := NewJWT(util.RandString(32))
	require.NoError(t, err, fmt.Sprintf("Expected no error but %v was encountered", err))

	token, err := maker.CreateToken(util.RandomOwner(), -time.Minute)
	require.NoError(t, err, fmt.Sprintf("Expected no error but %v was encountered", err))
	require.NotEmpty(t, token)

	pl, err := maker.Verify(token)
	require.Error(t, err)
	require.Nil(t, pl)

}

func TestInvalidJWTTokenAlgNone(t *testing.T) {
	pl, err := NewPayLoad(util.RandomOwner(), time.Minute)
	require.NoError(t, err, fmt.Sprintf("Expected no error but %v was encountered", err))

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, pl)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err, fmt.Sprintf("Expected no error but %v was encountered", err))

	maker, err := NewJWT(util.RandString(32))
	require.NoError(t, err)

	pl, err = maker.Verify(token)
	require.Error(t, err)
	require.Nil(t, pl)

}
