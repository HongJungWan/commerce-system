package usecases

import (
	"errors"

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

func (mi *MemberInteractor) Register(req *request.CreateMemberRequest) (*response.RegisterMemberResponse, error) {
	member, err := req.CreateToEntity()
	if err != nil {
		return nil, err
	}

	existingMember, _ := mi.MemberRepository.GetByUserName(member.Username)
	if existingMember != nil {
		return nil, errors.New("이미 존재하는 사용자 ID입니다.")
	}

	if err := mi.MemberRepository.Create(member); err != nil {
		return nil, err
	}

	memberResponse := response.NewMemberResponse(member)

	return &response.RegisterMemberResponse{
		Message: "회원 가입이 완료되었습니다.",
		User:    *memberResponse,
	}, nil
}

func (mi *MemberInteractor) GetMyInfo(username string) (*response.MemberResponse, error) {
	member, err := mi.MemberRepository.GetByUserName(username)
	if err != nil {
		return nil, err
	}

	memberResponse := response.NewMemberResponse(member)
	return memberResponse, nil
}

func (mi *MemberInteractor) UpdateMyInfo(username string, req *request.UpdateMemberRequest) error {
	member, err := mi.MemberRepository.GetByUserName(username)
	if err != nil {
		return err
	}

	if req.NickName != "" {
		member.FullName = req.NickName
	}
	if req.Email != "" {
		member.Email = req.Email
	}
	if req.Password != "" {
		if err := member.AssignPassword(req.Password); err != nil {
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
		memberResponses = append(memberResponses, response.NewMemberResponse(member))
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
