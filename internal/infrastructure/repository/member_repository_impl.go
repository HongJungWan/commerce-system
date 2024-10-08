package repository

import (
	"time"

	"github.com/HongJungWan/commerce-system/internal/domain"
	"gorm.io/gorm"
)

type MemberRepositoryImpl struct {
	db *gorm.DB
}

func NewMemberRepository(db *gorm.DB) *MemberRepositoryImpl {
	return &MemberRepositoryImpl{db: db}
}

func (r *MemberRepositoryImpl) Create(member *domain.Member) error {
	return r.db.Create(member).Error
}

func (r *MemberRepositoryImpl) GetByID(id uint) (*domain.Member, error) {
	var member domain.Member
	if err := r.db.First(&member, id).Error; err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *MemberRepositoryImpl) GetByAccountId(userName string) (*domain.Member, error) {
	var member domain.Member
	if err := r.db.Where("account_id = ?", userName).First(&member).Error; err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *MemberRepositoryImpl) GetByMemberNumber(memberNumber string) (*domain.Member, error) {
	var member domain.Member
	if err := r.db.Where("member_number = ?", memberNumber).First(&member).Error; err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *MemberRepositoryImpl) Update(member *domain.Member) error {
	return r.db.Save(member).Error
}

func (r *MemberRepositoryImpl) Delete(id uint) error {
	member, err := r.GetByID(id)
	if err != nil {
		return err
	}

	now := time.Now()
	member.IsWithdrawn = true
	member.WithdrawnAt = &now

	return r.Update(member)
}

func (r *MemberRepositoryImpl) GetAll() ([]*domain.Member, error) {
	var members []*domain.Member
	if err := r.db.Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

func (r *MemberRepositoryImpl) GetStatsByMonth(month string) (int, int, error) {
	var joinedCount int64
	var deletedCount int64

	startDate, err := time.Parse("2006-01", month)
	if err != nil {
		return 0, 0, err
	}
	endDate := startDate.AddDate(0, 1, 0)

	// 가입한 멤버 수
	if err := r.db.Model(&domain.Member{}).
		Where("created_at >= ? AND created_at < ?", startDate, endDate).
		Count(&joinedCount).Error; err != nil {
		return 0, 0, err
	}

	// 삭제된 멤버 수
	if err := r.db.Model(&domain.Member{}).
		Where("is_withdrawn = ? AND withdrawn_at >= ? AND withdrawn_at < ?", true, startDate, endDate).
		Count(&deletedCount).Error; err != nil {
		return 0, 0, err
	}

	return int(joinedCount), int(deletedCount), nil
}
