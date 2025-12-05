package auth

import "golang.org/x/crypto/bcrypt"

func HashPassword(p string, cost int) (string, error) {
    b, err := bcrypt.GenerateFromPassword([]byte(p), cost)
    if err != nil { return "", err }
    return string(b), nil
}

func CheckPassword(hash, p string) bool {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(p)) == nil
}
