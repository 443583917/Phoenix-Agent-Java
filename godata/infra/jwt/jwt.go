package jwt

import (
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secret []byte
	expire time.Duration
}

type Claims struct {
	UserID   uint64 `json:"userId"`
	Username string `json:"username"`
	jwtlib.RegisteredClaims
}

func NewJWTManager(secret string, expire time.Duration) *JWTManager {
	return &JWTManager{
		secret: []byte(secret),
		expire: expire,
	}
}

func (m *JWTManager) GenerateToken(userID uint64, username string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwtlib.RegisteredClaims{
			IssuedAt:  jwtlib.NewNumericDate(now),
			ExpiresAt: jwtlib.NewNumericDate(now.Add(m.expire)),
		},
	}
	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

func (m *JWTManager) ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwtlib.ParseWithClaims(tokenStr, &Claims{},
		func(t *jwtlib.Token) (interface{}, error) {
			return m.secret, nil
		})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwtlib.ErrSignatureInvalid
}
