package jwt

import (
	"errors"
	"strings"
	"time"

	jwtv5 "github.com/golang-jwt/jwt/v5"
)

const (
	TokenAccess  = "access"
	TokenRefresh = "refresh"
	TokenPreview = "preview"
)

type Claims struct {
	UserID       int64    `json:"user_id"`
	UserCode     string   `json:"user_code"`
	Username     string   `json:"username"`
	OperatorID   int64    `json:"operator_id"`
	OperatorCode string   `json:"operator_code"`
	RoleCodes    []string `json:"role_codes"`
	Salt         string   `json:"salt"`
	TokenType    string   `json:"token_type"`
	IsPlatform   bool     `json:"is_platform"`
	jwtv5.RegisteredClaims
}

func (c Claims) WithType(tokenType string) Claims {
	c.TokenType = tokenType
	return c
}

func Sign(secret string, expireSec int64, c Claims) (string, int64, error) {
	if secret == "" {
		return "", 0, errors.New("jwt secret is empty")
	}
	if expireSec <= 0 {
		expireSec = 86400
	}
	if c.TokenType == "" {
		c.TokenType = TokenAccess
	}
	now := time.Now()
	exp := now.Add(time.Duration(expireSec) * time.Second)
	c.RegisteredClaims = jwtv5.RegisteredClaims{
		IssuedAt:  jwtv5.NewNumericDate(now),
		ExpiresAt: jwtv5.NewNumericDate(exp),
	}
	t := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, c)
	s, err := t.SignedString([]byte(secret))
	return s, exp.Unix(), err
}

func Parse(secret, token string) (*Claims, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return nil, errors.New("missing token")
	}
	parsed, err := jwtv5.ParseWithClaims(token, &Claims{}, func(t *jwtv5.Token) (any, error) {
		if _, ok := t.Method.(*jwtv5.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	c, ok := parsed.Claims.(*Claims)
	if !ok || !parsed.Valid {
		return nil, errors.New("invalid token")
	}
	return c, nil
}

func ParseTyped(secret, token, wantType string) (*Claims, error) {
	c, err := Parse(secret, token)
	if err != nil {
		return nil, err
	}
	if c.TokenType != wantType {
		return nil, errors.New("invalid token type")
	}
	return c, nil
}

func StripBearer(v string) string {
	v = strings.TrimSpace(v)
	if len(v) > 7 && strings.EqualFold(v[:7], "Bearer ") {
		return strings.TrimSpace(v[7:])
	}
	return v
}

func IsExpired(err error) bool {
	return errors.Is(err, jwtv5.ErrTokenExpired)
}
