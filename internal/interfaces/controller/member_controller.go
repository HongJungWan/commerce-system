package controller

import (
	"net/http"

	"github.com/HongJungWan/commerce-system/internal/interfaces/dto/request"
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
	var req request.RegisterMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 요청입니다."})
		return
	}

	responseData, err := mc.memberInteractor.Register(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, responseData)
}

func (mc *MemberController) GetMyInfo(c *gin.Context) {
	username := c.GetString("username")
	responseData, err := mc.memberInteractor.GetMyInfo(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "사용자 정보를 가져올 수 없습니다."})
		return
	}

	c.JSON(http.StatusOK, responseData)
}

func (mc *MemberController) UpdateMyInfo(c *gin.Context) {
	username := c.GetString("username")
	var req request.UpdateMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 요청입니다."})
		return
	}

	if err := mc.memberInteractor.UpdateMyInfo(username, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "업데이트에 실패했습니다."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "정보가 수정되었습니다."})
}
func (mc *MemberController) DeleteMyAccount(c *gin.Context) {
	username := c.GetString("username")
	if err := mc.memberInteractor.DeleteByUserName(username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "계정 삭제에 실패했습니다."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "계정이 삭제되었습니다."})
}

func (mc *MemberController) GetAllMembers(c *gin.Context) {
	isAdmin := c.GetBool("is_admin")
	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "접근 권한이 없습니다."})
		return
	}
	responseData, err := mc.memberInteractor.GetAllMembers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "회원 목록을 가져올 수 없습니다."})
		return
	}

	c.JSON(http.StatusOK, responseData)
}

func (mc *MemberController) GetMemberStats(c *gin.Context) {
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
	stats, err := mc.memberInteractor.GetMemberStats(month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "통계 정보를 가져올 수 없습니다."})
		return
	}
	c.JSON(http.StatusOK, stats)
}
