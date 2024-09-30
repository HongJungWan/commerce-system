package usecases

import (
	"errors"

	"github.com/HongJungWan/commerce-system/internal/domain"
	"github.com/HongJungWan/commerce-system/internal/domain/repository"
)

type MemberInteractor struct {
	MemberRepository repository.MemberRepository
}

func NewMemberInteractor(repo repository.MemberRepository) *MemberInteractor {
	return &MemberInteractor{MemberRepository: repo}
}

func (mi *MemberInteractor) Register(member *domain.Member) error {
	if err := member.Validate(); err != nil {
		return err
	}

	existingMember, _ := mi.MemberRepository.GetByUserID(member.Username)
	if existingMember != nil {
		return errors.New("이미 존재하는 사용자 ID입니다.")
	}

	if err := member.SetPassword(member.HashedPassword); err != nil {
		return err
	}

	return mi.MemberRepository.Create(member)
}

func (mi *MemberInteractor) GetByUserID(userID string) (*domain.Member, error) {
	return mi.MemberRepository.GetByUserID(userID)
}

func (mi *MemberInteractor) UpdateByUserID(userID string, updateData *domain.Member) error {
	member, err := mi.MemberRepository.GetByUserID(userID)
	if err != nil {
		return err
	}

	if updateData.FullName != "" {
		member.FullName = updateData.FullName
	}
	if updateData.Email != "" {
		member.Email = updateData.Email
	}
	if updateData.HashedPassword != "" {
		if err := member.SetPassword(updateData.HashedPassword); err != nil {
			return err
		}
	}

	return mi.MemberRepository.Update(member)
}

func (mi *MemberInteractor) DeleteByUserID(userID string) error {
	member, err := mi.MemberRepository.GetByUserID(userID)
	if err != nil {
		return err
	}
	return mi.MemberRepository.Delete(member.ID)
}

func (mi *MemberInteractor) GetAll() ([]*domain.Member, error) {
	return mi.MemberRepository.GetAll()
}

func (mi *MemberInteractor) GetStatsByMonth(month string) (int, int, error) {
	return mi.MemberRepository.GetStatsByMonth(month)
}
