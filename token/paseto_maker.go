package token

import (
	"aidanwoods.dev/go-paseto"
	"encoding/json"
	"github.com/google/uuid"
	"strings"
	"time"
)

type PasetoMaker struct {
	symmetricKey paseto.V4SymmetricKey
	implicit     []byte
}

func NewPasetoMaker(implicit string) (Maker, error) {
	maker := &PasetoMaker{
		symmetricKey: paseto.NewV4SymmetricKey(),
		implicit:     []byte(implicit),
	}
	return maker, nil
}

func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", payload, err
	}

	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return "", payload, err
	}

	token, err := paseto.NewTokenFromClaimsJSON(payloadJson, nil)
	if err != nil {
		return "", payload, err
	}

	encryptedToken := token.V4Encrypt(maker.symmetricKey, maker.implicit)

	return encryptedToken, payload, err
}

func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	parser := paseto.NewParser()
	parser.AddRule(paseto.NotExpired())
	parsedToken, err := parser.ParseV4Local(maker.symmetricKey, token, maker.implicit)
	if err != nil {
		if strings.Contains(err.Error(), "expired") {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, err := getPayloadFromToken(parsedToken)
	return payload, nil
}

func getPayloadFromToken(t *paseto.Token) (*Payload, error) {
	id, err := t.GetString("id")
	if err != nil {
		return nil, ErrInvalidToken
	}
	username, err := t.GetString("username")
	if err != nil {
		return nil, ErrInvalidToken
	}
	issuedAt, err := t.GetIssuedAt()
	if err != nil {
		return nil, ErrInvalidToken
	}
	expiredAt, err := t.GetExpiration()
	if err != nil {
		return nil, ErrInvalidToken
	}

	return &Payload{
		ID:        uuid.MustParse(id),
		Username:  username,
		IssuedAt:  issuedAt,
		ExpiredAt: expiredAt,
	}, nil
}
