package repostestssession

import (
	"api_chat/api/layers/domain/db"
	"api_chat/api/layers/repos"
	"testing"
	"time"
)

type TestReposSerssionCaller struct{}

type TestsReposSession interface {
	TestsSessionAll(repos *repos.ReposContext, t *testing.T)

	TestSessionCreate(repos *repos.ReposContext, t *testing.T)
	TestSessionValidate(repos *repos.ReposContext, t *testing.T)
	TestSessionDestroy(repos *repos.ReposContext, t *testing.T)
}

var session string
var user1Id string

const (
	userFN    = "tom"
	userLN    = "lastnm"
	userLogin = "lgoin"
	UserPass  = "password1"
)

func NewTestReposSerssionCaller() *TestReposSerssionCaller {
	return &TestReposSerssionCaller{}
}

func (c *TestReposSerssionCaller) TestsSessionAll(repos *repos.ReposContext, t *testing.T) {
	c.TestSessionCreate(repos, t)
	c.TestSessionValidate(repos, t)
	c.TestSessionDestroy(repos, t)
}

func (c *TestReposSerssionCaller) TestSessionCreate(repos *repos.ReposContext, t *testing.T) {
	user := db.User{
		FirstName: userFN,
		LastName:  userLN,
		Login:     userLogin,
		Password:  UserPass,
	}
	t.Run("test repos create session", func(t *testing.T) {
		var err error
		if user1Id, err = repos.User().Create(&user); err != nil {
			t.Errorf("error create user err:%v", err)
		}

		if session, err = repos.Session().SessionCreate(user1Id); err != nil {
			t.Errorf("error create err:%v", err)
		}
	})
}

func (c *TestReposSerssionCaller) TestSessionValidate(repos *repos.ReposContext, t *testing.T) {
	t.Run("test repos valid session", func(t *testing.T) {
		if !repos.Session().SessionValidate(session) {
			t.Errorf("error validate before session %v", session)
		}

		time.Sleep(3 * time.Second)

		if repos.Session().SessionValidate(session) {
			t.Errorf("error validate after session %v", session)
		}
	})
}

func (c *TestReposSerssionCaller) TestSessionDestroy(repos *repos.ReposContext, t *testing.T) {
	t.Run("test repos destroy session", func(t *testing.T) {
		if _, err := repos.User().Delete(user1Id); err != nil {
			t.Errorf("error delete user id:%v err:%v", user1Id, err)
		}

		if err := repos.Session().SessionDestroy(session); err != nil {
			t.Errorf("error destroy session err:%v", err)
		}
	})
}
