package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

func NewJWTManager(secretKey string, tokenDuration time.Duration) *JwtManager {
	return &JwtManager{SecretKey: secretKey, Duration: tokenDuration}
}

type JwtManager struct {
	SecretKey string
	Duration  time.Duration
}

type UserClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
	Role     string `json:"role"`
}

func (m *JwtManager) GenerateKey(user *User) (string, error) {
	claims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(m.Duration).Unix(),
		},
		Username: user.Username,
		Role:     user.Role,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.SecretKey))
}

func (m *JwtManager) Verify(token string) (*UserClaims, error) {
	jwtToken, err := jwt.ParseWithClaims(
		token,
		&UserClaims{},
		func(t *jwt.Token) (interface{}, error) {
			_, ok := t.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("Invalid signing method")
			}
			return []byte(m.SecretKey), nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := jwtToken.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("Invalid token claims")
	}

	return claims, nil
}
