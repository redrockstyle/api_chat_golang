package dbtestschat

import (
	"api_chat/api/layers/controller/database"
	"api_chat/api/layers/domain/db"
	"fmt"
	"testing"
)

type TestChat struct {
	dbc *database.DbController
	id  uint64
	t   *testing.T
}

func NewTestChat(dbc *database.DbController, t *testing.T) *TestChat {
	return &TestChat{dbc: dbc, t: t}
}

func (tku *TestChat) TestsAll() {
	tku.TestCreateChat()
	tku.TestGetChat()
	tku.TestDeleteChat()
}

func (tku *TestChat) TestCreateChat() {
	chat := db.Chat{
		Desc:    "deschat",
		Creator: "1q2w3e4r5t",
	}
	tku.t.Run(fmt.Sprintf("create chat desc:%s", chat.Desc), func(t *testing.T) {
		id, err := tku.dbc.TypeCreate("", &chat)
		if err != nil {
			t.Errorf("error create user: %v", err)
		}
		tku.id = id.(uint64)
	})
}

func (tku *TestChat) TestGetChat() {
	chat := db.Chat{
		Id: tku.id,
	}
	tku.t.Run(fmt.Sprintf("get chat id:%v", chat.Id), func(t *testing.T) {
		id, err := tku.dbc.TypeGet("", &chat, nil)
		if err != nil {
			t.Errorf("error get chat: %v", err)
		} else if id != chat.Id {
			tku.t.Errorf("error get chat: error returned id '%v' != '%v'", id, chat.Id)
		}
	})
}

func (tku *TestChat) TestDeleteChat() {
	chat := db.Chat{
		Id: tku.id,
	}
	tku.t.Run(fmt.Sprintf("delete chat id:%v", chat.Id), func(t *testing.T) {
		id, err := tku.dbc.TypeDelete("", &chat, nil)
		if err != nil {
			t.Errorf("error delete chat: %v", err)
		} else if id != chat.Id {
			tku.t.Errorf("error delete chat: error returned id '%v' != '%v'", id, chat.Id)
		}
	})
}
