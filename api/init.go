package api

import (
	"api_chat/api/layers/base/hasher"
	"api_chat/api/layers/base/logx"
	database "api_chat/api/layers/controller/database"
	mysqlsql "api_chat/api/layers/controller/database/mysql_sql"
	psql_init "api_chat/api/layers/controller/database/postrge_sql"
	httpserver "api_chat/api/layers/controller/server"
	httphandler "api_chat/api/layers/controller/server/http_handler"
	httprouter "api_chat/api/layers/controller/server/http_router"
	config "api_chat/api/layers/domain/cfg"
	dbt "api_chat/api/layers/domain/db"
	mwrole "api_chat/api/layers/middleware/mw_role"
	"api_chat/api/layers/repos"
	"api_chat/api/layers/usecase"
	"log"
	"time"

	"gorm.io/gorm"
)

type Initializer struct {
	filename string
	log      logx.Logger
	cfg      *config.Configuration
}

func NewInitializerAPI(filename string) *Initializer {
	return &Initializer{filename: filename, cfg: nil, log: nil}
}

func (i *Initializer) InitConfig() {
	vpr, err := config.LoadConfig(i.filename)
	if err != nil {
		log.Fatalf("Error open config: %v", err)
	}
	if i.cfg, err = config.ParseConfig(vpr); err != nil {
		log.Fatalf("Error parse config: %v", err)
	}
	i.log = logx.NewApiLogger(i.cfg)
	i.log.InitLogger()
	i.log.Info("Read config: success")
}

func (i *Initializer) Startup() {
	i.log.Info("Initialization API")
	i.log.Infof("VersionAPI:%v Mode:%v SSL:%v", i.cfg.Server.AppVersion, i.cfg.Server.Mode, i.cfg.SSL.Active)
	i.log.Infof("LoggerLevel:%v", i.cfg.Logger.Level)

	i.log.Info("Open database")
	var db *gorm.DB
	var err error
	switch i.cfg.Database.Driver {
	case config.DriverMysql:
		db, err = mysqlsql.OpenMysqlDatabase(i.cfg)
		if err != nil {
			i.log.Fatalf("Open database is failed: %v", err)
		} else {
			i.log.Info("Success init mysql database")
		}
	case config.DriverPostgres:
		db, err = psql_init.OpenPsqlDatabase(i.cfg)
		if err != nil {
			i.log.Fatalf("Open database is failed: %v", err)
		} else {
			i.log.Info("Success init psql database")
		}
	default:
		i.log.Fatalf("Database name %v is not supported", i.cfg.Database.Driver)
		return
	}

	i.log.Info("Mirgrate database")
	dbc := database.NewDbController(i.cfg, db, i.log)
	if i.cfg.Database.CleanStart {
		if err := dbc.CleanInit(); err != nil {
			i.log.Fatal("CleanInit database is failed: %v", err)
		}
	} else {
		if err := dbc.Init(); err != nil {
			i.log.Fatal("Init database is failed: %v", err)
		}
	}

	i.log.Info("Init repos")
	repos := repos.NewReposContext(dbc, hasher.NewHasher(), i.log, time.Duration(i.cfg.Server.SessionLifeTimeMin*int64(time.Minute)))

	if i.cfg.Database.CreateAdmin {
		i.log.Infof("Create user:%v role:%v", i.cfg.Admin.Username, mwrole.RoleAdmin)
		if _, err := repos.User().Create(&dbt.User{Login: i.cfg.Admin.Username, Password: i.cfg.Admin.Password, Role: mwrole.RoleAdmin}); err != nil {
			i.log.Warnf("Error create %v: %v", i.cfg.Admin.Username, err)
		}
	}

	i.log.Info("Init usecase")
	uc := usecase.NewUsecaseOparetor(i.cfg, repos, i.log)

	i.log.Info("Init HTTP server")
	handler := httphandler.NewRestApiHandler(i.cfg, uc, i.log)
	router := httprouter.NewApiRouter(i.cfg, handler)
	server := httpserver.NewRestApiServer(router, i.cfg, i.log)
	server.Runtime()
}
