package handler

import (
	"github.com/gin-gonic/gin"
	"rest_API/pkg/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	api := router.Group("/api", h.userIdentity)
	{
		messenger := api.Group("/chat")
		{
			messenger.GET("/load_user_nick", h.loadNick)
			messenger.GET("/load_chat_list", h.loadChatList)
			//			messenger.GET("/load_chat", h.loadChatById)
			messenger.POST("/create_chat", h.createChat)
			messenger.GET("/search_chat_byName", h.searchChatByName)
			messenger.POST("/accept_invite", h.acceptInvite)
			messenger.DELETE("/deny_invite", h.denyInvite)

			chat := messenger.Group("/:id", h.checkChatAndUser)
			{
				chat.GET("/load_chat_new_messages", h.loadNewMsgById)
				chat.GET("/load_chat_old_messages=:id_msg", h.loadOldMsgById)
				chat.POST("/send_message", h.sendMessage)
				chat.GET("/get_users_of_chat", h.getUsersOfChat)
				chat.DELETE("/log_out_of_the_chat", h.LogOut)
				chat.POST("/invite_user", h.inviteUser)
				//chat.PATCH("/update_message", h.updateMessage)
				//chat.DELETE("/delete_message", h.deleteMessage)

				admin := chat.Group("/admin", h.adminAccessVerification)
				{
					admin.DELETE("/delete_user", h.deleteUser)
					admin.PATCH("/rename_chat", h.renameChat)
					admin.DELETE("/delete_chat", h.deleteChat)
					admin.PATCH("/give_admin_status", h.giveAdminStatus)
				}
			}
		}
	}
	return router
}
