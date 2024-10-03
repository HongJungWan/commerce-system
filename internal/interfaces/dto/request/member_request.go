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
	AccountId   string `json:"account_id" example:"hong43ok"`
	Password    string `json:"password" example:"ghdwjddhks"`
	NickName    string `json:"nick_name" example:"hongmang"`
	Email       string `json:"email" example:"hong43ok@gmail.com"`
	IsAdmin     bool   `json:"is_admin" example:"true"`
	IsWithdrawn bool   `json:"is_withdrawn" example:"false"`
}

type UpdateMemberRequest struct {
	NickName string `json:"nick_name,omitempty" example:"hong"`
	Email    string `json:"email,omitempty" example:"hong43ok@naver.com"`
	Password string `json:"password,omitempty" example:"hong"`
}

func (req *CreateMemberRequest) CreateToEntity() (*domain.Member, error) {
	member := &domain.Member{
		MemberNumber: MEMBER + uuid.New().String(),
		AccountId:    req.AccountId,
		NickName:     req.NickName,
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
		NickName: req.NickName,
		Email:    req.Email,
	}

	if req.Password != "" {
		if err := member.AssignPassword(req.Password); err != nil {
			return nil, err
		}
	}

	return member, nil
}
