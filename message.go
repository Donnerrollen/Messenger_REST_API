package todo

type Message struct {
	id          int    `json:"-" db:"id"`
	title       string `json:"title"`
	description string `json:"description"`
}

type UserMessage struct {
	Id        int
	UserId    int
	MessageId int
}

type UpdateMessage struct {
	idChat   string `json:"idChat"`
	messsage string `json:"messsage"`
}
