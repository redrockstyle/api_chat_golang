package mwrole

import (
	"api_chat/api/layers/base/logx"
	"api_chat/api/layers/domain/cfg"
	"api_chat/api/layers/domain/db"
	"api_chat/api/layers/repos"
	"errors"
)

type RoleCtx struct {
	cfg *cfg.Configuration
	rep *repos.ReposContext
	log logx.Logger
}

const (
	RoleAdmin  = "admin"
	RoleMiddle = "middle"
	RoleCommon = "common"
)

func NewRoleContext(cfg *cfg.Configuration, rep *repos.ReposContext, logx logx.Logger) *RoleCtx {
	return &RoleCtx{cfg: cfg, rep: rep, log: logx}
}

func (rc *RoleCtx) GetAllUsers(session string, offset int, limit int) ([]db.User, error) {
	userId, err := rc.rep.Session().SessionGetUserId(session)
	if err != nil {
		return nil, err
	}
	if !rc.IsPermDeleteByUserId(userId) {
		return nil, errors.New("operation is not permitted")
	}

	users := []db.User{}
	if _, err := rc.rep.User().Get(&users, nil, offset, limit); err != nil {
		return nil, err
	}

	lenSlice := len(users)
	for i := 0; i < lenSlice; i++ {
		users[i].Password = ""
	}
	return users, nil
}

func (rc *RoleCtx) GetAllChats(session string, offset int, limit int) ([]db.Chat, error) {
	userId, err := rc.rep.Session().SessionGetUserId(session)
	if err != nil {
		return nil, err
	}
	if !rc.IsPermDeleteByUserId(userId) {
		return nil, errors.New("operation is not permitted")
	}

	chats := []db.Chat{}
	if _, err := rc.rep.Chat().Get(&chats, nil, offset, limit); err != nil {
		return nil, err
	}

	lenSlice := len(chats)
	for i := 0; i < lenSlice; i++ {
		user := db.User{}
		if _, err := rc.rep.User().Get(&user, map[string]interface{}{"id": chats[i].Creator}, 0, 0); err != nil {
			return nil, err
		}
		chats[i].Creator = user.Login
	}
	return chats, nil
}

func (rc *RoleCtx) SetRoleUser(session string, login string, role string) error {
	if !rc.IsCorrectRole(role) {
		return errors.New("incorrect role")
	}

	userId, err := rc.rep.Session().SessionGetUserId(session)
	if err != nil {
		return err
	}

	user := db.User{}
	if _, err := rc.rep.User().Get(&user, map[string]interface{}{"id": userId}, 0, 0); err != nil {
		return err
	}
	if !rc.IsAdminRole(&user) {
		return errors.New("user is not admin")
	}

	userTarget := db.User{}
	if _, err := rc.rep.User().Get(&userTarget, map[string]interface{}{"login": login}, 0, 0); err != nil {
		return err
	}

	if _, err := rc.rep.User().ModifyCol(userTarget.Id, map[string]interface{}{"role": role}); err != nil {
		return err
	}

	rc.log.Infof("Set role:%v for username:%v", role, userTarget.Login)
	return nil
}

func (rc *RoleCtx) IsAdminRoleById(userId string) bool {
	user := db.User{Id: userId}
	if _, err := rc.rep.User().Get(&user, nil, 0, 0); err != nil {
		return false
	}
	return rc.IsAdminRole(&user)
}

func (rc *RoleCtx) IsAdminRole(user *db.User) bool {
	return user.Role == RoleAdmin
}

func (rc *RoleCtx) IsPermDeleteByUserId(userId string) bool {
	user := db.User{Id: userId}
	if _, err := rc.rep.User().Get(&user, nil, 0, 0); err != nil {
		return false
	}

	return rc.IsPermDeleteByUser(&user)
}

func (rc *RoleCtx) IsPermDeleteByUser(user *db.User) bool {
	switch user.Role {
	case RoleAdmin:
		return true
	case RoleMiddle:
		return true
	}
	return false
}

func (rc *RoleCtx) IsCorrectRole(role string) bool {
	switch role {
	case RoleAdmin:
		return true
	case RoleMiddle:
		return true
	case RoleCommon:
		return true
	default:
		return false
	}
}
