package request

type RegisterMemberRequest struct {
	MemberNumber string `json:"member_number"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	FullName     string `json:"full_name"`
	Email        string `json:"email"`
}

type UpdateMemberRequest struct {
	FullName string `json:"full_name,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}
