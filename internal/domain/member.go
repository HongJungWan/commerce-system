package domain

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Member struct {
	ID           uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	MemberNumber string     `gorm:"unique;not null" json:"member_number"`
	UserID       string     `gorm:"unique;not null" json:"user_id"`
	Password     string     `gorm:"not null" json:"password"`
	Name         string     `gorm:"not null" json:"name"`
	Email        string     `gorm:"unique;not null" json:"email"`
	IsAdmin      bool       `gorm:"default:false" json:"is_admin"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt    *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (m *Member) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	m.Password = string(hashedPassword)
	return nil
}

func (m *Member) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(m.Password), []byte(password))
	return err == nil
}

func (m *Member) Validate() error {
	if m.UserID == "" || m.Password == "" || m.Name == "" || m.Email == "" {
		return errors.New("필수 필드가 누락되었습니다.")
	}
	// FIXME: 추가적인 검증 로직을 여기에 추가
	return nil
}
