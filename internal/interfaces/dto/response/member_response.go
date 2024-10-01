package response

import (
	"github.com/HongJungWan/commerce-system/internal/domain"
	"github.com/HongJungWan/commerce-system/internal/helper"
	"time"
)

type MemberResponse struct {
	ID           uint   `json:"id"`
	MemberNumber string `json:"member_number"`
	Username     string `json:"username"`
	FullName     string `json:"full_name"`
	Email        string `json:"email"`
	CreatedAt    string `json:"created_at"`
	IsAdmin      bool   `json:"is_admin"`
	IsWithdrawn  bool   `json:"is_withdrawn"`
	WithdrawnAt  string `json:"withdrawn_at,omitempty"`
}

type RegisterMemberResponse struct {
	Message string         `json:"message"`
	User    MemberResponse `json:"user"`
}

func NewMemberResponse(member *domain.Member) *MemberResponse {
	return &MemberResponse{
		ID:           member.ID,
		MemberNumber: member.MemberNumber,
		Username:     member.Username,
		FullName:     member.FullName,
		Email:        member.Email,
		CreatedAt:    member.CreatedAt.Format(time.RFC3339),
		IsAdmin:      member.IsAdmin,
		IsWithdrawn:  member.IsWithdrawn,
		WithdrawnAt:  helper.FormatTime(member.WithdrawnAt),
	}
}
