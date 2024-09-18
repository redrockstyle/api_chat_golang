package repos

import (
	"api_chat/api/layers/base/hasher"
	"api_chat/api/layers/base/logx"
	"api_chat/api/layers/controller/database"
	"api_chat/api/layers/domain/db"
	"api_chat/api/layers/repos/entity/chatctx"
	"api_chat/api/layers/repos/entity/sessionctx"
	"api_chat/api/layers/repos/entity/userctx"
	"time"
)

type ReposContext struct {
	dbc  *database.DbController
	logx logx.Logger
	uCtx *userctx.UserContext
	cCtx *chatctx.ChatContext
	sCtx *sessionctx.SessionCtx
}

func NewReposContext(
	dbc *database.DbController,
	hasher *hasher.Hasher,
	logx logx.Logger,
	lifetime time.Duration,
) *ReposContext {
	return &ReposContext{
		dbc:  dbc,
		logx: logx,
		uCtx: userctx.NewUserContext(dbc, hasher, logx),
		cCtx: chatctx.NewChatContext(dbc, logx),
		sCtx: sessionctx.NewSessionCtx(dbc, hasher, lifetime),
	}
}

func (mc *ReposContext) User() *userctx.UserContext {
	return mc.uCtx
}

func (mc *ReposContext) Chat() *chatctx.ChatContext {
	return mc.cCtx
}

func (mc *ReposContext) Session() *sessionctx.SessionCtx {
	return mc.sCtx
}

func (mc *ReposContext) GetUsersByChat(userId string, chatId uint64, offset int, limit int) (*[]db.User, error) {
	cu := []db.ChatUser{}
	writed, err := mc.dbc.TypeGets("", &cu, map[string]interface{}{"id_chat": chatId}, offset, limit)
	if writed == 0 {
		return nil, err
	}

	users := make([]db.User, writed)
	for i := 0; i < int(writed); i++ {
		user := db.User{Id: cu[i].IdUser}
		if _, err := mc.dbc.TypeGet("", &user, nil); err != nil {
			return nil, err
		}
		user.Password = ""
		user.Id = ""
		user.Role = ""
		users[i] = user
	}
	return &users, nil
}

func (mc *ReposContext) GetChatsByUser(userId string, offset int, limit int) (*[]db.Chat, error) {
	cu := []db.ChatUser{}
	writed, err := mc.dbc.TypeGets("", &cu, map[string]interface{}{"id_user": userId}, offset, limit)
	if writed == 0 {
		return nil, err
	}

	sliceChat := make([]db.Chat, writed)
	for i := 0; i < int(writed); i++ {
		chat := db.Chat{Id: cu[i].IdChat}
		if _, err := mc.dbc.TypeGet("", &chat, nil); err != nil {
			return nil, err
		}

		user := db.User{Id: chat.Creator}
		if _, err := mc.dbc.TypeGet("", &user, nil); err != nil {
			return nil, err
		}
		chat.Creator = user.Login
		sliceChat[i] = chat
	}
	return &sliceChat, nil
}

func (mc *ReposContext) AddUserToChat(userId string, chatId uint64) (uint64, error) {
	cuid, err := mc.dbc.TypeCreate("", &db.ChatUser{IdUser: userId, IdChat: chatId})
	if err != nil {
		return 0, err
	}
	return cuid.(uint64), nil
}

func (mc *ReposContext) IsExistsUserInChat(userId string, chatId uint64) bool {
	return mc.dbc.TypeExists("", &db.ChatUser{}, map[string]interface{}{"id_user": userId, "id_chat": chatId})
}

func (mc *ReposContext) DelUserFromChat(userId string, chatId uint64) (uint64, error) {
	cuid, err := mc.dbc.TypeDelete("", &db.ChatUser{}, map[string]interface{}{"id_user": userId, "id_chat": chatId})
	if err != nil {
		return 0, err
	}
	return cuid.(uint64), nil
}

func (mc *ReposContext) AddMsgToChat(userId string, chatId uint64, msg string) (uint64, error) {
	tableName, err := mc.dbc.TypeTableIsExists(chatId, &db.Message{Id: chatId})
	if err != nil {
		if tableName, err = mc.dbc.TypeTableCreate(chatId, &db.Message{}); err != nil {
			return 0, err
		}
	}
	mid, err := mc.dbc.TypeCreate(tableName, &db.Message{IdUser: userId, Text: msg})
	if err != nil {
		return 0, err
	}
	return mid.(uint64), nil
}

func (mc *ReposContext) GetCountMessages(chatId uint64, conds map[string]interface{}) (int64, error) {
	return mc.dbc.TypeCount(chatId, &db.Message{}, conds)
}

func (mc *ReposContext) GetMessagesFromChat(chatId uint64, offset int, limit int) (*[]db.Message, error) {
	tableName, err := mc.dbc.TypeTableIsExists(chatId, &db.Message{})
	if err != nil {
		return nil, err
	}
	msgs := []db.Message{}
	if _, err := mc.dbc.TypeGets(tableName, &msgs, nil, offset, limit); err != nil {
		return nil, err
	}
	return &msgs, nil
}

func (mc *ReposContext) GetMessageFromChat(chatId uint64, msgId uint64) (*db.Message, error) {
	tableName, err := mc.dbc.TypeTableIsExists(chatId, &db.Message{})
	if err != nil {
		return nil, err
	}
	msg := db.Message{Id: msgId}
	if _, err := mc.dbc.TypeGet(tableName, &msg, nil); err != nil {
		return nil, err
	}
	return &msg, nil
}

func (mc *ReposContext) ModifyMsgByChat(msgId uint64, chatId uint64, msg string) (int64, error) {
	return 0, nil
}

func (mc *ReposContext) DelMsgFromChat(msgId uint64, chatId uint64) (uint64, error) {
	tableName, err := mc.dbc.TypeTableIsExists(chatId, &db.Message{})
	if err != nil {
		return 0, err
	}
	mid, err := mc.dbc.TypeDelete(tableName, &db.Message{Id: msgId}, nil)
	if err != nil {
		return 0, err
	}
	count, err := mc.GetCountMessages(chatId, nil)
	if err != nil {
		return 0, err
	}
	if count == 0 {
		if err := mc.dbc.TypeTableDropByName(tableName); err != nil {
			return 0, err
		}
	}
	return mid.(uint64), nil
}
