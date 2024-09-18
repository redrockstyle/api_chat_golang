package repostestsuser

import (
	"api_chat/api/layers/domain/db"
	"api_chat/api/layers/repos"
	"fmt"
	"testing"
)

var userId string
var user2Id string

const (
	userFN    = "tom"
	userLN    = "lastnm"
	userLogin = "lgoin"
	UserPass  = "password1"
	UserPass2 = "password2"

	user2FN    = "tom"
	user2LN    = "sawyer"
	user2Login = "tomsaw"
	user2Pass  = "passtom"
)

type TestReposUserCaller struct{}

type TestsReposUser interface {
	TestsUserAll(repos *repos.ReposContext, t *testing.T)

	TestUserCreate(repos *repos.ReposContext, t *testing.T)
	TestUserGet(repos *repos.ReposContext, t *testing.T)
	TestUserModify(repos *repos.ReposContext, t *testing.T)
	TestUserDelete(repos *repos.ReposContext, t *testing.T)
	TestUserCmpPass(repos *repos.ReposContext, t *testing.T)
}

func NewTestReposUserCaller() *TestReposUserCaller {
	return &TestReposUserCaller{}
}

func (c *TestReposUserCaller) TestsUserAll(repos *repos.ReposContext, t *testing.T) {
	c.TestUserCreate(repos, t)
	c.TestUserGet(repos, t)
	c.TestUserGets(repos, t)
	// c.TestUserCount(repos, t)
	c.TestUserCmpPass(repos, UserPass, t)
	c.TestUserModify(repos, t)
	//c.TestUserCmpPass(repos, UserPass2, t)
	c.TestUserDelete(repos, t)
}

func (c *TestReposUserCaller) TestUserCreate(repos *repos.ReposContext, t *testing.T) {
	user := db.User{
		FirstName: userFN,
		LastName:  userLN,
		Login:     userLogin,
		Password:  UserPass,
	}
	user2 := db.User{
		FirstName: user2FN,
		LastName:  user2LN,
		Login:     user2Login,
		Password:  user2Pass,
	}
	var err error
	t.Run("test repos create user", func(t *testing.T) {
		if userId, err = repos.User().Create(&user); err != nil || userId == "" || len(userId) != 36 {
			t.Errorf("error create user login:%v retId:%v err:%v", user.Login, userId, err)
		}
		if user2Id, err = repos.User().Create(&user2); err != nil || user2Id == "" || len(user2Id) != 36 {
			t.Errorf("error create user login:%v retId:%v err:%v", user2.Login, user2Id, err)
		}
	})
}

func (c *TestReposUserCaller) TestUserGet(repos *repos.ReposContext, t *testing.T) {
	user := db.User{
		Id: userId,
	}
	t.Run("test repos get user", func(t *testing.T) {
		if retUser, err := repos.User().Get(userId, nil, 0, 0); err != nil ||
			retUser.(*db.User).Id != userId || retUser.(*db.User).FirstName != userFN ||
			retUser.(*db.User).LastName != userLN || retUser.(*db.User).Login != userLogin {
			t.Errorf("error get user (string) err:%v", err)
		}

		if id, err := repos.User().Get(&user, nil, 0, 0); err != nil || id.(string) != userId {
			t.Errorf("error get user (struct) userId:%v err:%v", userId, err)
		}
	})
}

func (c *TestReposUserCaller) TestUserGets(repos *repos.ReposContext, t *testing.T) {
	users := make([]db.User, 2)
	t.Run("test repos get's users", func(t *testing.T) {
		_, err := repos.User().Get(&users, map[string]interface{}{"first_name": user2FN}, 0, 0)
		if err != nil {
			t.Errorf("error get's users err:%v", err)
		}

		// t.Errorf("ID:%v ID2:%v PASS:%v PASS2:%v", users[0].Id, users[1].Id, users[0].Password, users[1].Password, err)
	})
}

// func (c *TestReposUserCaller) TestUserCount(repos *repos.ReposContext, t *testing.T) {
// 	t.Run("test repos user count", func(t *testing.T) {
// 		if count, err := repos.User().Count(nil); err != nil {
// 			t.Errorf("error count users retCnt:%v err:%v", count, err)
// 		}
// 	})
// }

func (c *TestReposUserCaller) TestUserModify(repos *repos.ReposContext, t *testing.T) {
	user := db.User{
		Id:        userId,
		FirstName: user2FN,
		LastName:  user2LN,
		Login:     userLogin,
		Password:  UserPass,
	}
	t.Run("test repos modify user", func(t *testing.T) {
		if retId, err := repos.User().Modify(&user); err != nil || retId != userId {
			t.Errorf("error modify user userId:%v retId:%v err:%v", userId, retId, err)
		}
	})
}

func (c *TestReposUserCaller) TestUserDelete(repos *repos.ReposContext, t *testing.T) {
	t.Run("test repos delete users", func(t *testing.T) {
		if retId, err := repos.User().Delete(userId); err != nil {
			t.Errorf("error delete user userId:%v retId:%v err:%v", userId, retId, err)
		}
		if retId, err := repos.User().Delete(user2Id); err != nil {
			t.Errorf("error delete user user2Id:%v retId:%v err:%v", user2Id, retId, err)
		}
	})
}
func (c *TestReposUserCaller) TestUserCmpPass(repos *repos.ReposContext, password string, t *testing.T) {
	user := db.User{
		Id: userId,
	}

	t.Run(fmt.Sprintf("test repos cmp pass:%v", password), func(t *testing.T) {
		if !repos.User().ComparePass(userId, password) {
			t.Errorf("error compare pass (string) userId:%v pass:%v", userId, password)
		}
		if !repos.User().ComparePass(&user, password) {
			t.Errorf("error compare pass (struct) userId:%v pass:%v", userId, password)
		}
		// hasher := hasher.NewHasher()
		// up1, _ := hasher.HashStr(UserPass)
		// up2, _ := hasher.HashStr(UserPass2)
		// t.Errorf("hasher str:%v hashed:%v cmp:%v",
		// 	UserPass, up1, hasher.CompareHashAndStr(up1, UserPass))
		// t.Errorf("hasher str:%v hashed:%v cmp:%v",
		// 	UserPass2, up2, hasher.CompareHashAndStr(up2, UserPass2))
	})
}
