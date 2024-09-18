package repostests

import (
	"api_chat/api/layers/base/hasher"
	"api_chat/api/layers/base/logx"
	"api_chat/api/layers/controller/database"
	"api_chat/api/layers/repos"
	repostestschat "api_chat/api/tests/repos_tests/repos_tests_chat"
	repostestsmsg "api_chat/api/tests/repos_tests/repos_tests_msg"
	repostestssession "api_chat/api/tests/repos_tests/repos_tests_session"
	repostestsuser "api_chat/api/tests/repos_tests/repos_tests_user"
	"testing"
	"time"
)

type TestsRepos struct {
	t        *testing.T
	dbc      *database.DbController
	hasher   *hasher.Hasher
	lifetime time.Duration
	repos    *repos.ReposContext
	tru      *repostestsuser.TestReposUserCaller
	trc      *repostestschat.TestReposChatCaller
	trs      *repostestssession.TestReposSerssionCaller
	trm      *repostestsmsg.TestReposMsgCaller
}

func NewTestsRepos(dbc *database.DbController, hasher *hasher.Hasher, lifetime time.Duration, t *testing.T, logx logx.Logger) *TestsRepos {
	return &TestsRepos{
		dbc:      dbc,
		hasher:   hasher,
		lifetime: lifetime,
		t:        t,
		repos:    repos.NewReposContext(dbc, hasher, logx, lifetime),
		tru:      repostestsuser.NewTestReposUserCaller(),
		trc:      repostestschat.NewTestReposChatCaller(),
		trs:      repostestssession.NewTestReposSerssionCaller(),
		trm:      repostestsmsg.NewTestReposMsgCaller(),
	}
}

func (tr *TestsRepos) TestsAll() {
	tr.TestUser()
	tr.TestChat()
	tr.TestMsg()
	tr.TestSession()
}

func (tr *TestsRepos) TestUser() {
	tr.tru.TestsUserAll(tr.repos, tr.t)
}

func (tr *TestsRepos) TestChat() {
	tr.trc.TestsChatAll(tr.repos, tr.t)
}

func (tr *TestsRepos) TestMsg() {
	tr.trs.TestsSessionAll(tr.repos, tr.t)
}

func (tr *TestsRepos) TestSession() {
	tr.trm.TestsMsgAll(tr.repos, tr.t)
}
