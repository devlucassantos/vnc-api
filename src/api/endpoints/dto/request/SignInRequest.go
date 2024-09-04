package request

type SignIn struct {
	Email    string `json:"email"    example:"example@email.com"`
	Password string `json:"password" example:"luc@ssantos123"`
}
