package ucmsg

import (
	"api_chat/api/layers/base/logx"
	"api_chat/api/layers/domain/db"
	mwrole "api_chat/api/layers/middleware/mw_role"
	"api_chat/api/layers/repos"
	"errors"
)

type UsecaseMessage struct {
	rep  *repos.ReposContext
	role *mwrole.RoleCtx
	log  logx.Logger
}

func NewUsecaseMessage(rep *repos.ReposContext, role *mwrole.RoleCtx, log logx.Logger) *UsecaseMessage {
	return &UsecaseMessage{rep: rep, role: role, log: log}
}

func (ucm *UsecaseMessage) GetMessagesFromChat(session string, chatDesc string, offset int, limit int) (*[]db.Message, error) {
	ucm.log.Debugf("GetMessagesFromChat called session:%v offset:%v limit:%v", session, offset, limit)
	userId, err := ucm.rep.Session().SessionGetUserId(session)
	if err != nil {
		return nil, errors.New("invalid session")
	}

	chat := db.Chat{}
	if _, err := ucm.rep.Chat().Get(&chat, map[string]interface{}{"desc": chatDesc}, 0, 0); err != nil {
		return nil, err
	}

	if !ucm.role.IsAdminRoleById(userId) {
		if !ucm.rep.IsExistsUserInChat(userId, chat.Id) {
			return nil, errors.New("this user is not follower")
		}
	}

	messages, err := ucm.rep.GetMessagesFromChat(chat.Id, offset, limit)
	if err != nil {
		return nil, err
	}

	lenMsg := len(*messages)
	for i := 0; i < lenMsg; i++ {
		user := db.User{Id: (*messages)[i].IdUser}
		if _, err := ucm.rep.User().Get(&user, nil, 0, 0); err != nil {
			return nil, err
		}
		(*messages)[i].IdUser = user.Login
	}
	return messages, nil
}

func (ucm *UsecaseMessage) PushMessageToChat(session string, chatDesc string, text string) error {
	lenText := len(text)
	ucm.log.Debugf("PushMessageToChat called session:%v chat:%v text:%v", session, chatDesc, lenText)
	userId, err := ucm.rep.Session().SessionGetUserId(session)
	if err != nil {
		return err
	}
	if lenText == 0 {
		return errors.New("empty message")
	}

	chat := db.Chat{}
	if _, err := ucm.rep.Chat().Get(&chat, map[string]interface{}{"desc": chatDesc}, 0, 0); err != nil {
		return err
	}

	if !ucm.role.IsAdminRoleById(userId) {
		if !ucm.rep.IsExistsUserInChat(userId, chat.Id) {
			return errors.New("this user is not follower")
		}
	}

	ucm.log.Infof("push message (%v-bytes) to %v", len(text), chatDesc)
	_, err = ucm.rep.AddMsgToChat(userId, chat.Id, text)
	return err
}

func (ucm *UsecaseMessage) DelMessageFromChat(session string, chatDesc string, msgId uint64) error {
	ucm.log.Debugf("DelMessageFromChat called session:%v chat:%v msgId:%v", session, chatDesc, msgId)
	userId, err := ucm.rep.Session().SessionGetUserId(session)
	if err != nil {
		return err
	}

	chat := db.Chat{}
	if _, err := ucm.rep.Chat().Get(&chat, map[string]interface{}{"desc": chatDesc}, 0, 0); err != nil {
		return err
	}

	if !ucm.role.IsAdminRoleById(userId) {
		if !ucm.role.IsPermDeleteByUserId(userId) {
			return errors.New("operation is not permitted")
		}
		if !ucm.rep.IsExistsUserInChat(userId, chat.Id) {
			return errors.New("this user is not follower")
		}
	}

	if _, err := ucm.rep.GetMessageFromChat(chat.Id, msgId); err != nil {
		return err
	}

	ucm.log.Infof("delete message id:%v from chat:%v", msgId, chatDesc)
	_, err = ucm.rep.DelMsgFromChat(msgId, chat.Id)
	return err
}
