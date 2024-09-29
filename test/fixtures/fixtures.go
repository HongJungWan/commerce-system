package fixtures

import (
	"github.com/HongJungWan/commerce-system/internal/domain"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// 테스트용 인메모리 데이터베이스 설정
func SetupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("테스트 데이터베이스 연결에 실패했습니다.")
	}

	// 스키마 마이그레이션: 필요한 모든 도메인 모델 추가
	err = db.AutoMigrate(&domain.Member{}, &domain.Product{}, &domain.Order{})
	if err != nil {
		panic("테스트 데이터베이스 마이그레이션에 실패했습니다.")
	}

	return db
}
