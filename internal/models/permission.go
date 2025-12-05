package models

type Permission struct {
    ID       uint   `gorm:"primaryKey"`
    Role     string `gorm:"index;size:32;not null"`
    Resource string `gorm:"index;size:64;not null"`
    Action   string `gorm:"index;size:32;not null"`
    Allowed  bool   `gorm:"not null;default:true"`
}
