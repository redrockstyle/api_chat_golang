package dbtestsmessage

import (
	"api_chat/api/layers/controller/database"
	"api_chat/api/layers/domain/db"
	"fmt"
	"testing"
)

type TestsMessage struct {
	dbc    *database.DbController
	t      *testing.T
	table  string
	user   *db.User
	chat   *db.Chat
	msg    *db.Message
	cu     *db.ChatUser
	cm     *db.ChatMessage
	chatId uint64
	userId string
	msgId  uint64
	cuId   uint64
	cmId   uint64
}

func NewTestsMessage(dbc *database.DbController, t *testing.T) *TestsMessage {
	user := db.User{
		Id:        "1a2b3c4d5e6f",
		FirstName: "fname",
		LastName:  "lname",
		Login:     "somelogin",
		Password:  "somepass",
	}
	chat := db.Chat{
		Desc:    "Nevermind",
		Creator: "1a2b3c4d5e6f",
	}
	msg := db.Message{
		IdUser: user.Id,
		Text:   "Something text",
	}
	cu := db.ChatUser{
		IdUser: user.Id,
	}
	cm := db.ChatMessage{}
	return &TestsMessage{dbc: dbc, t: t, user: &user, chat: &chat, msg: &msg, cu: &cu, cm: &cm}
}

func (tku *TestsMessage) TestsAll() {
	tku.TestCreateTable()
	tku.TestCreateMsg()
	tku.TestGetMsg()
	tku.TestModifyMsg()
	tku.TestDeleteMsg()
	tku.TestDropTable()
}

func (tku *TestsMessage) TestCreateTable() {
	tku.t.Run("create table", func(t *testing.T) {
		id, err := tku.dbc.TypeCreate("", tku.chat)
		if err != nil {
			tku.t.Errorf("error create chat: %v", err)
			return
		}
		tku.chatId = id.(uint64)
		tku.table, err = tku.dbc.TypeTableCreate(tku.chatId, &db.Message{})
		if err != nil {
			tku.t.Errorf("error creeate table: %v", err)
		} else if tku.table != (fmt.Sprintf("msg_%v", tku.chatId)) {
			tku.t.Errorf("error create table: returned %v but expected msg_301", tku.table)
		}
	})
}

func (tku *TestsMessage) TestDropTable() {
	tku.t.Run(fmt.Sprintf("drop table name:%v", tku.table), func(t *testing.T) {
		if err := tku.dbc.TypeTableDropByName(tku.table); err != nil {
			tku.t.Errorf("error drop table: %v", err)
		}
	})
}

func (tku *TestsMessage) TestCreateMsg() {
	tku.t.Run(fmt.Sprintf("create msg table with id:%v", tku.chatId), func(t *testing.T) {
		id, err := tku.dbc.TypeCreate("", tku.user)
		if err != nil {
			tku.t.Errorf("error create user: %v", err)
			return
		}
		tku.userId = id.(string)

		tku.cu.IdChat = tku.chatId
		id, err = tku.dbc.TypeCreate("", tku.cu)
		if err != nil {
			tku.t.Errorf("error create chatuser: %v", err)
			return
		}
		tku.cuId = id.(uint64)

		id, err = tku.dbc.TypeCreate(tku.table, tku.msg)
		if err != nil {
			tku.t.Errorf("error create msg: %v", err)
			return
		}
		tku.msgId = id.(uint64)

		tku.cm.IdChat = tku.chatId
		tku.cm.IdMsg = tku.table
		id, err = tku.dbc.TypeCreate("", tku.cm)
		if err != nil {
			tku.t.Errorf("error create chatmessage: %v", err)
			return
		}
		tku.cmId = id.(uint64)
	})
}

func (tku *TestsMessage) TestGetMsg() {
	msg := db.Message{
		Id: tku.msgId,
	}
	tku.t.Run(fmt.Sprintf("get msg id%v from table id:%v", msg.Id, tku.table), func(t *testing.T) {
		id, err := tku.dbc.TypeGet(tku.table, &msg, nil)
		if err != nil {
			tku.t.Errorf("error get msg from table: %v", err)
		} else if id != tku.msg.Id || msg.IdUser != tku.msg.IdUser || msg.Text != tku.msg.Text {
			tku.t.Errorf("error get msg: error compare: '%v'!='%v' '%v'!='%v' '%v'!='%v'",
				id, tku.msg.Id, msg.IdUser, tku.msg.IdUser, msg.Text, tku.msg.Text,
			)
		}
	})
}

func (tku *TestsMessage) TestModifyMsg() {
	msg := db.Message{
		Id:     tku.msgId,
		IdUser: tku.userId,
		Text:   "somesomesome",
	}
	msgGet := db.Message{
		Id: tku.msgId,
	}
	tku.t.Run(fmt.Sprintf("get msg id%v from table id:%v", msg.Id, tku.table), func(t *testing.T) {
		id, err := tku.dbc.TypeModify(tku.table, &msg)
		if err != nil {
			tku.t.Errorf("error get msg from table: %v", err)
		} else if id.(uint64) != msg.Id {
			tku.t.Errorf("error get msg from table: error compare '%v'!='%v'", id.(uint64), msg.Id)
		}

		id, err = tku.dbc.TypeGet(tku.table, &msgGet, nil)
		if err != nil {
			tku.t.Errorf("error get msg from table: %v", err)
		} else if id.(uint64) != tku.msg.Id || msgGet.IdUser != tku.msg.IdUser || msgGet.Text == tku.msg.Text {
			tku.t.Errorf("error get msg: error compare: '%v'!='%v' '%v'!='%v' (be modified)'%v'=='%v'",
				id.(uint64), tku.msg.Id, msgGet.IdUser, tku.msg.IdUser, msgGet.Text, tku.msg.Text,
			)
		}
	})
}

func (tku *TestsMessage) TestDeleteMsg() {
	tku.t.Run(fmt.Sprintf("delete msg of the table msg_%v", tku.chatId), func(t *testing.T) {
		_, err := tku.dbc.TypeDelete("", tku.chat, nil)
		if err != nil {
			tku.t.Errorf("error delete chat: %v", err)
		}
		_, err = tku.dbc.TypeDelete("", tku.user, nil)
		if err != nil {
			tku.t.Errorf("error delete user: %v", err)
		}
		_, err = tku.dbc.TypeDelete(tku.table, tku.msg, nil)
		if err != nil {
			tku.t.Errorf("error delete msg: %v", err)
		}
		_, err = tku.dbc.TypeDelete("", tku.cu, nil)
		if err != nil {
			tku.t.Errorf("error delete chatuser: %v", err)
		}
		_, err = tku.dbc.TypeDelete("", tku.cm, nil)
		if err != nil {
			tku.t.Errorf("error delete chatmessage: %v", err)
		}
	})
}
