package usecases

import (
	"errors"
	"time"

	"github.com/HongJungWan/commerce-system/internal/domain"
	"github.com/HongJungWan/commerce-system/internal/domain/repository"
	"github.com/HongJungWan/commerce-system/internal/interfaces/dto/request"
	"github.com/HongJungWan/commerce-system/internal/interfaces/dto/response"
)

type MemberInteractor struct {
	MemberRepository repository.MemberRepository
	AuthUseCase      *AuthUseCase
}

func NewMemberInteractor(repo repository.MemberRepository, auth *AuthUseCase) *MemberInteractor {
	return &MemberInteractor{
		MemberRepository: repo,
		AuthUseCase:      auth,
	}
}

func (mi *MemberInteractor) Register(req *request.RegisterMemberRequest) (*response.RegisterMemberResponse, error) {
	member := mi.toEntity(req)

	if err := member.Validate(); err != nil {
		return nil, err
	}

	existingMember, _ := mi.MemberRepository.GetByUserName(member.Username)
	if existingMember != nil {
		return nil, errors.New("이미 존재하는 사용자 ID입니다.")
	}

	if err := member.SetPassword(req.Password); err != nil {
		return nil, err
	}

	if err := mi.MemberRepository.Create(member); err != nil {
		return nil, err
	}

	token, _ := mi.AuthUseCase.GenerateToken(member)

	return &response.RegisterMemberResponse{
		Message: "회원 가입이 완료되었습니다.",
		Token:   token,
		User:    mi.toDTO(member),
	}, nil
}

func (mi *MemberInteractor) GetMyInfo(username string) (*response.MemberResponse, error) {
	member, err := mi.MemberRepository.GetByUserName(username)
	if err != nil {
		return nil, err
	}

	return mi.toDTO(member), nil
}

func (mi *MemberInteractor) UpdateMyInfo(username string, req *request.UpdateMemberRequest) error {
	member, err := mi.MemberRepository.GetByUserName(username)
	if err != nil {
		return err
	}

	if req.FullName != "" {
		member.FullName = req.FullName
	}
	if req.Email != "" {
		member.Email = req.Email
	}
	if req.Password != "" {
		if err := member.SetPassword(req.Password); err != nil {
			return err
		}
	}

	return mi.MemberRepository.Update(member)
}

func (mi *MemberInteractor) DeleteByUserName(username string) error {
	member, err := mi.MemberRepository.GetByUserName(username)
	if err != nil {
		return err
	}
	return mi.MemberRepository.Delete(member.ID)
}

func (mi *MemberInteractor) GetAllMembers() ([]*response.MemberResponse, error) {
	members, err := mi.MemberRepository.GetAll()
	if err != nil {
		return nil, err
	}

	var memberResponses []*response.MemberResponse
	for _, member := range members {
		memberResponses = append(memberResponses, mi.toDTO(member))
	}

	return memberResponses, nil
}

func (mi *MemberInteractor) GetMemberStats(month string) (map[string]interface{}, error) {
	joined, deleted, err := mi.MemberRepository.GetStatsByMonth(month)
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"month":           month,
		"joined_members":  joined,
		"deleted_members": deleted,
	}
	return stats, nil
}

func (mi *MemberInteractor) toEntity(req *request.RegisterMemberRequest) *domain.Member {
	return &domain.Member{
		MemberNumber: req.MemberNumber,
		Username:     req.Username,
		FullName:     req.FullName,
		Email:        req.Email,
		CreatedAt:    time.Now(),
	}
}

func (mi *MemberInteractor) toDTO(member *domain.Member) *response.MemberResponse {
	return &response.MemberResponse{
		ID:           member.ID,
		MemberNumber: member.MemberNumber,
		Username:     member.Username,
		FullName:     member.FullName,
		Email:        member.Email,
		CreatedAt:    member.CreatedAt.Format(time.RFC3339),
		IsAdmin:      member.IsAdmin,
		IsWithdrawn:  member.IsWithdrawn,
		WithdrawnAt:  formatTime(member.WithdrawnAt),
	}
}

func formatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}
