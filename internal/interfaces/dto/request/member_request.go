package request

import (
	"github.com/HongJungWan/commerce-system/internal/domain"
	"time"
)

type RegisterMemberRequest struct {
	MemberNumber string `json:"member_number"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	FullName     string `json:"full_name"`
	Email        string `json:"email"`
}

type UpdateMemberRequest struct {
	FullName string `json:"full_name,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

func (req *RegisterMemberRequest) ToEntity() (*domain.Member, error) {
	member := &domain.Member{
		MemberNumber: req.MemberNumber,
		Username:     req.Username,
		FullName:     req.FullName,
		Email:        req.Email,
		CreatedAt:    time.Now(),
	}

	if err := member.SetPassword(req.Password); err != nil {
		return nil, err
	}

	if err := member.Validate(); err != nil {
		return nil, err
	}

	return member, nil
}

func (req *UpdateMemberRequest) ToEntity() (*domain.Member, error) {
	member := &domain.Member{
		FullName: req.FullName,
		Email:    req.Email,
	}

	if req.Password != "" {
		if err := member.SetPassword(req.Password); err != nil {
			return nil, err
		}
	}

	// 업데이트의 경우 필드가 선택적이므로 유효성 검사를 생략하거나 필요한 경우 추가
	return member, nil
}
