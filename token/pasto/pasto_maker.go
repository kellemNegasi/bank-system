package token

import (
	"time"

	"aidanwoods.dev/go-paseto"
)

type PasetoMaker struct {
	secretKey paseto.V4SymmetricKey
}

// NewPastorMaker returns an new `PastoMaker`.
func NewPastoMaker(key string) (*PasetoMaker, error) {
	sk, err := paseto.V4SymmetricKeyFromBytes([]byte(key))
	if err != nil {
		return nil, err
	}
	return &PasetoMaker{
		secretKey: sk,
	}, nil
}

// CreateToken creates a symmetrically signed token string from the provided payload values.
func (pm *PasetoMaker) CreateToken(user string, duration time.Duration) (string, error) {
	pl, err := NewPayLoad(user, "", duration)

	if err != nil {
		return "", err
	}

	jsonData, err := pl.ToJason()

	if err != nil {
		return "", err
	}

	token, err := paseto.NewTokenFromClaimsJSON(jsonData, nil)

	if err != nil {
		return "", err
	}

	token.SetExpiration(pl.ExpiresAt)
	token.SetIssuedAt(pl.IssuedAt)

	signedString := token.V4Encrypt(pm.secretKey, nil)
	return signedString, nil
}

func (pm *PasetoMaker) VerifyToken(token string) (*PasetoPayload, error) {
	parser := paseto.NewParser()

	tk, err := parser.ParseV4Local(pm.secretKey, token, nil)
	if err != nil {
		return nil, err
	}

	pl := PasetoPayload{}

	err = pl.LoadFromMap(tk.Claims())
	if err != nil {
		return nil, err
	}

	return &pl, nil
}
