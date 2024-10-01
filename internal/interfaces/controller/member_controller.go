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

// Register godoc
// @Summary      회원 가입
// @Description  새로운 회원을 등록합니다.
// @Tags         members
// @Accept       json
// @Produce      json
// @Param        registerRequest body request.RegisterMemberRequest true "회원 가입 정보"
// @Success      201 {object} response.MemberResponse "가입 성공"
// @Failure      400 {object} map[string]string "잘못된 요청"
// @Failure      500 {object} map[string]string "서버 오류"
// @Router       /members [post]
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

// GetMyInfo godoc
// @Summary      내 정보 조회
// @Description  인증된 사용자의 정보를 조회합니다.
// @Tags         members
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Success      200 {object} response.MemberResponse "내 정보"
// @Failure      500 {object} map[string]string "정보 조회 실패"
// @Router       /members/me [get]
func (mc *MemberController) GetMyInfo(c *gin.Context) {
	username := c.GetString("username")
	responseData, err := mc.memberInteractor.GetMyInfo(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "사용자 정보를 가져올 수 없습니다."})
		return
	}

	c.JSON(http.StatusOK, responseData)
}

// UpdateMyInfo godoc
// @Summary      내 정보 수정
// @Description  인증된 사용자의 정보를 수정합니다.
// @Tags         members
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        updateRequest body request.UpdateMemberRequest true "수정할 정보"
// @Success      200 {object} map[string]string "수정 성공"
// @Failure      400 {object} map[string]string "잘못된 요청"
// @Failure      500 {object} map[string]string "수정 실패"
// @Router       /members/me [put]
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

// DeleteMyAccount godoc
// @Summary      회원 탈퇴
// @Description  인증된 사용자의 계정을 삭제합니다.
// @Tags         members
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]string "삭제 성공"
// @Failure      500 {object} map[string]string "삭제 실패"
// @Router       /members/me [delete]
func (mc *MemberController) DeleteMyAccount(c *gin.Context) {
	username := c.GetString("username")
	if err := mc.memberInteractor.DeleteByUserName(username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "계정 삭제에 실패했습니다."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "계정이 삭제되었습니다."})
}

// GetAllMembers godoc
// @Summary      회원 목록 조회
// @Description  모든 회원의 목록을 조회합니다. (관리자 전용)
// @Tags         members
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Success      200 {array} response.MemberResponse "회원 목록"
// @Failure      403 {object} map[string]string "권한 없음"
// @Failure      500 {object} map[string]string "목록 조회 실패"
// @Router       /members [get]
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

// GetMemberStats godoc
// @Summary      회원 통계 조회
// @Description  특정 월의 회원 가입 통계를 조회합니다. (관리자 전용)
// @Tags         members
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        month query string true "조회할 월 (YYYY-MM)"
// @Success      200 {object} response.MemberStatsResponse "통계 정보"
// @Failure      400 {object} map[string]string "잘못된 요청"
// @Failure      403 {object} map[string]string "권한 없음"
// @Failure      500 {object} map[string]string "통계 조회 실패"
// @Router       /members/stats [get]
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
