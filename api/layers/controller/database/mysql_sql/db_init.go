package mysqlsql

import (
	config "api_chat/api/layers/domain/cfg"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func OpenMysqlDatabase(cfg *config.Configuration) (db *gorm.DB, err error) {

	dataSourceName := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local",
		cfg.Mysql.MysqlUser,
		cfg.Mysql.MysqlPass,
		cfg.Mysql.MysqlHost,
		cfg.Mysql.MysqlPort,
		cfg.Mysql.MysqlDB,
	)
	gormDB, err := gorm.Open(mysql.Open(dataSourceName), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return gormDB, nil
}
