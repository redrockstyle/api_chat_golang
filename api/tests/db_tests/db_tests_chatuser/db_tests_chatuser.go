package dbtestschatuser

import (
	"api_chat/api/layers/controller/database"
	"api_chat/api/layers/domain/db"
	"fmt"
	"testing"
)

type TestsChatUser struct {
	dbc *database.DbController
	t   *testing.T
	id  uint64
}

func NewTestsChatUser(dbc *database.DbController, t *testing.T) *TestsChatUser {
	return &TestsChatUser{dbc: dbc, t: t}
}

func (tku *TestsChatUser) TestsAll() {
	tku.TestCreateChatUser()
	tku.TestGetChatUser()
	tku.TestModifyChatUser()
	tku.TestDeleteChatUser()
}

func (tku *TestsChatUser) TestCreateChatUser() {
	cu := db.ChatUser{
		IdChat: 100,
		IdUser: "1a2b3c4d5e",
	}
	tku.t.Run(fmt.Sprintf("craete cu idchat:%v iduser:%v", cu.IdChat, cu.IdUser), func(t *testing.T) {
		id, err := tku.dbc.TypeCreate("", &cu)
		if err != nil {
			tku.t.Errorf("error create cu: %v", err)
		}
		tku.id = id.(uint64)
	})
}

func (tku *TestsChatUser) TestGetChatUser() {
	cu := db.ChatUser{
		Id: tku.id,
	}
	tku.t.Run(fmt.Sprintf("get cu by id:%v", cu.Id), func(t *testing.T) {
		id, err := tku.dbc.TypeGet("", &cu, nil)
		if err != nil {
			tku.t.Errorf("error get cu: %v", err)
		} else if id != cu.Id {
			tku.t.Errorf("error get cu: error returned id '%v' != '%v'", id, cu.Id)
		}
	})
}

func (tku *TestsChatUser) TestModifyChatUser() {
	cu := db.ChatUser{
		Id:     tku.id,
		IdChat: 9999,
		IdUser: "1a2b3c4d5e6f",
	}
	cuGet := db.ChatUser{
		Id: tku.id,
	}
	tku.t.Run(fmt.Sprintf("get modify by id:%v", cu.Id), func(t *testing.T) {
		id, err := tku.dbc.TypeModify("", &cu)
		if err != nil {
			tku.t.Errorf("error modify cu: %v", err)
		} else if id.(uint64) != cu.Id {
			tku.t.Errorf("error modify cu: error returned id '%v' != '%v'", id.(uint64), cu.Id)
		}

		_, err = tku.dbc.TypeGet("", &cuGet, nil)
		if err != nil {
			tku.t.Errorf("error get cu: %v", err)
		} else if cu != cuGet {
			tku.t.Error("error modify: fields is not modified")
		}
	})
}

func (tku *TestsChatUser) TestDeleteChatUser() {
	cu := db.ChatUser{
		Id: tku.id,
	}
	tku.t.Run(fmt.Sprintf("delete cu by id:%v", cu.Id), func(t *testing.T) {
		id, err := tku.dbc.TypeDelete("", &cu, nil)
		if err != nil {
			tku.t.Errorf("error gelete cu: %v", err)
		} else if id.(uint64) != cu.Id {
			tku.t.Errorf("error delete cu: error returned id '%v' != '%v'", id.(uint64), cu.Id)
		}
	})
}
