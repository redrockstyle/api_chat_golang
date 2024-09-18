package repostestsmsg

import (
	"api_chat/api/layers/domain/db"
	"api_chat/api/layers/repos"
	"testing"
)

type TestReposMsgCaller struct{}

type TestsReposMsg interface {
	TestsMsgAll(repos *repos.ReposContext, t *testing.T)

	TestMsgCreate(repos *repos.ReposContext, t *testing.T)
	TestMsgGet(repos *repos.ReposContext, t *testing.T)
	TestMsgModify(repos *repos.ReposContext, t *testing.T)
	TestMsgDelete(repos *repos.ReposContext, t *testing.T)
}

var user1Id string
var user2Id string
var chat1Id uint64

var msg1Id uint64
var msg2Id uint64
var msg3Id uint64
var msg4Id uint64

const (
	// user1
	user1FN    = "tomas"
	user1LN    = "lastnm"
	user1Login = "lgoin"
	User1Pass  = "password1"

	// user2
	user2FN    = "tom"
	user2LN    = "sawyer"
	user2Login = "tomsaw"
	user2Pass  = "passtom"

	// chat1
	chat1Desc = "desc1"

	// msg
	msg1 = "hello, mthfckrs"
	msg2 = "r u gay?"
	msg3 = "i'm not gay"
	msg4 = "no, u"
)

func (c *TestReposMsgCaller) TestsMsgAll(repos *repos.ReposContext, t *testing.T) {
	// c.TestAddUserToChat(repos, t)
	// c.TestAddMsgToChat(repos, t)

	// c.TestModifyMsgByChat(repos, t)
	// c.TestGetCountMessages(repos, t)
	// c.TestGetMessagesFromChat(repos, t)

	// c.TestDelMsgFromChat(repos, t)
	// c.TestDelUserFromChat(repos, t)
}

func NewTestReposMsgCaller() *TestReposMsgCaller {
	return &TestReposMsgCaller{}
}

func (c *TestReposMsgCaller) TestAddUserToChat(repos *repos.ReposContext, t *testing.T) {
	user1 := db.User{
		FirstName: user1FN,
		LastName:  user1LN,
		Login:     user1Login,
		Password:  User1Pass,
	}
	user2 := db.User{
		FirstName: user2FN,
		LastName:  user2LN,
		Login:     user2Login,
		Password:  user2Pass,
	}
	chat1 := db.Chat{
		Desc:    chat1Desc,
		Creator: user1Id,
	}
	var err error
	t.Run("test repos msg", func(t *testing.T) {
		if user1Id, err = repos.User().Create(&user1); err != nil {
			t.Errorf("error create user: %v", err)
		}
		if user2Id, err = repos.User().Create(&user2); err != nil {
			t.Errorf("error create user: %v", err)
		}
		if chat1Id, err = repos.Chat().Create(&chat1); err != nil {
			t.Errorf("error create chat: %v", err)
		}

		if _, err = repos.AddUserToChat(user1Id, chat1Id); err != nil {
			t.Errorf("error adding user:%v to chat err:%v", user1Id, err)
		}
		if _, err = repos.AddUserToChat(user2Id, chat1Id); err != nil {
			t.Errorf("error adding user:%v to chat err:%v", user2Id, err)
		}

		if msg1Id, err = repos.AddMsgToChat(user1Id, chat1Id, msg1); err != nil {
			t.Errorf("error push msg:%v err:%v", msg1Id, err)
		}
		if msg2Id, err = repos.AddMsgToChat(user2Id, chat1Id, msg2); err != nil {
			t.Errorf("error push msg:%v err:%v", msg1Id, err)
		}
		if msg3Id, err = repos.AddMsgToChat(user1Id, chat1Id, msg3); err != nil {
			t.Errorf("error push msg:%v err:%v", msg1Id, err)
		}
		if msg4Id, err = repos.AddMsgToChat(user1Id, chat1Id, msg4); err != nil {
			t.Errorf("error push msg:%v err:%v", msg1Id, err)
		}
	})
}
func (c *TestReposMsgCaller) TestAddMsgToChat(repos *repos.ReposContext, t *testing.T) {}

func (c *TestReposMsgCaller) TestModifyMsgByChat(repos *repos.ReposContext, t *testing.T)     {}
func (c *TestReposMsgCaller) TestGetCountMessages(repos *repos.ReposContext, t *testing.T)    {}
func (c *TestReposMsgCaller) TestGetMessagesFromChat(repos *repos.ReposContext, t *testing.T) {}

func (c *TestReposMsgCaller) TestDelMsgFromChat(repos *repos.ReposContext, t *testing.T) {}
func (c *TestReposMsgCaller) TestDelUserFromChat(repos *repos.ReposContext, t *testing.T) {
	// t.Run("test repos ")
}
