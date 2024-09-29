package controller

import (
	"net/http"

	"github.com/HongJungWan/commerce-system/internal/domain"
	"github.com/HongJungWan/commerce-system/internal/usecases"
	"github.com/gin-gonic/gin"
)

type MemberController struct {
	memberInteractor *usecases.MemberInteractor
	authUseCase      *usecases.AuthUseCase
}

func NewMemberController(mi *usecases.MemberInteractor, au *usecases.AuthUseCase) *MemberController {
	return &MemberController{
		memberInteractor: mi,
		authUseCase:      au,
	}
}

func (mc *MemberController) Register(c *gin.Context) {
	var member domain.Member
	if err := c.ShouldBindJSON(&member); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 요청입니다."})
		return
	}
	if err := mc.memberInteractor.Register(&member); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	token, _ := mc.authUseCase.GenerateToken(member.UserID, member.IsAdmin)
	c.JSON(http.StatusCreated, gin.H{"message": "회원 가입이 완료되었습니다.", "token": token})
}

func (mc *MemberController) GetMyInfo(c *gin.Context) {
	userID := c.GetString("user_id")
	member, err := mc.memberInteractor.GetByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "사용자 정보를 가져올 수 없습니다."})
		return
	}
	c.JSON(http.StatusOK, member)
}

func (mc *MemberController) UpdateMyInfo(c *gin.Context) {
	userID := c.GetString("user_id")
	var updateData domain.Member
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 요청입니다."})
		return
	}
	if err := mc.memberInteractor.UpdateByUserID(userID, &updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "업데이트에 실패했습니다."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "정보가 수정되었습니다."})
}

func (mc *MemberController) DeleteMyAccount(c *gin.Context) {
	userID := c.GetString("user_id")
	if err := mc.memberInteractor.DeleteByUserID(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "계정 삭제에 실패했습니다."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "계정이 삭제되었습니다."})
}

func (mc *MemberController) GetAllMembers(c *gin.Context) {
	// FIXME: DB 찔러서 가져오도록 수정하기.
	isAdmin := c.GetBool("is_admin")
	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "접근 권한이 없습니다."})
		return
	}
	members, err := mc.memberInteractor.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "회원 목록을 가져올 수 없습니다."})
		return
	}
	c.JSON(http.StatusOK, members)
}

func (mc *MemberController) GetMemberStats(c *gin.Context) {
	// FIXME: DB 찔러서 가져오도록 수정하기.
	isAdmin := c.GetBool("is_admin")
	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "접근 권한이 없습니다."})
		return
	}
	month := c.Query("month")
	if month == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "month 파라미터가 필요합니다."})
		return
	}
	joined, deleted, err := mc.memberInteractor.GetStatsByMonth(month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "통계 정보를 가져올 수 없습니다."})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"month":           month,
		"joined_members":  joined,
		"deleted_members": deleted,
	})
}
