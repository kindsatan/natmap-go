package models

import "time"

type RefreshToken struct {
    ID         uint `gorm:"primaryKey"`
    UserID     uint `gorm:"index;not null"`
    TokenHash  string `gorm:"size:128;not null;uniqueIndex"`
    ExpiresAt  time.Time `gorm:"not null"`
    Revoked    bool `gorm:"not null;default:false"`
    CreatedAt  time.Time
    UpdatedAt  time.Time
}
