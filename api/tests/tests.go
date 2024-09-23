package tests

import (
	"api_chat/api/layers/base/hasher"
	"api_chat/api/layers/base/logx"
	"api_chat/api/layers/controller/database"
	"api_chat/api/layers/controller/database/postrge_sql"
	config "api_chat/api/layers/domain/cfg"
	dbtests "api_chat/api/tests/db_tests"
	repostests "api_chat/api/tests/repos_tests"
	"testing"
	"time"

	"gorm.io/gorm"
)

const config_path = "../../config/config.dev"

type TestKeys struct {
	dbGorm *gorm.DB
	dbc    *database.DbController
	cfg    *config.Configuration
	log    logx.Logger
	t      *testing.T
}

func NewTestsKeys(t *testing.T) (*TestKeys, error) {
	var tk TestKeys
	tk.t = t
	vpr, err := config.LoadConfig(config_path)
	if err != nil {
		return nil, err
	}
	tk.cfg, err = config.ParseConfig(vpr)
	if err != nil {
		return nil, err
	}
	//tk.dbGorm, err = mysql_sql.OpenMysqlDatabase(tk.cfg)
	tk.dbGorm, err = postrge_sql.OpenPsqlDatabase(tk.cfg)
	if err != nil {
		return nil, err
	}

	tk.cfg.Logger.DisableCaller = true
	tk.cfg.Logger.Development = false
	tk.cfg.Logger.DisableStacktrace = true
	tk.cfg.Logger.Level = "fatal"
	tk.log = logx.NewApiLogger(tk.cfg)
	tk.log.InitLogger()
	tk.dbc = database.NewDbController(tk.cfg, tk.dbGorm, tk.log)
	tk.dbc.Init()
	return &tk, nil
}

func (tk *TestKeys) TestDatabase() {
	tk.t.Run("tests database", func(t *testing.T) {
		tdb := dbtests.NewTestDatabase(tk.dbc, tk.t)

		tdb.TestsAll()
	})
}

func (tk *TestKeys) TestRepos() {
	tk.t.Run("tests repository", func(t *testing.T) {
		tr := repostests.NewTestsRepos(
			tk.dbc,
			hasher.NewHasher(),
			time.Duration(1*time.Second),
			tk.t,
			tk.log,
		)

		tr.TestsAll()
	})
}
