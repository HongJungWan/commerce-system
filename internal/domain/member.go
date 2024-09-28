package domain

import (
	"time"
)

type Member struct {
	MemberNumber string     `gorm:"primaryKey;column:member_number" json:"member_number"`
	ID           string     `gorm:"unique;not null" json:"id"`
	Password     string     `gorm:"not null" json:"password"`
	Name         string     `gorm:"not null" json:"name"`
	Email        string     `gorm:"unique;not null" json:"email"`
	JoinedAt     time.Time  `gorm:"not null" json:"joined_at"`
	LeftAt       *time.Time `json:"left_at"`
	Left         bool       `gorm:"not null;default:false" json:"left"`
	Orders       []Order    `gorm:"foreignKey:MemberNumber;references:MemberNumber" json:"orders"`
}
