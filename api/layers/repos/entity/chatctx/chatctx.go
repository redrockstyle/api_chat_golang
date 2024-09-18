package chatctx

import (
	//"api_chat/api/domain/cfg"
	"api_chat/api/layers/base/logx"
	"api_chat/api/layers/controller/database"
	"api_chat/api/layers/domain/db"
	"errors"
	"reflect"

	"github.com/go-playground/validator/v10"
)

type ChatContext struct {
	dbc  *database.DbController
	logx logx.Logger
	vld  *validator.Validate
}

func NewChatContext(dbc *database.DbController, logx logx.Logger) *ChatContext {
	return &ChatContext{
		dbc:  dbc,
		logx: logx,
		vld:  validator.New(validator.WithRequiredStructEnabled()),
	}
}

/*
 * Create chat
 * Input: string name or chat struct
 * Return: id created chat or error if creation is failed
 */
func (cr *ChatContext) Create(chat interface{}) (uint64, error) {
	cr.logx.Debug("Create called")
	var id interface{}
	var err error
	chatStruct := GetChatByInter(chat)
	if err := cr.vld.Struct(chatStruct); err != nil {
		return 0, err
	}
	if chatStruct != nil {
		if cr.dbc.TypeExists("", chatStruct, map[string]interface{}{"desc": chatStruct.Desc}) {
			return 0, errors.New("chat is already created")
		}
		if id, err = cr.dbc.TypeCreate("", chatStruct); err != nil {
			return 0, err
		}
	} else {
		return 0, errors.New("unsupport argument")
	}
	return id.(uint64), nil
}

/*
 * Delete chat
 * Input: string name or chat struct
 * Return: id deleted chat or error if deletion is failed
 */
func (cr *ChatContext) Delete(chat interface{}) (uint64, error) {
	cr.logx.Debug("Delete called")
	var id uint64
	chatStruct := GetChatByInter(chat)
	// if err := cr.vld.Struct(chatStruct); err != nil {
	// 	return 0, err
	// }
	if chatStruct != nil {
		if !cr.dbc.TypeExists("", chatStruct, nil) {
			return 0, errors.New("record is not found")
		}
		idType, err := cr.dbc.TypeDelete("", chatStruct, nil)
		if err != nil {
			return 0, err
		}
		id = idType.(uint64)
		if cr.dbc.TypeExists("", &db.ChatUser{}, map[string]interface{}{"id_chat": id}) {
			if _, err := cr.dbc.TypeDelete("", &db.ChatUser{}, map[string]interface{}{"id_chat": id}); err != nil {
				return 0, err
			}
		}
		if cr.dbc.TypeExists("", &db.ChatMessage{}, map[string]interface{}{"id_chat": id}) {
			if _, err := cr.dbc.TypeDelete("", &db.ChatMessage{}, map[string]interface{}{"id_chat": id}); err != nil {
				return 0, err
			}

			if tableName, err := cr.dbc.TypeTableIsExists(id, &chatStruct); err == nil {
				return 0, cr.dbc.TypeTableDropByName(tableName)

			}
			return id, nil
		}
	} else {
		return 0, errors.New("unsupport argument")
	}
	return id, nil
}

/*
 * Modify chat
 * Return: id modified chat or error if modification is failed
 */
func (cr *ChatContext) Modify(chat *db.Chat) (uint64, error) {
	cr.logx.Debug("Modify called")
	if err := cr.vld.Struct(chat); err != nil {
		return 0, err
	}
	if !cr.dbc.TypeExists("", chat, nil) {
		return 0, errors.New("record is not found")
	}
	id, err := cr.dbc.TypeModify("", chat)
	if err != nil {
		return 0, err
	}
	return id.(uint64), err
}

/*
 * Get chat (or get's chats)
 * Input1: chat->id (string): returned struct user
 * Input2: struct of chat: returned new struct of chat
 * Input3: slice struct of chats (used conds and idCnt): returned len chats and write to a passed slice struct of chats
 * Proc1: conds is {"key":"value"} example
 * Proc2: idCnt is []int{1,2,3} example index
 */
func (cr *ChatContext) Get(t interface{}, conds map[string]interface{}, offset int, limit int) (interface{}, error) {
	cr.logx.Debug("Get called")
	var ids interface{}
	var err error
	if t != nil {
		chatStruct := GetChatByInter(t)
		if chatStruct != nil {
			// if err := cr.vld.Struct(chatStruct); err != nil {
			// 	return 0, err
			// }
			if _, err = cr.dbc.TypeGet("", chatStruct, conds); err != nil {
				return nil, err
			}
			return *chatStruct, nil
		} else {
			val := reflect.ValueOf(t)
			if val.Kind() == reflect.Ptr && val.Elem().Kind() == reflect.Slice {
				if ids, err = cr.dbc.TypeGets("", t, conds, offset, limit); err != nil {
					return nil, err
				}
			} else {
				return nil, errors.New("unsupport argument")
			}
		}
	} else {
		return nil, errors.New("arg one most be nil")
	}
	return ids, nil
}

func (cr *ChatContext) Count(conds map[string]interface{}) (int64, error) {
	return cr.dbc.TypeCount(0, &db.Chat{}, conds)
}

/*
 * Get chat struct from interface
 * Return: chat struct if interface==typeOf(string) or interface==typeOf(*Chat)
 */
func GetChatByInter(chat interface{}) *db.Chat {
	var chatStruct *db.Chat
	if chat != nil {
		val := reflect.ValueOf(chat)
		if val.Kind() == reflect.Uint64 {
			chatStruct = &db.Chat{Id: chat.(uint64)}
		} else if val.Kind() == reflect.Ptr && val.Elem().Kind() == reflect.Struct {
			chatStruct = chat.(*db.Chat)
		}
	}
	return chatStruct
}
