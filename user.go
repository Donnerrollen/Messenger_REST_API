package todo

type User struct {
	Id       int    `json:"-" db:"id"`
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"` //Это логин, использующийся при входе пользователя
	Password string `json:"password" binding:"required"`
}
