package database

import (
	"api_chat/api/layers/base/logx"
	config "api_chat/api/layers/domain/cfg"
	"api_chat/api/layers/domain/db"
	"errors"
	"fmt"
	"reflect"

	"gorm.io/gorm"
)

const (
	MaxOffsetValue = 50
	MaxLimitValue  = 100
)

type DbController struct {
	db  *gorm.DB
	cfg *config.Configuration
	log logx.Logger
}

func NewDbController(cfg *config.Configuration, db *gorm.DB, logx logx.Logger) *DbController {
	return &DbController{cfg: cfg, db: db, log: logx}
}

func (dbc *DbController) Init() error {
	dbc.log.Debugf("Init called")
	if err := MigrateTable(dbc.db,
		db.User{}, db.Chat{}, db.ChatUser{}, db.ChatMessage{}, db.Session{},
	); err != nil {
		return err
	}
	return nil
}

func (dbc *DbController) CleanInit() error {
	dbc.log.Debug("CleanInit called")
	var tables []string
	if err := dbc.db.Table("information_schema.tables").Where("table_schema = ?", "public").Pluck("table_name", &tables).Error; err != nil {
		return err
	}
	// func(dbg *gorm.DB, args ...interface{}) {
	// 	lenArgs := len(args)
	// 	for i := 0; i < lenArgs; i++ {
	// 		if dbg.Migrator().HasTable(&args[i]) {
	// 			dbg.Migrator().DropTable(&args[i])
	// 		}
	// 	}
	// }(dbc.db, db.User{}, db.Chat{}, db.ChatUser{}, db.ChatMessage{}, db.Session{})
	lenArgs := len(tables)
	for i := 0; i < lenArgs; i++ {
		dbc.db.Migrator().DropTable(tables[i])
	}
	return dbc.Init()
}

/*
 * Add record to the db (struct db: User, Chat, ChatUser, Message)
 * Returned id (string) or error if add is failed
 */
func (dbc *DbController) TypeCreate(table string, t interface{}) (interface{}, error) {
	dbc.log.Debugf("Create called with %v", t)
	if err := ValidateStruct(t); err != nil {
		return nil, err
	}
	if err := dbc.db.Table(table).Create(t); err.Error != nil {
		return nil, err.Error
	}
	return GetIdField(t), nil
}

/*
 * Del record from db (struct db: User, Chat, ChatUser, Message)
 * Returned id (string) or error if add is failed
 */
func (dbc *DbController) TypeDelete(table string, t interface{}, conds map[string]interface{}) (interface{}, error) {
	dbc.log.Debugf("Delete called with %v", t)
	if err := ValidateStruct(t); err != nil {
		return nil, err
	}
	if err := dbc.db.Table(table).Where(conds).Delete(t); err.Error != nil {
		return nil, err.Error
	}
	return GetIdField(t), nil
}

/*
 * Modify record from db (struct db: User, Chat, ChatUser, Message)
 * Returned id (string) or error if mod is failed
 */
func (dbc *DbController) TypeModify(table string, t interface{}) (interface{}, error) {
	dbc.log.Debugf("Modify called with %v", t)
	if err := ValidateStruct(t); err != nil {
		return nil, err
	}
	if err := dbc.db.Table(table).Save(t); err.Error != nil {
		return nil, err.Error
	}
	return GetIdField(t), nil
}

func (dbc *DbController) TypeModifyCol(table string, t interface{}, conds map[string]interface{}) (interface{}, error) {
	dbc.log.Debugf("Mdify_col called with %v", t)
	if err := ValidateStruct(t); err != nil {
		return nil, err
	}
	if conds == nil {
		return nil, errors.New("conds most be nil")
	}
	if ctx := dbc.db.Table(table).Model(t).Updates(conds); ctx.Error != nil {
		return nil, ctx.Error
	}
	return GetIdField(t), nil
}

/*
 * Get record from db (struct db: User, Chat, ChatUser, Message)
 * Returned id (string) or error if get is failed
 */
func (dbc *DbController) TypeGet(table string, t interface{}, conds map[string]interface{}) (interface{}, error) {
	dbc.log.Debugf("Get called with %v", t)
	if err := ValidateStruct(t); err != nil {
		return nil, err
	}
	val := reflect.ValueOf(t)
	if val.Kind() == reflect.Ptr && val.Elem().Kind() == reflect.Struct {
		if !dbc.TypeExists(table, t, conds) {
			return nil, errors.New("record is not found")
		}
		if err := dbc.db.Table(table).Where(conds).First(t); err.Error != nil {
			return nil, err.Error
		}
		return GetIdField(t), nil
	} else {
		return nil, fmt.Errorf("error data type %v", val.Kind())
	}
}

/*
 * Get records from db (struct db: User, Chat, ChatUser, Message)
 * conds example: {"key":"value"}
 * ids example: []int{1,2,3}
 * Returned count writed or error if get is failed
 */
func (dbc *DbController) TypeGets(table string, t interface{}, conds map[string]interface{}, offset int, limit int) (int64, error) {
	dbc.log.Debugf("Gets called with t:%v conds:%v offset:%v limit:%v", t, conds, offset, limit)
	if t == nil {
		return 0, errors.New("passed nil")
	}
	val := reflect.ValueOf(t)
	if val.Kind() == reflect.Ptr && val.Elem().Kind() == reflect.Slice {
		if offset > MaxOffsetValue || offset < 0 {
			dbc.log.Warnf("Invalid using offset:%v set 0", offset)
			offset = 0
		}
		if limit > MaxLimitValue || limit <= 0 {
			dbc.log.Warnf("Invalid using limit:%v set 10", limit)
			limit = 10
		}
		ctx := dbc.db.Table(table).Model(t).Where(conds).Offset(offset).Limit(limit).Find(t)
		if ctx.Error != nil {
			return 0, ctx.Error
		}
		return ctx.RowsAffected, nil
	} else {
		return 0, fmt.Errorf("error data type %v", val.Kind())
	}
}

/*
 * Chech value exists in the db
 * conds has a json format {"key":"value"}
 */
func (dbc *DbController) TypeExists(table string, t interface{}, conds map[string]interface{}) bool {
	dbc.log.Debugf("Exists called with %v %v %v", table, t, conds)
	if err := ValidateStruct(t); err != nil {
		return false
	}
	var exists bool
	if dbc.db.Table(table).Model(t).Select("count(*) > 0").Where(conds).Find(&exists).Error != nil {
		return false
	}
	return exists
}

/*
 * Get type count
 * Proc: use table name "t.(type)_suffix" or suffix ignored if suffix==0
 * Return: count rows in db or error if get count is failed
 */
func (dbc *DbController) TypeCount(suffix uint64, t interface{}, conds map[string]interface{}) (int64, error) {
	dbc.log.Debugf("Count called with %v", t)
	if t == nil {
		return 0, errors.New("passed nil")
	}
	if err := ValidateStruct(t); err != nil {
		return 0, err
	}
	val := reflect.ValueOf(t)
	if val.Kind() == reflect.Ptr && val.Elem().Kind() == reflect.Struct {
		var err error
		count := int64(0)
		tableName := ""
		if suffix != 0 {
			tableName = MuTableStr(t, fmt.Sprintf("%v", suffix))
		}
		err = dbc.db.Table(tableName).Model(t).Where(conds).Count(&count).Error
		if err != nil {
			return 0, err
		}
		return count, nil
	} else {
		return 0, fmt.Errorf("error data type %v", val.Kind())
	}
}

/*
 * Create table name: "t.(type)_suffix"
 * Return: name of the created table or error if create is failed
 */
func (dbc *DbController) TypeTableCreate(suffix uint64, t interface{}) (string, error) {
	dbc.log.Debugf("TableCreate called with %v", t)
	if err := ValidateStruct(t); err != nil {
		return "", err
	}
	tableName := MuTableStr(t, fmt.Sprintf("%v", suffix))
	//if err := dbc.db.Scopes(MuTable(t, suffix)).Create(t); err.Error != nil {
	if err := dbc.db.Table(tableName).AutoMigrate(t); err != nil {
		return "", err
	}
	return tableName, nil
}

func (dbc *DbController) TypeTableDropByName(name string) error {
	return dbc.TypeTableDrop(name)
}

func (dbc *DbController) TypeTableDrop(t interface{}) error {
	dbc.log.Debugf("TableDrop called with %v", t)
	return dbc.db.Migrator().DropTable(t)
}

/*
 * Is exists table for name "t.(type)_suffix"
 */
func (dbc *DbController) TypeTableIsExists(suffix uint64, t interface{}) (string, error) {
	dbc.log.Debugf("TableIsExists called with %v", t)
	if err := ValidateStruct(t); err != nil {
		return "", err
	}
	tableName := MuTableStr(t, fmt.Sprintf("%v", suffix))
	if dbc.db.Migrator().HasTable(tableName) {
		return tableName, nil
	}
	return "", errors.New("table is not found")
}

func MigrateTable(db *gorm.DB, t ...interface{}) error {
	return db.AutoMigrate(t...)
}

/*
 * Get tx.Table() with "t.(type)_suffix"
 */
func MuTableTx(t interface{}, suffix string) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Table(MuTable(t, suffix))
	}
}

/*
 * Get string: "t.(type)_suffix"
 */
func MuTableStr(t interface{}, suffix string) string {
	return MuTable(t, suffix)
}

func MuTable(t interface{}, suffix string) string {
	switch t.(type) {
	case *db.User:
		return "user_" + suffix
	case *db.Chat:
		return "chat_" + suffix
	case *db.Message:
		return "msg_" + suffix
	case *db.ChatMessage:
		return "cm_" + suffix
	case *db.ChatUser:
		return "cu_" + suffix
	default:
		return ""
	}
}

/*
 * Get id field from interface (reflect)
 * https://stackoverflow.com/questions/49460380/recursive-struct-reflection-error-panic-reflect-field-of-non-struct-type
 */
func GetIdField(t interface{}) interface{} {
	if t == nil {
		return nil
	}
	val := reflect.ValueOf(t)
	if val.Kind() == reflect.Ptr && val.Elem().Kind() == reflect.Struct {
		val = val.Elem()
	} else {
		return nil
	}

	switch t.(type) {
	case *db.User:
		return val.Field(0).String()
	case *db.Session:
		return val.Field(0).String()
	default:
		return val.Field(0).Uint()
	}
}

func ValidateStruct(t interface{}) error {
	if t != nil {
		val := reflect.ValueOf(t)
		if val.Kind() == reflect.Ptr && val.Elem().Kind() == reflect.Struct {
			switch t.(type) {
			case *db.User:
				return nil
			case *db.Chat:
				return nil
			case *db.ChatUser:
				return nil
			case *db.ChatMessage:
				return nil
			case *db.Message:
				return nil
			case *db.Session:
				return nil
			case *[]db.User:
				return nil
			case *[]db.Chat:
				return nil
			case *[]db.ChatUser:
				return nil
			case *[]db.ChatMessage:
				return nil
			case *[]db.Message:
				return nil
			case *[]db.Session:
				return nil
			default:
				return fmt.Errorf("unsupported struct '%v' '%v'", reflect.TypeOf(t).Kind(), reflect.TypeOf(t).Elem())
			}
		} else {
			return fmt.Errorf("unsupported format type: %v", val.Kind())
		}
	} else {
		return errors.New("passed nil")
	}
}
