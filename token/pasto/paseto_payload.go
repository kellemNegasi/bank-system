package token

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type PasetoPayload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	ExpiresAt time.Time `json:"expires_at"`
	IssuedAt  time.Time `json:"issued_at"`
}

func NewPayLoad(user string, role string, duration time.Duration) (*PasetoPayload, error) {

	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	issuedAt := time.Now()
	expiresAt := issuedAt.Add(duration)

	return &PasetoPayload{
		ID:        id,
		Username:  user,
		Role:      role,
		ExpiresAt: expiresAt,
		IssuedAt:  issuedAt,
	}, nil
}

func (pl *PasetoPayload) LoadFromMap(m map[string]interface{}) error {
	data, err := json.Marshal(m)
	if err == nil {
		err = json.Unmarshal(data, pl)
	}
	return err
}

func (pl *PasetoPayload) ToJason() ([]byte, error) {
	data, err := json.Marshal(pl)
	if err != nil {
		return nil, err
	}

	return data, nil
}
