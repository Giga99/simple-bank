package token

import (
	"aidanwoods.dev/go-paseto"
	"encoding/json"
	"github.com/google/uuid"
	"strings"
	"time"
)

type PasetoMaker struct {
	symetricKey paseto.V4SymmetricKey
	implicit    []byte
}

func NewPasetoMaker(implicit string) (Maker, error) {
	maker := &PasetoMaker{
		symetricKey: paseto.NewV4SymmetricKey(),
		implicit:    []byte(implicit),
	}
	return maker, nil
}

func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	println(payloadJson)

	token, err := paseto.NewTokenFromClaimsJSON(payloadJson, nil)
	if err != nil {
		return "", err
	}
	println(token)

	return token.V4Encrypt(maker.symetricKey, maker.implicit), nil
}

func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	parser := paseto.NewParser()
	parser.AddRule(paseto.NotExpired())
	parsedToken, err := parser.ParseV4Local(maker.symetricKey, token, maker.implicit)
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
