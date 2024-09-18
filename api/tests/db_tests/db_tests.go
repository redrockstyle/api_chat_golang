package dbtests

import (
	"api_chat/api/layers/controller/database"
	dbtestschat "api_chat/api/tests/db_tests/db_tests_chat"
	dbtestschatuser "api_chat/api/tests/db_tests/db_tests_chatuser"
	dbtestsmessage "api_chat/api/tests/db_tests/db_tests_message"
	dbtestsuser "api_chat/api/tests/db_tests/db_tests_user"
	"testing"
)

type TestDatabase struct {
	tku  *dbtestsuser.TestsUser
	tkc  *dbtestschat.TestChat
	tkcu *dbtestschatuser.TestsChatUser
	tkm  *dbtestsmessage.TestsMessage
}

func NewTestDatabase(dbc *database.DbController, t *testing.T) *TestDatabase {
	return &TestDatabase{
		tku:  dbtestsuser.NewTestsUser(dbc, t),
		tkc:  dbtestschat.NewTestChat(dbc, t),
		tkcu: dbtestschatuser.NewTestsChatUser(dbc, t),
		tkm:  dbtestsmessage.NewTestsMessage(dbc, t),
	}
}

func (td *TestDatabase) TestsAll() {
	td.TestsUser()
	td.TestChat()
	td.TestChatUser()
	td.TestMessage()
}

func (td *TestDatabase) TestsUser() {
	td.tku.TestsAll()
}

func (td *TestDatabase) TestChat() {
	td.tkc.TestsAll()
}

func (td *TestDatabase) TestChatUser() {
	td.tkcu.TestsAll()
}

func (td *TestDatabase) TestMessage() {
	td.tkm.TestsAll()
}
