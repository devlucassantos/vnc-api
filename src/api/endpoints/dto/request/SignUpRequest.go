package request

type SignUp struct {
	FirstName string `json:"first_name" example:"Lucas"`
	LastName  string `json:"last_name"  example:"Santos"`
	Email     string `json:"email"      example:"example@email.com"`
	Password  string `json:"password"   example:"luc@ssantos123"`
}
