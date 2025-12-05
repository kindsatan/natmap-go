package db

import (
    "errors"
    "os"
    "path/filepath"

    "github.com/glebarez/sqlite"
    "gorm.io/gorm"
)

func SetupDB(path string) (*gorm.DB, error) {
    dir := filepath.Dir(path)
    if dir != "." && dir != "" {
        if err := os.MkdirAll(dir, 0755); err != nil { return nil, err }
    }
    d, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
    if err != nil { return nil, err }
    sqlDB, err := d.DB()
    if err != nil { return nil, err }
    if sqlDB == nil { return nil, errors.New("db nil") }
    return d, nil
}
