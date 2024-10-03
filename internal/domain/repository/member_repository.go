package repository

import "github.com/HongJungWan/commerce-system/internal/domain"

type MemberRepository interface {
	Create(member *domain.Member) error
	GetByID(id uint) (*domain.Member, error)
	GetByAccountId(userName string) (*domain.Member, error)
	GetByMemberNumber(memberNumber string) (*domain.Member, error)
	Update(member *domain.Member) error
	Delete(id uint) error
	GetAll() ([]*domain.Member, error)
	GetStatsByMonth(month string) (int, int, error)
}
