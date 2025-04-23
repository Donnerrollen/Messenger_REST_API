package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	todo "rest_API"
	"strconv"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
	userAccess          = "access"
	adminAccess         = "admin_access"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		NewErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		NewErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}

	userId, err := h.services.Autorization.ParseToken(headerParts[1])
	if err != nil {
		NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set(userCtx, userId)
}

// Функция проверяет есть ли этот пользователь в этом чате
func (h *Handler) checkChatAndUser(c *gin.Context) {
	var UserCheck todo.UserVerification

	ChatId := c.Param("id")
	UserCheck.ChatId, _ = strconv.Atoi(ChatId)

	id, ok := c.Get(userCtx)
	if !ok {
		NewErrorResponse(c, http.StatusInternalServerError, "id users not found")
	}

	UserCheck.UserId = id.(int)

	access := h.services.Chat.CheckChatAndUser(UserCheck)

	c.Set(userAccess, access)
}

func (h *Handler) adminAccessVerification(c *gin.Context) {
	var AdminCheck todo.UserVerification

	ChatId := c.Param("id")
	AdminCheck.ChatId, _ = strconv.Atoi(ChatId)

	id, ok := c.Get(userCtx)
	if !ok {
		NewErrorResponse(c, http.StatusInternalServerError, "id users not found")
	}

	AdminCheck.UserId = id.(int)

	admin_access := h.services.Chat.AdminAccessVerification(AdminCheck)

	c.Set(adminAccess, admin_access)
}
