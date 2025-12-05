package auth

import "testing"

func TestHashAndCheckPassword(t *testing.T) {
    h, err := HashPassword("secret", 12)
    if err != nil { t.Fatal(err) }
    if !CheckPassword(h, "secret") { t.Fatal("check failed") }
    if CheckPassword(h, "wrong") { t.Fatal("should fail") }
}
