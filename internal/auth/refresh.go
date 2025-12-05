package auth

import (
    "crypto/rand"
    "crypto/sha256"
    "encoding/hex"
)

func GenerateRefreshToken() (string, string, error) {
    b := make([]byte, 32)
    if _, err := rand.Read(b); err != nil { return "", "", err }
    raw := hex.EncodeToString(b)
    h := sha256.Sum256([]byte(raw))
    return raw, hex.EncodeToString(h[:]), nil
}

func HashRefreshToken(raw string) string {
    h := sha256.Sum256([]byte(raw))
    return hex.EncodeToString(h[:])
}
