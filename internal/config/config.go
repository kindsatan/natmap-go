package config

import (
    "os"
    "strconv"
    "time"

    "github.com/joho/godotenv"
)

type Config struct {
    HTTPAddr    string
    SQLitePath  string
    JWTSecret   string
    TokenTTL    time.Duration
    RefreshTTL  time.Duration
    BcryptCost  int
    SeedUser    bool
    SeedUsername string
    SeedPassword string
}

func Load() Config {
    _ = godotenv.Load()
    addr := getenv("HTTP_ADDR", ":8080")
    dbpath := getenv("SQLITE_PATH", "./data/app.db")
    secret := getenv("JWT_SECRET", "change-me")
    ttlStr := getenv("TOKEN_TTL", "24h")
    ttl, err := time.ParseDuration(ttlStr)
    if err != nil { ttl = 24 * time.Hour }
    rttlStr := getenv("REFRESH_TTL", "168h")
    rttl, err := time.ParseDuration(rttlStr)
    if err != nil { rttl = 7 * 24 * time.Hour }
    costStr := getenv("BCRYPT_COST", "12")
    cost, err := strconv.Atoi(costStr)
    if err != nil || cost < 10 { cost = 12 }
    seed := getenv("SEED_USER", "0") == "1"
    return Config{
        HTTPAddr: addr,
        SQLitePath: dbpath,
        JWTSecret: secret,
        TokenTTL: ttl,
        RefreshTTL: rttl,
        BcryptCost: cost,
        SeedUser: seed,
        SeedUsername: getenv("SEED_USERNAME", ""),
        SeedPassword: getenv("SEED_PASSWORD", ""),
    }
}

func getenv(k, def string) string {
    v := os.Getenv(k)
    if v == "" { return def }
    return v
}
