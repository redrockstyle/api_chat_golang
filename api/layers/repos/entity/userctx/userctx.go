package userctx

import (
	"api_chat/api/layers/base/hasher"
	"api_chat/api/layers/base/ident"
	"api_chat/api/layers/base/logx"
	"api_chat/api/layers/controller/database"
	"api_chat/api/layers/domain/db"
	"errors"
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
)

const (
	kekid = "come as you are"
)

type UserContext struct {
	//cfg  *cfg.Configuration
	dbc    *database.DbController
	logx   logx.Logger
	hasher *hasher.Hasher
	ident  *ident.Ident
	vld    *validator.Validate
}

func NewUserContext(dbc *database.DbController, hasher *hasher.Hasher, logx logx.Logger) *UserContext {
	return &UserContext{
		dbc:    dbc,
		logx:   logx,
		hasher: hasher,
		ident:  ident.NewIdent(),
		vld:    validator.New(validator.WithRequiredStructEnabled()),
	}
}

/*
 * Create user
 * Input: user struct (id is not required)
 * Proc: generate UUID and write to user->id, hash password and write to user->password
 * Return: id created user or error if user creation is failed
 */
func (ur *UserContext) Create(user *db.User) (string, error) {
	ur.logx.Debug("Create called")
	if err := ur.vld.Struct(user); err != nil {
		return "", err
	}
	if !ur.IsUniqueLogin(user.Login) {
		return "", errors.New("this login is already used")
	}

	var err error
	user.Id = ur.ident.GenerateUUID(user.Login + kekid)
	user.Password, err = ur.hasher.HashStr(user.Password)
	if err != nil {
		return "", err
	}
	retId, err := ur.dbc.TypeCreate("", user)
	if err != nil {
		return "", err
	}
	if user.Id != retId.(string) {
		return "", fmt.Errorf("error create: returned '%v' != '%v'", retId.(string), user.Id)
	}
	ur.logx.Infof("Created new user login:%v", user.Login)
	return user.Id, err
}

func (ur *UserContext) IsUniqueLogin(login string) bool {
	if login == "" {
		return false
	}

	if ur.dbc.TypeExists("", &db.User{}, map[string]interface{}{"login": login}) {
		return false
	}
	return true
}

/*
 * Delete user
 * Input: id_user or user struct (with required id)
 * Return: id deleted user or error if delete is failed
 */
func (ur *UserContext) Delete(t interface{}) (string, error) {
	ur.logx.Debug("Delete called")
	val := reflect.ValueOf(t)
	var id interface{}
	var err error
	ur.logx.Debugf("called delete user with type %v", val)
	if val.Kind() == reflect.String {
		str := t.(string)
		if id, err = ur.dbc.TypeDelete("", &db.User{Id: str}, nil); err != nil {
			return "", err
		}
	} else if val.Kind() == reflect.Ptr && val.Elem().Kind() == reflect.Struct {
		if err := ur.vld.Struct(t); err != nil {
			return "", err
		}
		if !ur.dbc.TypeExists("", t, nil) {
			return "", errors.New("record is not found")
		}
		if id, err = ur.dbc.TypeDelete("", t, nil); err != nil {
			return "", err
		}
	} else {
		ur.logx.Warnf("call delete user with unsupport %v", t)
		return "", errors.New("unsupport argument")
	}
	ur.logx.Infof("User id:%v deleted", id.(string))
	return id.(string), nil
}

/*
 * Modify user
 * Input: struct user with required user->id
 * Proc: fields will be modified if get(user->id).password == user->password
 * Return: id modified user or error if modification user fields is failed
 */
func (ur *UserContext) Modify(user *db.User) (string, error) {
	ur.logx.Debug("Modify called")
	if err := ur.vld.Struct(user); err != nil {
		return "", err
	}
	getUser := db.User{Id: user.Id}
	if _, err := ur.dbc.TypeGet("", &getUser, nil); err != nil {
		return "", err
	}
	if !ur.hasher.CompareHashAndStr(getUser.Password, user.Password) {
		return "", errors.New("pass is not equal")
	}
	user.Password = getUser.Password
	id, err := ur.dbc.TypeModify("", user)
	if err != nil {
		return "", err
	}
	return id.(string), nil
}

func (ur *UserContext) ModifyPass(userId string, oldPass string, newPass string) error {
	ur.logx.Debug("ModifyPass called")
	if userId == "" || newPass == "" {
		return errors.New("passed nil")
	}
	getUser := db.User{Id: userId}
	if _, err := ur.dbc.TypeGet("", &getUser, nil); err != nil {
		return err
	}
	if !ur.hasher.CompareHashAndStr(getUser.Password, oldPass) {
		return errors.New("pass is not equal")
	}

	pass, err := ur.hasher.HashStr(newPass)
	if err != nil {
		return err
	}

	if _, err := ur.dbc.TypeModifyCol("", &db.User{}, map[string]interface{}{"password": pass}); err != nil {
		return err
	}
	return nil
}

func (ur *UserContext) ModifyCol(userId string, conds map[string]interface{}) (string, error) {
	// if !ur.ident.CheckUUIDv5(userId) {
	// 	return "", errors.New("invalid user id")
	// }
	id, err := ur.dbc.TypeModifyCol("", &db.User{Id: userId}, conds)
	if err != nil {
		return "", err
	}
	return id.(string), nil
}

/*
 * Get user (or gets users)
 * Input1: user->id (string): returned struct user
 * Input2: struct user: returned id user and write to a passed struct user
 * Input3: slice struct users (used conds and idCnt): returned len users and write to a passed slice struct user
 * Proc1: conds is {"key":"value"} example
 * Proc2: idCnt is []string{"id1","id2","id3"} example index
 */
func (ur *UserContext) Get(t interface{}, conds map[string]interface{}, offset int, limit int) (interface{}, error) {
	ur.logx.Debug("Get called")
	val := reflect.ValueOf(t)
	var ids interface{}
	var err error
	ur.logx.Debugf("called get user with type %v", val)
	if t == nil {
		return nil, errors.New("arg one is most be nil")
	}

	if val.Kind() == reflect.String {
		user := db.User{
			Id: t.(string),
		}
		if _, err = ur.dbc.TypeGet("", &user, conds); err != nil {
			return nil, err
		}
		return &user, nil
	} else if val.Kind() == reflect.Ptr && val.Elem().Kind() == reflect.Struct {
		if ids, err = ur.dbc.TypeGet("", t, conds); err != nil {
			return nil, err
		}
	} else if val.Kind() == reflect.Ptr && val.Elem().Kind() == reflect.Slice {
		if ids, err = ur.dbc.TypeGets("", t, conds, offset, limit); err != nil {
			return nil, err
		}
	} else {
		ur.logx.Warnf("call get user with unsupport %v", t)
		return nil, errors.New("unupport argument")
	}
	return ids, nil
}

// func (ur *UserContext) GetIdUser(login string, pass string) (string, error) {
// 	user := db.User{}
// 	if _, err := ur.Get(&user, map[string]interface{}{"login": login}, 0, 0); err != nil {
// 		return "", err
// 	}
// 	if ur.ComparePass(user.Id, pass) {
// 		return user.Id, nil
// 	} else {
// 		return "", errors.New("pass is not equal")
// 	}
// }

// func (ur *UserContext) Count(conds map[string]interface{}) (int64, error) {
// 	return ur.dbc.TypeCount(0, &db.User{}, conds)
// }

/*
 * Compare pass
 * Input1: user->id (string) and pass
 * Input2: struct user and pass
 * Proc: get user, hashed pass and compare
 * Return: true if pass==user->pass or false if pass!=user->pass
 */
func (ur *UserContext) ComparePass(user interface{}, pass string) bool {
	ur.logx.Debug("ComparePass called")
	if pass == "" {
		return false
	}

	var userStruct *db.User
	if user != nil {
		refUser := reflect.ValueOf(user)
		if refUser.Kind() == reflect.String {
			userStr := user.(string)
			if userStr != "" {
				use := db.User{
					Id: userStr,
				}
				if _, err := ur.dbc.TypeGet("", &use, nil); err != nil {
					return false
				}
				userStruct = &use
			}
		} else if refUser.Kind() == reflect.Ptr && refUser.Elem().Kind() == reflect.Struct {
			if _, err := ur.dbc.TypeGet("", user, nil); err != nil {
				return false
			}
			userStruct = user.(*db.User)
		}
	}
	if userStruct != nil {
		if ur.hasher.CompareHashAndStr(userStruct.Password, pass) {
			return true
		} else {
			ur.logx.Warnf("input wrong password for %v", userStruct.Login)
			return false
		}
	} else {
		return false
	}
}

func (ur *UserContext) ComparePassAndLogin(user interface{}, login string, pass string) bool {
	ur.logx.Debug("ComparePassAndLogin called")
	if pass == "" || login == "" {
		return false
	}

	var userStruct *db.User
	if user != nil {
		refUser := reflect.ValueOf(user)
		if refUser.Kind() == reflect.String {
			userStr := user.(string)
			if userStr != "" {
				use := db.User{
					Id: userStr,
				}
				if _, err := ur.dbc.TypeGet("", &use, nil); err != nil {
					return false
				}
				userStruct = &use
			}
		} else if refUser.Kind() == reflect.Ptr && refUser.Elem().Kind() == reflect.Struct {
			if _, err := ur.dbc.TypeGet("", user, nil); err != nil {
				return false
			}
			userStruct = user.(*db.User)
		}
	}
	if userStruct != nil {
		if ur.hasher.CompareHashAndStr(userStruct.Password, pass) && userStruct.Login == login {
			return true
		} else {
			ur.logx.Warnf("input wrong compare for %v", userStruct.Login)
			return false
		}
	} else {
		return false
	}
}
