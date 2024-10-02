package domain

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Member struct {
	ID             uint       `gorm:"primaryKey;autoIncrement" json:"id"`   // 기본 키
	MemberNumber   string     `gorm:"unique;not null" json:"member_number"` // 회원번호
	Username       string     `gorm:"unique;not null" json:"username"`      // 아이디
	HashedPassword string     `gorm:"not null" json:"hashed_password"`      // 패스워드 (해싱 처리)
	FullName       string     `gorm:"not null" json:"full_name"`            // 회원명
	Email          string     `gorm:"unique;not null" json:"email"`         // 이메일
	CreatedAt      time.Time  `gorm:"autoCreateTime" json:"created_at"`     // 가입일
	IsAdmin        bool       `gorm:"default:false" json:"is_admin"`        // 관리자 여부
	IsWithdrawn    bool       `gorm:"default:false" json:"is_withdrawn"`    // 탈퇴 여부
	WithdrawnAt    *time.Time `gorm:"index" json:"withdrawn_at,omitempty"`  // 탈퇴일
}

func NewMember(memberNumber, username, fullName, email string, isAdmin bool) (*Member, error) {
	member := &Member{
		MemberNumber: memberNumber,
		Username:     username,
		FullName:     fullName,
		Email:        email,
		IsAdmin:      isAdmin,
		IsWithdrawn:  false,
		CreatedAt:    time.Now(),
	}
	return member, nil
}

func (m *Member) Validate() error {
	if m.MemberNumber == "" {
		return errors.New("회원번호가 누락되었습니다.")
	}
	if m.Username == "" {
		return errors.New("아이디가 누락되었습니다.")
	}
	if m.HashedPassword == "" {
		return errors.New("패스워드가 누락되었습니다.")
	}
	if m.FullName == "" {
		return errors.New("회원명이 누락되었습니다.")
	}
	if m.Email == "" {
		return errors.New("이메일이 누락되었습니다.")
	}
	return nil
}

func (m *Member) AssignPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	m.HashedPassword = string(hashedPassword)
	return nil
}

func (m *Member) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(m.HashedPassword), []byte(password))
	return err == nil
}
