package sessionctx

import (
	"api_chat/api/layers/base/hasher"
	"api_chat/api/layers/base/ident"
	"api_chat/api/layers/controller/database"
	"api_chat/api/layers/domain/db"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"sync"
	"time"
)

type SessionCtx struct {
	dbc      *database.DbController
	lock     sync.Mutex
	hash     *hasher.Hasher
	lifetime time.Duration
	ident    *ident.Ident
}

func NewSessionCtx(dbc *database.DbController, hash *hasher.Hasher, lifetime time.Duration) *SessionCtx {
	return &SessionCtx{dbc: dbc, hash: hash, lifetime: lifetime, ident: ident.NewIdent()}
}

/*
 * Create session value
 * Return: 32-bytes value or error if creation is failed
 */
func (sc *SessionCtx) SessionCreate(userId string) (string, error) {
	if userId == "" || sc.ident.CheckUUIDv5(userId) {
		return "", errors.New("invalid user id format")
	}
	if _, err := sc.dbc.TypeGet("", &db.User{Id: userId}, nil); err != nil {
		return "", err
	}

	sessAct := db.Session{}
	if sc.dbc.TypeExists("", &sessAct, map[string]interface{}{"id_user": userId}) {
		_ = sc.SessionDestroy(sessAct.Id)
	}

	sessionId := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, sessionId); err != nil {
		return "", errors.New("create session is failed")
	}
	retSessionId := hex.EncodeToString(sessionId)
	if _, err := sc.dbc.TypeCreate(
		"",
		&db.Session{Id: retSessionId, IdUser: userId, Time: time.Now().UnixMilli()},
	); err != nil {
		return "", err
	}
	return retSessionId, nil
}

/*
 * Validate session value
 * Proc: used sync.Mutex.Lock()
 * Return: true if session is correct or false if session is incorrect
 */
func (sc *SessionCtx) SessionValidate(session string) bool {
	sc.lock.Lock()
	defer sc.lock.Unlock()
	if session == "" {
		return false
	}
	// if !sc.dbc.TypeExists("", &db.Session{}, map[string]interface{}{"id": session}) {
	// 	return false
	// }

	sess := db.Session{}
	if _, err := sc.dbc.TypeGet("", &sess, map[string]interface{}{"id": session}); err != nil {
		return false
	}
	if time.Now().UnixMilli()-sess.Time > sc.lifetime.Milliseconds() {
		_ = sc.SessionDestroy(session)
		return false
	}
	return true
}

// used SessionValidate
func (sc *SessionCtx) SessionGetUserId(session string) (string, error) {
	if !sc.SessionValidate(session) {
		return "", errors.New("invalid session")
	}

	sess := db.Session{Id: session}
	if _, err := sc.dbc.TypeGet("", &sess, nil); err != nil {
		return "", err
	}
	return sess.IdUser, nil
}

func (sc *SessionCtx) SesstionCloseById(userId string) error {
	// if !sc.dbc.TypeExists("", &db.Session{}, map[string]interface{}{"id_user": userId}) {
	// 	return nil
	// }

	ses := db.Session{}
	if !sc.dbc.TypeExists("", &ses, map[string]interface{}{"id_user": userId}) {
		return nil
	}

	if _, err := sc.dbc.TypeGet("", &ses, map[string]interface{}{"id_user": userId}); err != nil {
		return err
	}

	if _, err := sc.dbc.TypeDelete("", &ses, nil); err != nil {
		return err
	}

	return nil
}

/*
 * Destroy session
 * Return: error if destroy is failed
 */
func (sc *SessionCtx) SessionDestroy(session string) error {
	_, err := sc.dbc.TypeDelete("", &db.Session{Id: session}, nil)
	return err
}
