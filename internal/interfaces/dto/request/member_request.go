package request

import (
	"github.com/HongJungWan/commerce-system/internal/domain"
	"github.com/google/uuid"
	"time"
)

const (
	MEMBER = "Member"
)

type RegisterMemberRequest struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	FullName    string `json:"full_name"`
	Email       string `json:"email"`
	IsAdmin     bool   `json:"is_admin"`
	IsWithdrawn bool   `json:"is_withdrawn"`
}

type UpdateMemberRequest struct {
	FullName string `json:"full_name,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

func (req *RegisterMemberRequest) ToEntity() (*domain.Member, error) {
	member := &domain.Member{
		MemberNumber: MEMBER + uuid.New().String(),
		Username:     req.Username,
		FullName:     req.FullName,
		Email:        req.Email,
		IsAdmin:      req.IsAdmin,
		IsWithdrawn:  false,
		CreatedAt:    time.Now(),
	}

	if err := member.AssignPassword(req.Password); err != nil {
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
		if err := member.AssignPassword(req.Password); err != nil {
			return nil, err
		}
	}

	return member, nil
}
