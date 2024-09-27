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
	// ap.r.GET(rootPath+"/test", ap.h.HandlerTest)

	ap.r.HEAD(rootPath+"/{catch:*}", ap.h.HandlerOptHead)
	ap.r.OPTIONS(rootPath+"/{catch:*}", ap.h.HandlerOptHead)

	// auth
	ap.r.POST(rootPath+"/{auth}", ap.h.HandlerAuth) // login, register, refresh token, logout

	// user
	ap.r.GET(rootPath+"/user/{login}", ap.h.HandlerUser)                   // info self and info current user
	ap.r.PATCH(rootPath+"/user/{login}", ap.h.HandlerUser)                 // modification self user data
	ap.r.PATCH(rootPath+"/user/self/password", ap.h.HandlerChangePassword) // change password for current user
	ap.r.DELETE(rootPath+"/user/{login}", ap.h.HandlerUser)                // delete user (do not drop chat's if chat owner)

	// chat
	ap.r.GET(rootPath+"/chat/{desc}", ap.h.HandlerChat)    // chat info data or self-creator chats
	ap.r.POST(rootPath+"/chat/{desc}", ap.h.HandlerChat)   // create chat
	ap.r.DELETE(rootPath+"/chat/{desc}", ap.h.HandlerChat) // delete chat

	// follow
	ap.r.GET(rootPath+"/chat/{desc}/followers", ap.h.HandlerFollow) // get followers for current chat
	ap.r.POST(rootPath+"/chat/{desc}/follow", ap.h.HandlerFollow)   // follow chat
	ap.r.DELETE(rootPath+"/chat/{desc}/follow", ap.h.HandlerFollow) // unfollow chat

	// msg
	ap.r.GET(rootPath+"/chat/{desc}/message", ap.h.HandlerMessage)         // get messages
	ap.r.POST(rootPath+"/chat/{desc}/message", ap.h.HandlerMessage)        // push message
	ap.r.DELETE(rootPath+"/chat/{desc}/message/{id}", ap.h.HandlerMessage) // delete message

	// role (admin)
	ap.r.POST(rootPath+"/admin/{entity}", ap.h.HandlerAdmin)           // create user or chat
	ap.r.DELETE(rootPath+"/admin/user/{login}", ap.h.HandlerAdmin)     // delete user
	ap.r.PATCH(rootPath+"/admin/user/{login}/role", ap.h.HandlerAdmin) // set role for current user
	ap.r.GET(rootPath+"/admin/{entity}", ap.h.HandlerAdmin)            // get all users and all chats
}

func (ap *ApiRouter) GetHandler() func(ctx *fasthttp.RequestCtx) {
	return ap.r.Handler
}
