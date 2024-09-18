package cfg

import (
	"errors"
	"log"

	"github.com/spf13/viper"
)

const (
	DriverPostgres = "psql"
	DriverMysql    = "mysql"
)

type Configuration struct {
	Postgres PostgresConfig
	Mysql    MysqlConfig
	Server   ServerConfig
	SSL      SslConfig
	Database DatabaseConfig
	Logger   LoggerConfig
	Admin    AdminConfig
}

type ServerConfig struct {
	AppVersion         string
	Port               string
	Domain             string
	Mode               string
	ServerName         string
	PrefixPath         string
	Registration       bool
	SessionLifeTimeMin int64
}

type SslConfig struct {
	Active  bool
	CrtPath string
	KeyPath string
}

type DatabaseConfig struct {
	Driver      string
	CleanStart  bool
	CreateAdmin bool
}

type PostgresConfig struct {
	PostgresqlHost string
	PostgresqlPort string
	PostgresqlUser string
	PostgresqlPass string
	PostgresqlDB   string
}

type MysqlConfig struct {
	MysqlHost string
	MysqlPort string
	MysqlUser string
	MysqlPass string
	MysqlDB   string
}

type LoggerConfig struct {
	Development       bool
	DisableCaller     bool
	DisableStacktrace bool
	Encoding          string
	Level             string
}

type AdminConfig struct {
	Username string
	Password string
	Role     string
}

func LoadConfig(filename string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config not found")
		}
		return nil, err
	}
	return v, nil
}

func ParseConfig(v *viper.Viper) (*Configuration, error) {
	var c Configuration
	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("error decode struct, %v", err)
		return nil, err
	}
	return &c, nil
}
