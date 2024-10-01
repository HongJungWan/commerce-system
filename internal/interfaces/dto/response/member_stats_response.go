package response

type MemberStatsResponse struct {
	Month          string `json:"month"`
	JoinedMembers  int    `json:"joined_members"`
	DeletedMembers int    `json:"deleted_members"`
}
