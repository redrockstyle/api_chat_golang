package httprouter

import (
	httphandler "api_chat/api/layers/controller/server/http_handler"
	"api_chat/api/layers/domain/cfg"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

type ApiRouter struct {
	cfg *cfg.Configuration
	r   *router.Router
	h   *httphandler.HttpApiHandler
}

func NewApiRouter(cfg *cfg.Configuration, h *httphandler.HttpApiHandler) *ApiRouter {
	return &ApiRouter{cfg: cfg, r: router.New(), h: h}
}

func (ap *ApiRouter) InitRouter(rootPath string) {
	if ap.cfg.Server.Mode == "Development" {
		ap.r.GET(rootPath+"/test", ap.h.HandlerTest)
	}

	// auth
	ap.r.POST(rootPath+"/login", ap.h.HandlerLogin) // login
	if ap.cfg.Server.Registration {
		ap.r.POST(rootPath+"/register", ap.h.HanlderRegister) // register
	}
	ap.r.POST(rootPath+"/refresh", ap.h.HanlderRefresh) // refresh token

	// user
	ap.r.GET(rootPath+"/user/{login}", ap.h.HandlerUser)          // info self and info current user
	ap.r.GET(rootPath+"/user/{login}/{follow}", ap.h.HandlerUser) // info self follow chats
	ap.r.POST(rootPath+"/user/{login}", ap.h.HandlerUser)         // modification self user data
	ap.r.DELETE(rootPath+"/user/{login}", ap.h.HandlerUser)       // delete user (do not drop chat's if chat owner)

	// chat
	ap.r.GET(rootPath+"/chat/{desc}", ap.h.HandlerChat)    // chat info data or self-creator chats
	ap.r.POST(rootPath+"/chat/{desc}", ap.h.HandlerChat)   // create chat
	ap.r.DELETE(rootPath+"/chat/{desc}", ap.h.HandlerChat) // delete chat

	// follow
	ap.r.GET(rootPath+"/chat/{desc}/followers", ap.h.HandlerFollow)   // get followers for current chat
	ap.r.POST(rootPath+"/chat/{desc}/follow", ap.h.HandlerFollow)     // follow chat
	ap.r.DELETE(rootPath+"/chat/{desc}/unfollow", ap.h.HandlerFollow) // unfollow chat

	// msg
	ap.r.GET(rootPath+"/msg/{chat}", ap.h.HandlerMessage) // get messages
	ap.r.POST(rootPath+"/msg", ap.h.HandlerMessage)       // push message
	ap.r.DELETE(rootPath+"/msg", ap.h.HandlerMessage)     // delete message

	// role (admin)
	ap.r.POST(rootPath+"/admin/create/{entity}", ap.h.HandlerAdminCreate) // create user or chat
	ap.r.POST(rootPath+"/admin/{entity}/role", ap.h.HandlerAdmin)         // set role for current user
	ap.r.GET(rootPath+"/admin/{entity}/show", ap.h.HandlerAdmin)          // get all users and all chats
}

func (ap *ApiRouter) GetHander() func(ctx *fasthttp.RequestCtx) {
	return ap.r.Handler
}
