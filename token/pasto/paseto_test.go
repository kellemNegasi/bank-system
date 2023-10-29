package token

import (
	"fmt"
	"testing"
	"time"

	"github.com/kellemNegasi/bank-system/util"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {

	pMaker, err := NewPastoMaker(util.RandString(32))
	require.NoError(t, err, fmt.Sprintf("Expected no error but %v was encountered", err))

	userName := util.RandomOwner()
	duration := time.Minute
	issuedAt := time.Now()
	expiresAt := issuedAt.Add(duration)

	token, err := pMaker.CreateToken(userName, duration)
	require.NoError(t, err, fmt.Sprintf("Expected no error but %v was encountered", err))
	require.NotEmpty(t, token)

	pl, err := pMaker.VerifyToken(token)
	require.NoError(t, err, fmt.Sprintf("Expected no error but %v was encountered", err))
	require.NotEmpty(t, pl)

	require.NotZero(t, pl.ID)
	require.Equal(t, userName, pl.Username)
	require.WithinDuration(t, issuedAt, pl.IssuedAt, time.Second)
	require.WithinDuration(t, expiresAt, pl.ExpiresAt, time.Second)

}

func TestExpiredJWTToken(t *testing.T) {
	maker, err := NewPastoMaker(util.RandString(32))
	require.NoError(t, err, fmt.Sprintf("Expected no error but %v was encountered", err))

	token, err := maker.CreateToken(util.RandomOwner(), -time.Minute)
	require.NoError(t, err, fmt.Sprintf("Expected no error but %v was encountered", err))
	require.NotEmpty(t, token)

	pl, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.Nil(t, pl)

}
