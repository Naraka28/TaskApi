package user

type User struct{
	Id int64 `json:"id"`
	Name string `json:"name"`
	Age int `json:"age"`
	Email string `json:"email"`
	Password string `json:"-"`
}
type RegisterUser struct {
	Name string `json:"name"`
	Age int `json:"age"`
	Email string `json:"email"`
	Password string `json:"password"`
}