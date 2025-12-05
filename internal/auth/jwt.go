package auth

import (
    "strconv"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    UserID   uint   `json:"user_id"`
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}

func GenerateToken(userID uint, username, role, secret string, ttl time.Duration) (string, time.Time, error) {
    exp := time.Now().Add(ttl)
    c := Claims{
        UserID: userID,
        Username: username,
        Role: role,
        RegisteredClaims: jwt.RegisteredClaims{Subject: strconv.FormatUint(uint64(userID), 10), ExpiresAt: jwt.NewNumericDate(exp), IssuedAt: jwt.NewNumericDate(time.Now())},
    }
    t := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
    s, err := t.SignedString([]byte(secret))
    if err != nil { return "", time.Time{}, err }
    return s, exp, nil
}

func ParseToken(token, secret string) (*Claims, error) {
    p, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (interface{}, error) { return []byte(secret), nil })
    if err != nil { return nil, err }
    if c, ok := p.Claims.(*Claims); ok && p.Valid { return c, nil }
    return nil, jwt.ErrTokenInvalidClaims
}
