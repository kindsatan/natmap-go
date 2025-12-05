package models

import "time"

type User struct {
    ID           uint `gorm:"primaryKey"`
    Username     string `gorm:"uniqueIndex;size:64;not null"`
    PasswordHash string `gorm:"size:255;not null"`
    Email        *string `gorm:"size:255"`
    Role         string `gorm:"size:32;not null;default:user"`
    IsActive     bool `gorm:"not null;default:true"`
    CreatedAt    time.Time `gorm:"autoCreateTime"`
    UpdatedAt    time.Time `gorm:"autoUpdateTime"`
    LastLoginAt  *time.Time
}
