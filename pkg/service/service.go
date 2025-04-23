package service

import (
	todo "rest_API"
	"rest_API/pkg/repository"
)

type Autorization interface {
	CreateUser(user todo.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type (
	Chat interface {
		CheckChatAndUser(userCheck todo.UserVerification) (access bool)
		AdminAccessVerification(userCheck todo.UserVerification) (access bool)
		LoadChatList(userId int) ([]todo.ChatList, error)
		CreateChat(chat todo.Chat) (todo.Chat, error)
		SearchChatByName(name string) ([]todo.SearchChat, error)
		SendMessage(message todo.SendMessage) (int, error)
		LoadNewMsgById(iuserId, chatId int) ([]todo.ReadMessage, error)
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
		LoadNick(UserId int) (string, error)
	}
)

type Service struct {
	Autorization
	Chat
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Autorization: NewAuthService(repos.Autorization),
		Chat:         NewChatService(repos.Chat),
	}
}
