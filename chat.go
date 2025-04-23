package todo

type UserID struct {
	UserId int `json:"userID"`
}

type ChatId struct {
	ChatId int `json:"chatId"`
}

type UserVerification struct {
	UserId int `json:"userId"`
	ChatId int `uri:"id" binding:"required" json:"chatId"`
}

type ChatList struct {
	Id          int    `bd:"id"`
	Name        string `bd:"name"`
	CountNewMsg int    `bd:"countNewMsg"`
	LastMsgChat int    `bd:"lastMsgChat"`
}

type Chat struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	UserId int    `json:"owner"`
}

type SearchChat struct {
	Id     int `json:"id"`
	UserId int `json:"owner"`
}

type ChatName struct {
	Name string `json:"name"`
}

type SendMessage struct {
	Id      int    `json:"id" db:"id"`
	UserId  int    `json:"idUser"`
	ChatId  int    `json:"idChat"`
	Message string `json:"message" binding:"required"`
}

type ReadMessage struct {
	Id      int    `json:"id"`
	User    string `json:"user"`
	Message string `json:"message"`
	Data    string `json:"date_create"`
}

type UsersList struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	IsOwner bool   `json:"is_owner"`
}

type RenameChat struct {
	ChatId   int
	Owner    int
	ChatName string
}
