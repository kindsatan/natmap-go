package auth

import (
    "testing"
    "time"
)

func TestGenerateAndParseToken(t *testing.T) {
    tok, _, err := GenerateToken(1, "user", "admin", "secret", time.Minute)
    if err != nil { t.Fatal(err) }
    c, err := ParseToken(tok, "secret")
    if err != nil { t.Fatal(err) }
    if c.UserID != 1 || c.Username != "user" || c.Role != "admin" { t.Fatal("claims mismatch") }
}
