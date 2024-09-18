package repostestschat

import (
	"api_chat/api/layers/domain/db"
	"api_chat/api/layers/repos"
	"testing"
)

type TestReposChatCaller struct{}

type TestsReposChat interface {
	TestsChatAll(repos *repos.ReposContext, t *testing.T)

	TestChatCreate(repos *repos.ReposContext, t *testing.T)
	TestChatGet(repos *repos.ReposContext, t *testing.T)
	TestChatModify(repos *repos.ReposContext, t *testing.T)
	TestChatDelete(repos *repos.ReposContext, t *testing.T)
}

var chat1Id uint64
var chat2Id uint64

const (
	chatDesc1    = "desc1"
	chatDesc2    = "desc2"
	chatDesc3    = "desc3"
	chatCreator1 = "1q2w3e4r5t"
	chatCreator2 = "5t4r3e2w1q"
)

func NewTestReposChatCaller() *TestReposChatCaller {
	return &TestReposChatCaller{}
}

func (c *TestReposChatCaller) TestsChatAll(repos *repos.ReposContext, t *testing.T) {
	c.TestChatCreate(repos, t)
	c.TestChatGet(repos, t)
	c.TestChatCount(repos, t)
	c.TestChatModify(repos, t)
	c.TestChatDelete(repos, t)
}

func (c *TestReposChatCaller) TestChatCreate(repos *repos.ReposContext, t *testing.T) {
	chat1 := db.Chat{
		Desc:    chatDesc1,
		Creator: chatCreator1,
	}
	chat2 := db.Chat{
		Desc:    chatDesc2,
		Creator: chatCreator1,
	}
	var err error
	t.Run("test repos create chat", func(t *testing.T) {
		if chat1Id, err = repos.Chat().Create(&chat1); err != nil {
			t.Errorf("error create chat desc:%v err:%v", chatDesc1, err)
		}

		if chat2Id, err = repos.Chat().Create(&chat2); err != nil {
			t.Errorf("error create chat desc:%v err:%v", chatDesc2, err)
		}
	})
}
func (c *TestReposChatCaller) TestChatGet(repos *repos.ReposContext, t *testing.T) {
	chat := db.Chat{Id: chat1Id}
	t.Run("test repos get chat", func(t *testing.T) {
		if retChat, err := repos.Chat().Get(chat1Id, nil, 0, 0); err != nil ||
			retChat.(db.Chat).Id != chat1Id ||
			retChat.(db.Chat).Desc != chatDesc1 {
			t.Errorf("error get chat (string) Id:%v desc:%v err:%v", chat1Id, chatDesc1, err)
		}

		if retChat, err := repos.Chat().Get(&chat, nil, 0, 0); err != nil ||
			retChat.(db.Chat).Id != chat1Id || retChat.(db.Chat).Desc != chatDesc1 {
			t.Errorf("error get chat (struct) id:%v desc:%v err:%v", chat1Id, chatDesc1, err)
		}
	})
}

func (c *TestReposChatCaller) TestChatCount(repos *repos.ReposContext, t *testing.T) {
	t.Run("test repos user count", func(t *testing.T) {
		if count, err := repos.Chat().Count(nil); err != nil {
			t.Errorf("error count chats retCnt:%v err:%v", count, err)
		}
	})
}

func (c *TestReposChatCaller) TestChatModify(repos *repos.ReposContext, t *testing.T) {
	chat := db.Chat{
		Id:      chat1Id,
		Desc:    chatDesc3,
		Creator: chatCreator2,
	}
	t.Run("repos modify chat", func(t *testing.T) {
		if id, err := repos.Chat().Modify(&chat); err != nil || id != chat1Id {
			t.Errorf("error modify chat retId:%v id:%v err:%v", id, chat1Id, err)
		}
	})
}

func (c *TestReposChatCaller) TestChatDelete(repos *repos.ReposContext, t *testing.T) {
	t.Run("delete chat", func(t *testing.T) {
		if id, err := repos.Chat().Delete(&db.Chat{Id: chat1Id}); err != nil || id != chat1Id {
			t.Errorf("error delete chat chat1Id:%v != id:%v err:%v", chat1Id, id, err)
		}

		if id, err := repos.Chat().Delete(chat2Id); err != nil || id != chat2Id {
			t.Errorf("error delete chat chat2Id:%v != id:%v err:%v", chat2Id, id, err)
		}
	})
}
