package auth

import (
    "testing"
    "time"
)

func TestGenerateRefreshToken(t *testing.T) {
    raw, hash, err := GenerateRefreshToken()
    if err != nil { t.Fatal(err) }
    if HashRefreshToken(raw) != hash { t.Fatal("hash mismatch") }
}

func TestJWTWithRole(t *testing.T) {
    tok, _, err := GenerateToken(1, "u", "user", "s", time.Minute)
    if err != nil { t.Fatal(err) }
    c, err := ParseToken(tok, "s")
    if err != nil { t.Fatal(err) }
    if c.Role != "user" { t.Fatal("role mismatch") }
}
