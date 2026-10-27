package service

import (
	"errors"
	"fmt"
	"shopAPI/internal/authserver/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JwtProvider struct {
	secretKey    []byte
	accessExpiry time.Duration
	//refreshExpiry time.Duration
}

func NewJwtProvider(cfg *config.Config) *JwtProvider {
	return &JwtProvider{
		secretKey:    []byte(cfg.Secret),
		accessExpiry: 30 * time.Minute,
		//refreshExpiry: 24 * time.Hour,
	}
}

// Claims структура для хранения данных внутри токена
type Claims struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
	Type   string `json:"type"`
}

func (p *JwtProvider) GenerateAccessToken(userID uuid.UUID) (string, error) {
	return p.generateToken(userID, p.accessExpiry, "access")
}

/*
func (p *JwtProvider) GenerateRefreshToken(userID uuid.UUID) (string, error) {
	return p.generateToken(userID, p.refreshExpiry, "refresh")
}*/

func (p *JwtProvider) generateToken(userID uuid.UUID, expiry time.Duration, tokenType string) (string, error) {
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID: userID.String(),
		Type:   tokenType,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(p.secretKey)
}

func (p *JwtProvider) ValidateAccessToken(tokenStr string) (string, error) {
	claims, err := p.parseToken(tokenStr)
	if err != nil {
		return "", err
	}
	if claims.Type != "access" {
		return "", errors.New("expected access token")
	}
	return claims.UserID, nil
}

/*
func (p *JwtProvider) ValidateRefreshToken(tokenStr string) (bool, error) {
	claims, err := p.parseToken(tokenStr)
	if err != nil {
		return false, err
	}
	if claims.Type != "refresh" {
		return false, errors.New("expected refresh token")
	}
	return true, nil
}*/
/*
func (p *JwtProvider) ExtractUUID(tokenStr string) (uuid.UUID, error) {
	claims, err := p.parseToken(tokenStr)
	if err != nil {
		return uuid.Nil, err
	}
	if claims.Type != "access" {
		return uuid.Nil, errors.New("this is not an access token")
	}
	return uuid.Parse(claims.UserID)
}

func (p *JwtProvider) ExtractUUIDForRefresh(tokenStr string) (uuid.UUID, error) {
	claims, err := p.parseToken(tokenStr)
	if err != nil {
		return uuid.Nil, err
	}
	return uuid.Parse(claims.UserID)
}*/

func (p *JwtProvider) parseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return p.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
