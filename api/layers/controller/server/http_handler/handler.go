package httphandler

import (
	"api_chat/api/layers/base/logx"
	"api_chat/api/layers/domain/cfg"
	"api_chat/api/layers/domain/db"
	"api_chat/api/layers/usecase"
	"encoding/json"
	"strconv"

	"github.com/valyala/fasthttp"
)

const (
	TargetSelf       = "self"
	strHeaderSession = "X-Allow-Session"

	//strContentType     = "Content-Type"
	strApplicationJSON = "application/json; charset=UTF-8"
)

type ReMessage struct {
	Chat    string `json:"chat"`
	Text    string `json:"text"`
	Message uint64 `json:"msg_id"`
}

type HttpApiHandler struct {
	cfg  *cfg.Configuration
	uc   *usecase.UsecaseOperator
	logx logx.Logger
}

func NewRestApiHandler(cfg *cfg.Configuration, uc *usecase.UsecaseOperator, logx logx.Logger) *HttpApiHandler {
	return &HttpApiHandler{cfg: cfg, uc: uc, logx: logx}
}

func (has *HttpApiHandler) HandlerTest(ctx *fasthttp.RequestCtx) {
	SetServerNameHeader(ctx, has.cfg.Server.ServerName)
	ctx.SetBody([]byte("Hello API"))
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func (has *HttpApiHandler) HandlerLogin(ctx *fasthttp.RequestCtx) {
	SetServerNameHeader(ctx, has.cfg.Server.ServerName)
	if string(ctx.Method()) == fasthttp.MethodPost {
		user := db.User{}
		if err := json.Unmarshal(ctx.Request.Body(), &user); err != nil {
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		} else {
			session, err := has.uc.Login(&user)
			if err != nil {
				has.logx.Warnf("Login user:%v err:%v", user.Login, err)
				ctx.SetStatusCode(fasthttp.StatusBadRequest)
			} else {
				ctx.Response.Header.Set(strHeaderSession, session)
				ctx.SetStatusCode(fasthttp.StatusOK)
			}
		}
	} else {
		ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
	}
}

func (has *HttpApiHandler) HanlderRegister(ctx *fasthttp.RequestCtx) {
	SetServerNameHeader(ctx, has.cfg.Server.ServerName)
	if has.cfg.Server.Registration {
		if string(ctx.Method()) == fasthttp.MethodPost {
			user := db.User{}
			if err := json.Unmarshal(ctx.Request.Body(), &user); err != nil {
				ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			} else {
				if err := has.uc.Register(&user); err != nil {
					has.logx.Warnf("Register user:%v err:%v", user.Login, err)
					ctx.SetStatusCode(fasthttp.StatusBadRequest)
				} else {
					SetLocationHeader(ctx, has.cfg.Server.PrefixPath+"/user/"+user.Login)
					ctx.SetStatusCode(fasthttp.StatusCreated)
				}
			}
		} else {
			ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
		}
	} else {
		ctx.SetStatusCode(fasthttp.StatusForbidden)
	}
}

func (has *HttpApiHandler) HanlderRefresh(ctx *fasthttp.RequestCtx) {
	SetServerNameHeader(ctx, has.cfg.Server.ServerName)
	if string(ctx.Method()) == fasthttp.MethodPost {
		session := ctx.Request.Header.Peek(strHeaderSession)
		if len(session) == 0 {
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		} else {
			newSession, err := has.uc.Refresh(string(session))
			if err != nil {
				has.logx.Warnf("Refresh: %v", err)
				ctx.SetStatusCode(fasthttp.StatusForbidden)
			} else {
				ctx.Response.Header.Set(strHeaderSession, newSession)
			}
		}
	} else {
		ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
	}
}

func (has *HttpApiHandler) HandlerUser(ctx *fasthttp.RequestCtx) {
	SetServerNameHeader(ctx, has.cfg.Server.ServerName)
	session := string(ctx.Request.Header.Peek(strHeaderSession))
	target := ctx.UserValue("login").(string)
	if len(session) == 0 {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
	} else if len(target) == 0 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
	} else {
		method := string(ctx.Method())
		switch method {
		case fasthttp.MethodGet:
			{
				switch target {
				case TargetSelf:
					{
						user, err := has.uc.User().GetInfoSelf(session)
						if err != nil {
							has.logx.Warnf("Get self is failed: %v", err)
							ctx.SetStatusCode(fasthttp.StatusForbidden)
						} else {
							SetBodyJson(ctx, &user)
						}
					}
				default:
					user, err := has.uc.User().GetUserInfo(session, target)
					if err != nil {
						has.logx.Warnf("Get user is failed: %v", err)
						ctx.SetStatusCode(fasthttp.StatusNotFound)
					} else {
						SetBodyJson(ctx, &user)
					}
				}
			}
		case fasthttp.MethodPatch:
			{
				switch target {
				case TargetSelf:
					user := db.User{}
					if err := json.Unmarshal(ctx.Request.Body(), &user); err != nil {
						ctx.SetStatusCode(fasthttp.StatusBadRequest)
					} else {
						if err := has.uc.User().ChangeInfoSelf(session, &user); err != nil {
							has.logx.Warnf("Change info user is failed: %v", err)
							ctx.SetStatusCode(fasthttp.StatusBadRequest)
						} else {
							ctx.SetStatusCode(fasthttp.StatusOK)
						}
					}
				default:
					ctx.SetStatusCode(fasthttp.StatusTeapot)
				}
			}
		case fasthttp.MethodDelete:
			{
				switch target {
				case TargetSelf:
					user := db.User{}
					if err := json.Unmarshal(ctx.Request.Body(), &user); err != nil {
						ctx.SetStatusCode(fasthttp.StatusBadRequest)
					} else {
						if err := has.uc.User().DeleteSelf(session, &user); err != nil {
							has.logx.Warnf("Delete user is failed: %v", err)
							ctx.SetStatusCode(fasthttp.StatusBadRequest)
						} else {
							ctx.SetStatusCode(fasthttp.StatusNoContent)
						}
					}
				default:
					ctx.SetStatusCode(fasthttp.StatusTeapot)
				}
			}
		default:
			ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
		}
	}
}

func (has *HttpApiHandler) HandlerChat(ctx *fasthttp.RequestCtx) {
	SetServerNameHeader(ctx, has.cfg.Server.ServerName)
	session := string(ctx.Request.Header.Peek(strHeaderSession))
	target := ctx.UserValue("desc").(string)
	offsetStr := string(ctx.QueryArgs().Peek("offset"))
	limitStr := string(ctx.QueryArgs().Peek("limit"))
	if len(session) == 0 {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
	} else {
		method := string(ctx.Method())
		switch method {
		case fasthttp.MethodGet:
			{
				switch target {
				case TargetSelf:
					{
						offset, _ := strconv.Atoi(offsetStr)
						limit, _ := strconv.Atoi(limitStr)
						chats, err := has.uc.Chat().GetChatsInfoByUser(session, offset, limit)
						if err != nil {
							has.logx.Warnf("Get self chats: %v", err)
							ctx.SetStatusCode(fasthttp.StatusBadRequest)
						} else {
							SetBodyJson(ctx, &chats)
						}
					}
				default:
					{
						chat, err := has.uc.Chat().GetChatInfo(session, target)
						if err != nil {
							has.logx.Warnf("Get info chat: %v", err)
							ctx.SetStatusCode(fasthttp.StatusBadRequest)
						} else {
							SetBodyJson(ctx, &chat)
						}
					}
				}
			}
		case fasthttp.MethodPost:
			{
				if err := has.uc.Chat().CreateChat(session, target); err != nil {
					has.logx.Warnf("Create chat: %v", err)
					ctx.SetStatusCode(fasthttp.StatusBadRequest)
				} else {
					SetLocationHeader(ctx, has.cfg.Server.PrefixPath+"/chat/"+target)
					ctx.SetStatusCode(fasthttp.StatusCreated)
				}
			}
		case fasthttp.MethodDelete:
			{
				if err := has.uc.Chat().DeleteChat(session, target); err != nil {
					has.logx.Warnf("Delete chat: %v", err)
					ctx.SetStatusCode(fasthttp.StatusBadRequest)
				} else {
					ctx.SetStatusCode(fasthttp.StatusNoContent)
				}
			}
		default:
			ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
		}
	}
}

func (has *HttpApiHandler) HandlerFollow(ctx *fasthttp.RequestCtx) {
	SetServerNameHeader(ctx, has.cfg.Server.ServerName)
	session := string(ctx.Request.Header.Peek(strHeaderSession))
	target := ctx.UserValue("desc").(string)
	offsetStr := string(ctx.QueryArgs().Peek("offset"))
	limitStr := string(ctx.QueryArgs().Peek("limit"))
	if len(session) == 0 {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
	} else {
		method := string(ctx.Method())
		switch method {
		case fasthttp.MethodGet:
			{
				offset, _ := strconv.Atoi(offsetStr)
				limit, _ := strconv.Atoi(limitStr)
				users, err := has.uc.Chat().GetUsersInChat(session, target, offset, limit)
				if err != nil {
					has.logx.Warnf("Get users from chat:%v err:%v", target, err)
					ctx.SetStatusCode(fasthttp.StatusBadRequest)
				} else {
					SetBodyJson(ctx, &users)
				}
			}
		case fasthttp.MethodPost:
			{
				if err := has.uc.AddUserToChat(session, target); err != nil {
					has.logx.Warnf("Add user to chat: %v", err)
					ctx.SetStatusCode(fasthttp.StatusBadRequest)
				} else {
					ctx.SetStatusCode(fasthttp.StatusOK)
				}
			}
		case fasthttp.MethodDelete:
			{
				if err := has.uc.DelUserFromChat(session, target); err != nil {
					ctx.SetStatusCode(fasthttp.StatusBadRequest)
				} else {
					ctx.SetStatusCode(fasthttp.StatusNoContent)
				}
			}
		default:
			{
				ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
			}
		}
	}
}

func (has *HttpApiHandler) HandlerMessage(ctx *fasthttp.RequestCtx) {
	SetServerNameHeader(ctx, has.cfg.Server.ServerName)
	session := string(ctx.Request.Header.Peek(strHeaderSession))
	if len(session) == 0 {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
	} else {
		method := string(ctx.Method())
		switch method {
		case fasthttp.MethodGet:
			{
				target := ctx.UserValue("chat").(string)
				offsetStr := string(ctx.QueryArgs().Peek("offset"))
				limitStr := string(ctx.QueryArgs().Peek("limit"))
				offset, _ := strconv.Atoi(offsetStr)
				limit, _ := strconv.Atoi(limitStr)
				messages, err := has.uc.Message().GetMessagesFromChat(session, target, offset, limit)
				if err != nil {
					has.logx.Warnf("Get messages from chat:%v err:%v", target, err)
					ctx.SetStatusCode(fasthttp.StatusBadRequest)
				} else {
					SetBodyJson(ctx, &messages)
				}
			}
		case fasthttp.MethodPost:
			{
				reqBody := ReMessage{}
				if err := json.Unmarshal(ctx.Request.Body(), &reqBody); err != nil {
					ctx.SetStatusCode(fasthttp.StatusBadRequest)
				} else {
					if err := has.uc.Message().PushMessageToChat(session, reqBody.Chat, reqBody.Text); err != nil {
						has.logx.Warnf("Push message: %v", err)
						ctx.SetStatusCode(fasthttp.StatusBadRequest)
					} else {
						ctx.SetStatusCode(fasthttp.StatusOK)
					}
				}
			}
		case fasthttp.MethodDelete:
			{
				reqBody := ReMessage{}
				if err := json.Unmarshal(ctx.Request.Body(), &reqBody); err != nil {
					ctx.SetStatusCode(fasthttp.StatusBadRequest)
				} else {
					if err := has.uc.Message().DelMessageFromChat(session, reqBody.Chat, reqBody.Message); err != nil {
						has.logx.Warnf("Delete message: %v", err)
						ctx.SetStatusCode(fasthttp.StatusBadRequest)
					} else {
						ctx.SetStatusCode(fasthttp.StatusNoContent)
					}
				}
			}
		default:
			{
				ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
			}
		}
	}
}

func (has *HttpApiHandler) HandlerAdmin(ctx *fasthttp.RequestCtx) {
	SetServerNameHeader(ctx, has.cfg.Server.ServerName)
	session := string(ctx.Request.Header.Peek(strHeaderSession))
	target := ctx.UserValue("entity").(string)
	role := string(ctx.QueryArgs().Peek("do"))
	offsetStr := string(ctx.QueryArgs().Peek("offset"))
	limitStr := string(ctx.QueryArgs().Peek("limit"))
	if len(session) == 0 {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
	} else {
		method := string(ctx.Method())
		switch method {
		case fasthttp.MethodGet:
			{
				switch target {
				case "users":
					{
						offset, _ := strconv.Atoi(offsetStr)
						limit, _ := strconv.Atoi(limitStr)
						users, err := has.uc.Role().GetAllUsers(session, offset, limit)
						if err != nil {
							has.logx.Warnf("Delete message: %v", err)
							ctx.SetStatusCode(fasthttp.StatusBadRequest)
						} else {
							SetBodyJson(ctx, &users)
						}
					}
				case "chats":
					{
						offset, _ := strconv.Atoi(offsetStr)
						limit, _ := strconv.Atoi(limitStr)
						chats, err := has.uc.Role().GetAllChats(session, offset, limit)
						if err != nil {
							has.logx.Warnf("Delete message: %v", err)
							ctx.SetStatusCode(fasthttp.StatusBadRequest)
						} else {
							SetBodyJson(ctx, &chats)
						}
					}
				default:
					ctx.SetStatusCode(fasthttp.StatusTeapot)
				}
			}
		case fasthttp.MethodPatch:
			{
				if err := has.uc.Role().SetRoleUser(session, target, role); err != nil {
					has.logx.Warnf("Set role: %v", err)
					ctx.SetStatusCode(fasthttp.StatusBadRequest)
				} else {
					ctx.SetStatusCode(fasthttp.StatusOK)
				}
			}
		}
	}
}

func (has *HttpApiHandler) HandlerAdminCreate(ctx *fasthttp.RequestCtx) {
	SetServerNameHeader(ctx, has.cfg.Server.ServerName)
	session := string(ctx.Request.Header.Peek(strHeaderSession))
	target := ctx.UserValue("entity").(string)
	if len(session) == 0 {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
	} else {
		method := string(ctx.Method())
		switch method {
		case fasthttp.MethodPost:
			{
				switch target {
				case "user":
					{
						user := db.User{}
						if err := json.Unmarshal(ctx.Request.Body(), &user); err != nil {
							has.logx.Warnf("Docode chat %v", err)
							ctx.SetStatusCode(fasthttp.StatusBadRequest)
						}
						if err := has.uc.Register(&user); err != nil {
							has.logx.Warnf("Create chat: %v", err)
							ctx.SetStatusCode(fasthttp.StatusBadRequest)
						} else {
							ctx.SetStatusCode(fasthttp.StatusOK)
						}
					}
				case "chat":
					{
						chat := db.Chat{}
						if err := json.Unmarshal(ctx.Request.Body(), &chat); err != nil {
							has.logx.Warnf("Docode chat %v", err)
							ctx.SetStatusCode(fasthttp.StatusBadRequest)
						}
						if err := has.uc.Chat().CreateChatAdmin(session, &chat); err != nil {
							has.logx.Warnf("Create chat: %v", err)
							ctx.SetStatusCode(fasthttp.StatusBadRequest)
						} else {
							ctx.SetStatusCode(fasthttp.StatusOK)
						}
					}
				default:
					ctx.SetStatusCode(fasthttp.StatusTeapot)
				}
			}
		default:
			ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
		}
	}
}

func SetServerNameHeader(ctx *fasthttp.RequestCtx, servername string) {
	ctx.Response.Header.Set("Server", servername)
}

func SetBodyJson(ctx *fasthttp.RequestCtx, some interface{}) {
	SetChacheControlDisable(ctx)
	ctx.Response.Header.Set(fasthttp.HeaderContentType, strApplicationJSON)
	bodyResp, _ := json.Marshal(some)
	ctx.SetBody(bodyResp)
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func SetChacheControlDisable(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set(fasthttp.HeaderCacheControl, "no-cache, no-store, must-revalidate")
	ctx.Response.Header.Set(fasthttp.HeaderPragma, "no-cache")
	ctx.Response.Header.Set(fasthttp.HeaderExpires, "0")
}

func SetLocationHeader(ctx *fasthttp.RequestCtx, url string) {
	ctx.Request.Header.Set(fasthttp.HeaderLocation, url)
}
