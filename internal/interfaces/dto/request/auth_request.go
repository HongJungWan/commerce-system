package request

type LoginRequest struct {
	AccountId string `json:"account_id" example:"hong43ok"`
	Password  string `json:"password" example:"ghdwjddhks"`
}
