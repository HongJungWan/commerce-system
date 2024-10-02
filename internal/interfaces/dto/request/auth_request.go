package request

type LoginRequest struct {
	AccountId string `json:"account_id"`
	Password  string `json:"password"`
}
