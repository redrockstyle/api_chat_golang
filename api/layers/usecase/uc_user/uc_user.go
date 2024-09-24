package ucuser

import (
	"api_chat/api/layers/base/logx"
	"api_chat/api/layers/domain/db"
	mwauth "api_chat/api/layers/middleware/mw_auth"
	mwrole "api_chat/api/layers/middleware/mw_role"
	"api_chat/api/layers/repos"
	"errors"
)

type UsecaseUser struct {
	rep  *repos.ReposContext
	role *mwrole.RoleCtx
	auth *mwauth.AuthCtx
	log  logx.Logger
}

func NewUsecaseUser(rep *repos.ReposContext, role *mwrole.RoleCtx, auth *mwauth.AuthCtx, logx logx.Logger) *UsecaseUser {
	return &UsecaseUser{rep: rep, role: role, auth: auth, log: logx}
}

func (ucu *UsecaseUser) Register(user *db.User) error {
	return ucu.auth.Register(user)
}

func (ucu *UsecaseUser) Login(user *db.User) (string, error) {
	return ucu.auth.Authentificate(user)
}

func (ucu *UsecaseUser) Refresh(session string) (string, error) {
	return ucu.auth.Refresh(session)
}

func (ucu *UsecaseUser) Logout(session string) error {
	return ucu.auth.Logout(session)
}

func (ucu *UsecaseUser) CreateUserByAdmin(session string, user *db.User) error {
	userId, err := ucu.rep.Session().SessionGetUserId(session)
	if err != nil {
		return err
	}
	if !ucu.role.IsAdminRoleById(userId) {
		return errors.New("operation is not permitted")
	}

	return ucu.auth.Register(user)
}

func (ucu *UsecaseUser) ChangeInfoSelf(session string, user *db.User) error {
	ucu.log.Debug("ChangeInfoSelf called")
	userId, err := ucu.rep.Session().SessionGetUserId(session)
	if err != nil {
		return err
	}
	user.Id = userId
	if _, err := ucu.rep.User().Modify(user); err != nil {
		return err
	}
	return nil
}

func (ucu *UsecaseUser) ChangePasswordSelf(session string, oldPass string, newPass string) error {
	ucu.log.Debug("ChangePasswordSelf called")
	userId, err := ucu.rep.Session().SessionGetUserId(session)
	if err != nil {
		return err
	}
	if err := ucu.rep.User().ModifyPass(userId, oldPass, newPass); err != nil {
		return err
	}
	return nil
}

// mb realesed later
// func (ucu *UsecaseUser) GetInfoFriends()   {}

func (ucu *UsecaseUser) GetUserInfo(session string, loginTarget string) (*db.User, error) {
	ucu.log.Debug("GetUserInfo called")
	if !ucu.rep.Session().SessionValidate(session) {
		return nil, errors.New("invalid session")
	}

	userTarget := db.User{}
	if _, err := ucu.rep.User().Get(&userTarget, map[string]interface{}{"login": loginTarget}, 0, 0); err != nil {
		return nil, err
	}
	userTarget.Password = ""
	userTarget.Id = ""
	return &userTarget, nil
}

func (ucu *UsecaseUser) GetInfoSelf(session string) (*db.User, error) {
	ucu.log.Debug("GetInfoSelf called")
	userId, err := ucu.rep.Session().SessionGetUserId(session)
	if err != nil {
		return nil, err
	}
	userInt, err := ucu.rep.User().Get(userId, nil, 0, 0)
	if err != nil {
		return nil, err
	}

	user := userInt.(*db.User)
	user.Password = ""
	user.Id = ""
	return user, nil
}

func (ucu *UsecaseUser) DeleteSelf(session string, user *db.User) error {
	ucu.log.Debug("DeleteSelf called")
	userId, err := ucu.rep.Session().SessionGetUserId(session)
	if err != nil {
		return err
	}

	userSelf := db.User{}
	if _, err := ucu.rep.User().Get(&userSelf, map[string]interface{}{"id": userId}, 0, 0); err != nil {
		return err
	}

	userTarget := db.User{}
	if !ucu.role.IsAdminRole(&userSelf) {
		if !ucu.role.IsPermDeleteByUser(&userSelf) {
			return errors.New("operation is not permitted")
		} else {
			userTarget = userSelf
			if !ucu.rep.User().ComparePassAndLogin(userId, user.Login, user.Password) {
				return errors.New("pass/login is not equal")
			}
			if err := ucu.rep.Session().SessionDestroy(session); err != nil {
				return err
			}
		}
	} else {
		if _, err := ucu.rep.User().Get(&userTarget, map[string]interface{}{"login": user.Login}, 0, 0); err != nil {
			return err
		}
	}
	_, err = ucu.rep.User().Delete(userTarget.Id)
	return err
}
