package response

type MemberResponse struct {
	ID           uint   `json:"id"`
	MemberNumber string `json:"member_number"`
	Username     string `json:"username"`
	FullName     string `json:"full_name"`
	Email        string `json:"email"`
	CreatedAt    string `json:"created_at"`
	IsAdmin      bool   `json:"is_admin"`
	IsWithdrawn  bool   `json:"is_withdrawn"`
	WithdrawnAt  string `json:"withdrawn_at,omitempty"`
}

type RegisterMemberResponse struct {
	Message string          `json:"message"`
	Token   string          `json:"token"`
	User    *MemberResponse `json:"user"`
}
