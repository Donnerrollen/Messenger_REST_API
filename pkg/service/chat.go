package service

import (
	todo "rest_API"
	"rest_API/pkg/repository"
	"strconv"
)

type ChatService struct {
	repo repository.Chat
}

func NewChatService(repo repository.Chat) *ChatService {
	return &ChatService{repo: repo}
}

func (s *ChatService) LoadChatList(id int) ([]todo.ChatList, error) {
	return s.repo.LoadChatList(id)
}

func (s *ChatService) CreateChat(chat todo.Chat) (todo.Chat, error) {
	return s.repo.CreateChat(chat)
}

func (s *ChatService) SearchChatByName(name string) ([]todo.SearchChat, error) {
	return s.repo.SearchChatByName(name)
}

func (s *ChatService) CheckChatAndUser(userCheck todo.UserVerification) (access bool) {
	return s.repo.CheckChatAndUser(userCheck)
}

func (s *ChatService) SendMessage(message todo.SendMessage) (int, error) {
	return s.repo.SendMessage(message)
}

func (s *ChatService) AdminAccessVerification(userCheck todo.UserVerification) (access bool) {
	return s.repo.AdminAccessVerification(userCheck)
}

func (s *ChatService) LoadNewMsgById(userId, chatId int) ([]todo.ReadMessage, error) {
	return s.repo.LoadNewMsgById(userId, chatId)
}

func (s *ChatService) LoadOldMsgById(chatId, msgId int) ([]todo.ReadMessage, error) {
	return s.repo.LoadOldMsgById(chatId, msgId)
}

func (s *ChatService) CheckUserInSystem(userId int) (bool, error) {
	return s.repo.CheckUserInSystem(userId)
}

func (s *ChatService) IsUserInChat(chatId, userId int) (bool, error) {
	return s.repo.IsUserInChat(chatId, userId)
}

func (s *ChatService) CreateInvite(chatId, userId int) (int, error) {
	return s.repo.CreateInvite(chatId, userId)
}

func (s *ChatService) GetUsersOfChat(chatId int) ([]todo.UsersList, error) {
	return s.repo.GetUsersOfChat(chatId)
}

func (s *ChatService) AcceptInvite(userId, chatId int) (int, error) {
	return s.repo.AcceptInvite(userId, chatId)
}

func (s *ChatService) DenyInvite(UserId, ChatId int) (int, error) {
	return s.repo.DenyInvite(UserId, ChatId)
}

func (s *ChatService) LogOut(UserId, ChatId int) (int, error) {
	return s.repo.LogOut(UserId, ChatId)
}

func (s *ChatService) DeleteUser(UserId, ChatId int) (int, error) {
	return s.repo.DeleteUser(UserId, ChatId)
}

func (s *ChatService) RenameChat(RenameChat todo.RenameChat) error {
	return s.repo.RenameChat(RenameChat)
}

func (s *ChatService) DeleteChat(ChatId int) error {
	return s.repo.DeleteChat(ChatId)
}

func (s *ChatService) GiveAdminStatus(ChatId, UserId int) error {
	return s.repo.GiveAdminStatus(ChatId, UserId)
}

func (s *ChatService) LoadNick(UserId int) (string, error) {
	Name, Id, err := s.repo.LoadNick(UserId)
	if err != nil {
		return " ", err
	}
	Name = Name + "#" + strconv.Itoa(Id)
	return Name, nil
}
