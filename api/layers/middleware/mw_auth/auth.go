package mw_auth

import (
	"api_chat/api/layers/base/logx"
	"api_chat/api/layers/domain/cfg"
	"api_chat/api/layers/domain/db"
	mwrole "api_chat/api/layers/middleware/mw_role"
	"api_chat/api/layers/repos"
	"errors"

	"github.com/go-playground/validator/v10"
)

type AuthCtx struct {
	cfg   *cfg.Configuration
	repos *repos.ReposContext
	logx  logx.Logger
	vld   *validator.Validate
}

func NewAuthCtx(cfg *cfg.Configuration, reposCtx *repos.ReposContext, logx logx.Logger) *AuthCtx {
	return &AuthCtx{cfg: cfg, repos: reposCtx, logx: logx, vld: validator.New(validator.WithRequiredStructEnabled())}
}

func (ac *AuthCtx) Register(user *db.User) error {
	ac.logx.Debug("Register called")
	user.Role = mwrole.RoleCommon
	if err := ac.IsIgnoreUsernames(user.Login); err != nil {
		return err
	}
	_, err := ac.repos.User().Create(user)
	return err
}

/*
 * Auth user
 * Return: session (string) if auth is success or error if auth is failed
 */
func (ac *AuthCtx) Authentificate(user *db.User) (string, error) {
	ac.logx.Debug("Auth called")
	if err := ac.vld.Struct(user); err != nil {
		return "", err
	}

	getUser := db.User{Login: user.Login}
	if _, err := ac.repos.User().Get(&getUser, map[string]interface{}{"login": user.Login}, 0, 0); err != nil {
		return "", err
	}
	if !ac.repos.User().ComparePass(getUser.Id, user.Password) {
		return "", errors.New("pass is not equal")
	}

	if err := ac.repos.Session().SesstionCloseById(getUser.Id); err != nil {
		return "", err
	}

	return ac.repos.Session().SessionCreate(getUser.Id)
}

// func (ac *AuthCtx) CheckSession(session string) bool {
// 	return ac.repos.Session().SessionValidate(session)
// }

func (ac *AuthCtx) Refresh(session string) (string, error) {
	ac.logx.Debug("Refresh called")
	userId, err := ac.repos.Session().SessionGetUserId(session)
	if err != nil {
		return "", err
	}

	if err := ac.repos.Session().SessionDestroy(session); err != nil {
		return "", err
	}

	newSession, err := ac.repos.Session().SessionCreate(userId)
	if err != nil {
		return "", err
	}

	return newSession, nil
}

func (ac *AuthCtx) Logout(session string) error {
	ac.logx.Debug("Logout called")
	return ac.repos.Session().SessionDestroy(session)
}

func (ac *AuthCtx) IsIgnoreUsernames(login string) error {
	err := errors.New("login is protected")
	switch login {
	case "self":
		return err
	case "target":
		return err
	case "admin":
		return err
	default:
		return nil
	}
}
