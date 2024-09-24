package ucchat

import (
	"api_chat/api/layers/base/logx"
	"api_chat/api/layers/domain/db"
	mwrole "api_chat/api/layers/middleware/mw_role"
	"api_chat/api/layers/repos"
	"errors"
)

type UsecaseChat struct {
	rep  *repos.ReposContext
	role *mwrole.RoleCtx
	log  logx.Logger
}

func NewUsecaseChat(rep *repos.ReposContext, role *mwrole.RoleCtx, logx logx.Logger) *UsecaseChat {
	return &UsecaseChat{rep: rep, role: role, log: logx}
}

func (ucu *UsecaseChat) CreateChatByAdmin(session string, chat *db.Chat) error {
	userId, err := ucu.rep.Session().SessionGetUserId(session)
	if err != nil {
		return err
	}

	if !ucu.role.IsAdminRoleById(userId) {
		return errors.New("operation is not permitted")
	}

	(*chat).Creator = userId
	_, err = ucu.rep.Chat().Create(chat)
	return err
}

func (ucu *UsecaseChat) CreateChat(session string, desc string) error {
	ucu.log.Debugf("CreateChat called session:%v desc:%v", session, desc)
	userId, err := ucu.rep.Session().SessionGetUserId(session)
	if err != nil {
		return err
	}
	chat := db.Chat{Desc: desc, Creator: userId}
	chatId, err := ucu.rep.Chat().Create(&chat)
	if err != nil {
		return err
	}
	_, err = ucu.rep.AddUserToChat(userId, chatId)
	return err
}

func (ucu *UsecaseChat) DeleteChat(session string, desc string) error {
	ucu.log.Debugf("DeleteChat called session:%v desc:%v", session, desc)
	userId, err := ucu.rep.Session().SessionGetUserId(session)
	if err != nil {
		return err
	}

	chat := db.Chat{Desc: desc}
	if _, err := ucu.rep.Chat().Get(&chat, map[string]interface{}{"desc": desc}, 0, 0); err != nil {
		return err
	}

	userSelf := db.User{}
	if _, err := ucu.rep.User().Get(&userSelf, map[string]interface{}{"id": userId}, 0, 0); err != nil {
		return err
	}

	if userSelf.Role != mwrole.RoleAdmin {
		if userId != chat.Creator && !ucu.role.IsPermDeleteByUser(&userSelf) {
			return errors.New("creator is not equal")
		}
	}
	_, err = ucu.rep.Chat().Delete(&chat)
	return err
}

func (ucu *UsecaseChat) GetChatInfo(session string, desc string) (*db.Chat, error) {
	ucu.log.Debugf("GetChatInfo called session:%v desc:%v", session, desc)
	if !ucu.rep.Session().SessionValidate(session) {
		return nil, errors.New("bad session")
	}
	// userId, err := ucu.rep.Session().SessionGetUserId(session)
	// if err != nil {
	// 	return nil, err
	// }

	chat := db.Chat{}
	if _, err := ucu.rep.Chat().Get(&chat, map[string]interface{}{"desc": desc}, 0, 0); err != nil {
		return nil, err
	}
	user := db.User{}
	if _, err := ucu.rep.User().Get(&user, map[string]interface{}{"id": chat.Creator}, 0, 0); err != nil {
		return nil, err
	}
	chat.Creator = user.Login
	return &chat, nil
}

func (ucu *UsecaseChat) GetChatsInfoByUser(session string, offset int, limit int) (*[]db.Chat, error) {
	ucu.log.Debugf("GetChatsInfoByUser called session:%v offset:%v limit:%v", session, offset, limit)
	userId, err := ucu.rep.Session().SessionGetUserId(session)
	if err != nil {
		return nil, err
	}
	return ucu.rep.GetChatsByUser(userId, offset, limit)
}

func (ucu *UsecaseChat) GetUsersInChat(session string, desc string, offset int, limit int) (*[]db.User, error) {
	ucu.log.Debugf("GetUsersInChat called session:%v offset:%v limit:%v", session, offset, limit)
	userId, err := ucu.rep.Session().SessionGetUserId(session)
	if err != nil {
		return nil, err
	}

	chat := db.Chat{}
	if _, err := ucu.rep.Chat().Get(&chat, map[string]interface{}{"desc": desc}, 0, 0); err != nil {
		return nil, err
	}

	if !ucu.role.IsAdminRoleById(userId) {
		if !ucu.rep.IsExistsUserInChat(userId, chat.Id) {
			return nil, errors.New("user is not follower")
		}
	}

	return ucu.rep.GetUsersByChat(userId, chat.Id, offset, limit)
}
