package db

const (
	PrefixUsrTable = "usr"
	PrefixChtTable = "cht"
	PrefixMsgTable = "msg"
)

type User struct {
	Id        string `json:"-" gorm:"size:36"` // sizeof(UUID) == 36 bytes
	FirstName string `json:"first_name" gorm:"size:255"`
	LastName  string `json:"last_name" gorm:"size:255"`
	Login     string `json:"login" validate:"required" gorm:"size:36;unique"`
	Password  string `json:"password" validate:"required" gorm:"size:255"`
	Role      string `json:"role" gorm:"size:16"`
	//gorm.Model
}

type Chat struct {
	Id      uint64 `json:"id" gorm:"unique;primaryKey;autoIncrement"`
	Desc    string `json:"desc" validate:"required" gorm:"unique;size:255"`
	Creator string `json:"creator" validate:"required" gorm:"size:36"`
	Role    string `json:"role" gorm:"16"`
	//gorm.Model
}

type ChatUser struct {
	Id     uint64 `json:"id" gorm:"unique;primaryKey;autoIncrement"`
	IdUser string `json:"id_user" validate:"required" gorm:"size:36"`
	IdChat uint64 `json:"id_chat" validate:"required,numeric" gorm:"size:32"`
}

type ChatMessage struct {
	Id     uint64 `json:"id" gorm:"unique;primaryKey;autoIncrement"`
	IdChat uint64 `json:"id_chat" validate:"required,numeric" gorm:"size:32"`
	IdMsg  string `json:"id_msg" validate:"required,numeric" gorm:"32"`
}

type Message struct {
	Id     uint64 `json:"id" gorm:"unique;primaryKey;autoIncrement"`
	IdUser string `json:"id_user" validate:"required" gorm:"size:36"`
	Text   string `json:"text" validate:"required" gorm:"511"`
}

type Session struct {
	Id     string `json:"id" validate:"required" gorm:"unique;size:32"`
	IdUser string `json:"id_user" validate:"required" gorm:"unique;size:36"`
	Time   int64  `json:"time" validate:"required"`
}
