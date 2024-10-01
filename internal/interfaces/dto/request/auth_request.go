package request

type LoginRequest struct {
	Username       string `json:"username"`
	HashedPassword string `json:"hashed_password"`
}
