package repository

import (
	"github.com/jackc/pgx/v4"
	todo "rest_API"
)

type Autorization interface {
	CreateUser(user todo.User) (int, error)
	GetUser(username, password string) (todo.User, error)
}

type Chat interface {
	CheckChatAndUser(userCheck todo.UserVerification) (access bool)
	AdminAccessVerification(userCheck todo.UserVerification) (access bool)
	LoadChatList(userId int) ([]todo.ChatList, error)
	CreateChat(chat todo.Chat) (todo.Chat, error)
	SearchChatByName(name string) ([]todo.SearchChat, error)
	SendMessage(message todo.SendMessage) (int, error)
	LoadNewMsgById(userId, chatId int) ([]todo.ReadMessage, error)
	LoadOldMsgById(chatId, msgId int) ([]todo.ReadMessage, error)
	CheckUserInSystem(userId int) (bool, error)
	IsUserInChat(chatId, userId int) (bool, error)
	CreateInvite(chatId, userId int) (int, error)
	GetUsersOfChat(chatId int) ([]todo.UsersList, error)
	AcceptInvite(userId, chatId int) (int, error)
	DenyInvite(UserId, ChatId int) (int, error)
	LogOut(UserId, ChatId int) (int, error)
	DeleteUser(UserId, ChatId int) (int, error)
	RenameChat(RenameChat todo.RenameChat) error
	DeleteChat(ChatId int) error
	GiveAdminStatus(ChatId, UserId int) error
	LoadNick(UserId int) (string, int, error)
}

type Repository struct {
	Autorization
	Chat
}

func NewRepository(db *pgx.Conn) *Repository {
	return &Repository{
		Autorization: NewAuthPostgres(db),
		Chat:         NewChatPostgres(db),
	}
}
