package usecase

import (
	"api_chat/api/layers/base/logx"
	"api_chat/api/layers/domain/cfg"
	"api_chat/api/layers/domain/db"
	mwauth "api_chat/api/layers/middleware/mw_auth"
	mwrole "api_chat/api/layers/middleware/mw_role"
	"api_chat/api/layers/repos"
	ucchat "api_chat/api/layers/usecase/uc_chat"
	ucmsg "api_chat/api/layers/usecase/uc_msg"
	ucuser "api_chat/api/layers/usecase/uc_user"
	"errors"
)

type UsecaseOperator struct {
	cfg    *cfg.Configuration
	role   *mwrole.RoleCtx
	log    logx.Logger
	ucuser *ucuser.UsecaseUser
	ucchat *ucchat.UsecaseChat
	ucmsg  *ucmsg.UsecaseMessage
	repos  *repos.ReposContext
}

func NewUsecaseOparetor(cfg *cfg.Configuration, reposCtx *repos.ReposContext, logx logx.Logger) *UsecaseOperator {
	role := mwrole.NewRoleContext(cfg, reposCtx, logx)
	return &UsecaseOperator{
		cfg:    cfg,
		repos:  reposCtx,
		log:    logx,
		ucuser: ucuser.NewUsecaseUser(reposCtx, role, mwauth.NewAuthCtx(cfg, reposCtx, logx), logx),
		ucchat: ucchat.NewUsecaseChat(reposCtx, role, logx),
		ucmsg:  ucmsg.NewUsecaseMessage(reposCtx, role, logx),
		role:   role,
	}
}

func (uc *UsecaseOperator) User() *ucuser.UsecaseUser {
	return uc.ucuser
}

func (uc *UsecaseOperator) Chat() *ucchat.UsecaseChat {
	return uc.ucchat
}

func (uc *UsecaseOperator) Message() *ucmsg.UsecaseMessage {
	return uc.ucmsg
}

func (uc *UsecaseOperator) Role() *mwrole.RoleCtx {
	return uc.role
}

// func (uc *UsecaseOperator) DestroySession(session string) error {
// 	return uc.repos.Session().SessionDestroy(session)
// }

func (uc *UsecaseOperator) AddUserToChat(session string, desc string) error {
	uc.log.Debugf("AddUserToChat called session:%v chat:%v", session, desc)
	userId, err := uc.repos.Session().SessionGetUserId(session)
	if err != nil {
		return err
	}

	chat := db.Chat{}
	if _, err := uc.repos.Chat().Get(&chat, map[string]interface{}{"desc": desc}, 0, 0); err != nil {
		return err
	}

	if uc.repos.IsExistsUserInChat(userId, chat.Id) {
		return errors.New("user is already added")
	}

	_, err = uc.repos.AddUserToChat(userId, chat.Id)
	return err
}

func (uc *UsecaseOperator) DelUserFromChat(session string, desc string) error {
	uc.log.Debugf("DelUserFromChat called session:%v chat:%v", session, desc)
	userId, err := uc.repos.Session().SessionGetUserId(session)
	if err != nil {
		return err
	}

	chat := db.Chat{Desc: desc}
	if _, err = uc.repos.Chat().Get(&chat, map[string]interface{}{"desc": desc}, 0, 0); err != nil {
		return err
	}

	_, err = uc.repos.DelUserFromChat(userId, chat.Id)
	return err
}
