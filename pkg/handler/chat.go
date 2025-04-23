package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	todo "rest_API"
	"strconv"
)

func (h *Handler) loadChatList(c *gin.Context) {
	var list []todo.ChatList
	id, ok := c.Get(userCtx)

	if !ok {
		NewErrorResponse(c, http.StatusInternalServerError, "id users not found")
		return
	}

	userId := id.(int)

	list, err := h.services.Chat.LoadChatList(userId)

	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"chats": list,
	})
}

func (h *Handler) createChat(c *gin.Context) {
	var input todo.Chat

	id, ok := c.Get(userCtx)

	if !ok {
		NewErrorResponse(c, http.StatusInternalServerError, "id users not found")
	}

	err := c.BindJSON(&input)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	input.UserId = id.(int)

	input, _ = h.services.Chat.CreateChat(input)

	c.JSON(http.StatusOK, map[string]interface{}{
		"chat": input,
	})
}

func (h *Handler) searchChatByName(c *gin.Context) {
	var chat []todo.SearchChat
	var Cname todo.ChatName
	_, ok := c.Get(userCtx)

	if !ok {
		NewErrorResponse(c, http.StatusInternalServerError, "id users not found")
		return
	}

	err := c.BindJSON(&Cname)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	name := Cname.Name
	chat, _ = h.services.Chat.SearchChatByName(name)

	c.JSON(http.StatusOK, map[string]interface{}{
		"chat": chat,
	})
}

func (h *Handler) loadNewMsgById(c *gin.Context) {
	var message []todo.ReadMessage

	access, ok := c.Get(userAccess)
	if !ok {
		NewErrorResponse(c, http.StatusInternalServerError, "id chats not found")
		return
	}

	id_user, ok := c.Get(userCtx)
	if !ok {
		NewErrorResponse(c, http.StatusInternalServerError, "id user not found")
		return
	}

	if access == false {
		NewErrorResponse(c, http.StatusUnauthorized, "you have not access to the chat")
		return
	}
	ChatId := c.Param("id")
	ChatId_int, _ := strconv.Atoi(ChatId)

	message, err := h.services.Chat.LoadNewMsgById(id_user.(int), ChatId_int)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "Error")
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"messages": message,
	})
}

func (h *Handler) loadOldMsgById(c *gin.Context) {
	var message []todo.ReadMessage

	access, ok := c.Get(userAccess)
	if !ok {
		NewErrorResponse(c, http.StatusInternalServerError, "id chats not found")
		return
	}

	if access == false {
		NewErrorResponse(c, http.StatusUnauthorized, "you have not access to the chat")
		return
	}
	ChatId := c.Param("id")
	ChatId_int, _ := strconv.Atoi(ChatId)

	id_msg := c.Param("id_msg")
	id_msg_int, _ := strconv.Atoi(id_msg)

	message, err := h.services.Chat.LoadOldMsgById(ChatId_int, id_msg_int)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "Error")
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"messages": message,
	})
}

func (h *Handler) sendMessage(c *gin.Context) {
	var message todo.SendMessage

	id, ok := c.Get(userCtx)
	if !ok {
		NewErrorResponse(c, http.StatusInternalServerError, "id users not found")
	}

	access, ok := c.Get(userAccess)
	if !ok {
		NewErrorResponse(c, http.StatusInternalServerError, "id chats not found")
		return
	}

	if access == false {
		NewErrorResponse(c, http.StatusUnauthorized, "you have not access to the chat")
		return
	}

	err := c.BindJSON(&message)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	ChatId := c.Param("id")
	message.ChatId, _ = strconv.Atoi(ChatId)
	message.UserId = id.(int)

	message.Id, _ = h.services.SendMessage(message)

	c.JSON(http.StatusOK, map[string]interface{}{
		"messageID": message.Id,
	})
}

func (h *Handler) updateMessage(c *gin.Context) {

}

func (h *Handler) deleteMessage(c *gin.Context) {

}

func (h *Handler) inviteUser(c *gin.Context) {

	var id_user todo.UserID
	var msg string

	access, ok := c.Get(userAccess)
	if !ok {
		NewErrorResponse(c, http.StatusInternalServerError, "id chats not found")
		return
	}

	if access == false {
		NewErrorResponse(c, http.StatusUnauthorized, "you have not access to the chat")
		return
	}

	chatId, _ := strconv.Atoi(c.Param("id"))

	err := c.BindJSON(&id_user)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	UserCheck, _ := h.services.Chat.CheckUserInSystem(id_user.UserId)
	if !UserCheck {
		NewErrorResponse(c, http.StatusBadRequest, "This user does not exist")
		return
	}

	isUserInChat, err := h.services.Chat.IsUserInChat(chatId, id_user.UserId)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if isUserInChat {
		c.JSON(http.StatusOK, map[string]interface{}{
			"message": "user in the chat",
		})
		return
	}

	status, err := h.services.Chat.CreateInvite(chatId, id_user.UserId)
	if (err != nil) || (status == 2) {
		NewErrorResponse(c, http.StatusBadRequest, "failed to send invitation")
		return
	}

	if status == 0 {
		msg = "the invitation has been sent"
	} else if status == 1 {
		msg = "the invitation is exists"
	} else {
		NewErrorResponse(c, http.StatusBadRequest, "failed to send invitation")
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": msg,
	})

}

func (h *Handler) deleteUser(c *gin.Context) {
	var User todo.UserID
	var ChatId int
	var msg string

	id, ok := c.Get(userCtx)
	if !ok {
		NewErrorResponse(c, http.StatusInternalServerError, "id users not found")
	}

	admin_access, _ := c.Get(adminAccess)
	if !admin_access.(bool) {
		NewErrorResponse(c, http.StatusBadRequest, "You have not admin status")
		return
	}

	err := c.BindJSON(&User)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if id == User.UserId {
		NewErrorResponse(c, http.StatusBadRequest, "You can not delete yourself")
		return
	}
	ChatId, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	status, err := h.services.Chat.DeleteUser(User.UserId, ChatId)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if status == 2 {
		msg = "This user is not in the chat"
	}

	if status == 0 {
		msg = "success"
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"data": msg,
	})
}

func (h *Handler) renameChat(c *gin.Context) {
	var RenameChat todo.RenameChat
	var ChatName todo.ChatName

	id, ok := c.Get(userCtx)
	RenameChat.Owner = id.(int)

	admin_access, ok := c.Get(adminAccess)
	if !ok {
		NewErrorResponse(c, http.StatusBadRequest, "id user or id chat not found")
		return
	}
	if !admin_access.(bool) {
		NewErrorResponse(c, http.StatusBadRequest, "You have not admin status")
		return
	}

	err := c.BindJSON(&ChatName)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	RenameChat.ChatName = ChatName.Name

	RenameChat.ChatId, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Chat.RenameChat(RenameChat)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"data": "success",
	})
}

func (h *Handler) acceptInvite(c *gin.Context) {
	var Invite todo.UserVerification
	var msg string

	id, ok := c.Get(userCtx)
	Invite.UserId = id.(int)

	err := c.BindJSON(&Invite.ChatId)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if !ok {
		NewErrorResponse(c, http.StatusInternalServerError, "id users not found")
		return
	}

	status, err := h.services.Chat.AcceptInvite(Invite.UserId, Invite.ChatId)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if status == 1 {
		msg = "You have not this invitation"
	}

	if status == 0 {
		msg = "you join to the chat"
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": msg,
	})
}

func (h *Handler) denyInvite(c *gin.Context) {
	var Invite todo.UserVerification
	var ChatId todo.ChatId
	var msg string

	id, ok := c.Get(userCtx)
	Invite.UserId = id.(int)

	err := c.BindJSON(&ChatId)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	Invite.ChatId = ChatId.ChatId

	if !ok {
		NewErrorResponse(c, http.StatusInternalServerError, "id users not found")
	}

	status, err := h.services.Chat.DenyInvite(Invite.UserId, Invite.ChatId)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if status == 1 {
		msg = "You have not this invitation"
	}

	if status == 0 {
		msg = "you declined the invitation"
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": msg,
	})
}

func (h *Handler) getUsersOfChat(c *gin.Context) {

	var listUsers []todo.UsersList

	_, ok := c.Get(userCtx)
	if !ok {
		NewErrorResponse(c, http.StatusInternalServerError, "id users not found")
		return
	}

	access, ok := c.Get(userAccess)
	if !ok {
		NewErrorResponse(c, http.StatusInternalServerError, "id chats not found")
		return
	}

	if access == false {
		NewErrorResponse(c, http.StatusUnauthorized, "you have not access to the chat")
		return
	}

	chatId, _ := strconv.Atoi(c.Param("id"))

	listUsers, err := h.services.Chat.GetUsersOfChat(chatId)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "Error")
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"users": listUsers,
	})

}

func (h *Handler) LogOut(c *gin.Context) {
	var UserLogOut todo.UserVerification
	var status int
	var msg string

	id, ok := c.Get(userCtx)
	if !ok {
		NewErrorResponse(c, http.StatusInternalServerError, "id users not found")
		return
	}

	ChatId := c.Param("id")

	UserLogOut.UserId = id.(int)
	UserLogOut.ChatId, _ = strconv.Atoi(ChatId)

	status, err := h.services.LogOut(UserLogOut.UserId, UserLogOut.ChatId)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if status == 1 {
		msg = "You are admin, you can not delete yourself"
	}
	if status == 0 {
		msg = "success"
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"data": msg,
	})
}

func (h *Handler) deleteChat(c *gin.Context) {
	admin_access, ok := c.Get(adminAccess)
	if !ok {
		NewErrorResponse(c, http.StatusBadRequest, "id user or id chat not found")
		return
	}
	if !admin_access.(bool) {
		NewErrorResponse(c, http.StatusBadRequest, "You have not admin status")
		return
	}

	ChatId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Chat.DeleteChat(ChatId)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"data": "success",
	})
}

func (h *Handler) giveAdminStatus(c *gin.Context) {
	var User todo.UserID

	admin_access, ok := c.Get(adminAccess)
	if !ok {
		NewErrorResponse(c, http.StatusBadRequest, "id user or id chat not found")
		return
	}
	if !admin_access.(bool) {
		NewErrorResponse(c, http.StatusBadRequest, "You have not admin status")
		return
	}

	ChatId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = c.BindJSON(&User)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Chat.GiveAdminStatus(ChatId, User.UserId)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"data": "success",
	})
}

func (h *Handler) loadNick(c *gin.Context) {
	id, ok := c.Get(userCtx)
	if !ok {
		NewErrorResponse(c, http.StatusInternalServerError, "id users not found")
		return
	}
	Nick, err := h.services.Chat.LoadNick(id.(int))
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"nickname": Nick,
	})
}
