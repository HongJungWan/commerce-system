package request

import (
	"github.com/HongJungWan/commerce-system/internal/domain"
	"github.com/google/uuid"
	"time"
)

const (
	MEMBER = "Member"
)

type CreateMemberRequest struct {
	AccountId   string `json:"account_id"`
	Password    string `json:"password"`
	NickName    string `json:"nick_name"`
	Email       string `json:"email"`
	IsAdmin     bool   `json:"is_admin"`
	IsWithdrawn bool   `json:"is_withdrawn"`
}

type UpdateMemberRequest struct {
	NickName string `json:"nick_name,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

func (req *CreateMemberRequest) CreateToEntity() (*domain.Member, error) {
	member := &domain.Member{
		MemberNumber: MEMBER + uuid.New().String(),
		Username:     req.AccountId,
		FullName:     req.NickName,
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

func (req *UpdateMemberRequest) UpdateToEntity() (*domain.Member, error) {
	member := &domain.Member{
		FullName: req.NickName,
		Email:    req.Email,
	}

	if req.Password != "" {
		if err := member.AssignPassword(req.Password); err != nil {
			return nil, err
		}
	}

	return member, nil
}
