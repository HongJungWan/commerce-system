package controller

import (
	"net/http"
	"time"

	"github.com/HongJungWan/commerce-system/internal/domain"
	"github.com/HongJungWan/commerce-system/internal/interfaces/dto/request"
	"github.com/HongJungWan/commerce-system/internal/interfaces/dto/response"
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

	member := &domain.Member{
		MemberNumber: req.MemberNumber,
		Username:     req.Username,
		FullName:     req.FullName,
		Email:        req.Email,
		CreatedAt:    time.Now(),
	}

	if err := member.SetPassword(req.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "비밀번호 설정에 실패했습니다."})
		return
	}

	if err := mc.memberInteractor.Register(member); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	token, _ := mc.authUseCase.GenerateToken(member)

	responseData := response.RegisterMemberResponse{
		Message: "회원 가입이 완료되었습니다.",
		Token:   token,
		User: response.MemberResponse{
			ID:           member.ID,
			MemberNumber: member.MemberNumber,
			Username:     member.Username,
			FullName:     member.FullName,
			Email:        member.Email,
			CreatedAt:    member.CreatedAt.Format(time.RFC3339),
			IsAdmin:      member.IsAdmin,
			IsWithdrawn:  member.IsWithdrawn,
			WithdrawnAt:  formatTime(member.WithdrawnAt),
		},
	}

	c.JSON(http.StatusCreated, responseData)
}

func (mc *MemberController) GetMyInfo(c *gin.Context) {
	username := c.GetString("username")
	member, err := mc.memberInteractor.GetByUserName(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "사용자 정보를 가져올 수 없습니다."})
		return
	}

	responseData := response.MemberResponse{
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

	c.JSON(http.StatusOK, responseData)
}

func (mc *MemberController) UpdateMyInfo(c *gin.Context) {
	username := c.GetString("username")
	var updateData request.UpdateMemberRequest
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 요청입니다."})
		return
	}

	memberUpdate := &domain.Member{
		FullName: updateData.FullName,
		Email:    updateData.Email,
	}

	if updateData.Password != "" {
		if err := memberUpdate.SetPassword(updateData.Password); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "비밀번호 설정에 실패했습니다."})
			return
		}
	}

	if err := mc.memberInteractor.UpdateByUserName(username, memberUpdate); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "업데이트에 실패했습니다."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "정보가 수정되었습니다."})
}

func (mc *MemberController) DeleteMyAccount(c *gin.Context) {
	userName := c.GetString("username")
	if err := mc.memberInteractor.DeleteByUserName(userName); err != nil {
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
	members, err := mc.memberInteractor.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "회원 목록을 가져올 수 없습니다."})
		return
	}

	var memberResponses []response.MemberResponse
	for _, member := range members {
		memberResponses = append(memberResponses, response.MemberResponse{
			ID:           member.ID,
			MemberNumber: member.MemberNumber,
			Username:     member.Username,
			FullName:     member.FullName,
			Email:        member.Email,
			CreatedAt:    member.CreatedAt.Format(time.RFC3339),
			IsAdmin:      member.IsAdmin,
			IsWithdrawn:  member.IsWithdrawn,
			WithdrawnAt:  formatTime(member.WithdrawnAt),
		})
	}

	c.JSON(http.StatusOK, memberResponses)
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

func formatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}
