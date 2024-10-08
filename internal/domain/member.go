package domain

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Member struct {
	ID           uint       `gorm:"primaryKey;autoIncrement" json:"id"`   // 기본 키
	MemberNumber string     `gorm:"unique;not null" json:"member_number"` // 회원번호
	AccountId    string     `gorm:"unique;not null" json:"account_id"`    // 아이디
	Password     string     `gorm:"not null" json:"password"`             // 패스워드 (해싱 처리)
	NickName     string     `gorm:"not null" json:"nick_name"`            // 회원명
	Email        string     `gorm:"unique;not null" json:"email"`         // 이메일
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at"`     // 가입일
	IsAdmin      bool       `gorm:"default:false" json:"is_admin"`        // 관리자 여부
	IsWithdrawn  bool       `gorm:"default:false" json:"is_withdrawn"`    // 탈퇴 여부
	WithdrawnAt  *time.Time `gorm:"index" json:"withdrawn_at,omitempty"`  // 탈퇴일
}

func (m *Member) Validate() error {
	if m.MemberNumber == "" {
		return errors.New("회원번호가 누락되었습니다.")
	}
	if m.AccountId == "" {
		return errors.New("아이디가 누락되었습니다.")
	}
	if m.Password == "" {
		return errors.New("패스워드가 누락되었습니다.")
	}
	if m.NickName == "" {
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
	m.Password = string(hashedPassword)
	return nil
}

func (m *Member) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(m.Password), []byte(password))
	return err == nil
}
