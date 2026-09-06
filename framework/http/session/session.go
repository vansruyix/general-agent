package session

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"sync"
	"time"
)

const (
	CK_NAME = "sid"
	CK_AGE  = 3600
)

var globalStore Storage

// Session concurrent safe mutex
var migrateMux sync.Mutex

// SessionLoader x
type SessionLoader interface {
	Set(key string, value interface{})
	Get(key string) interface{}
	Remove(key string) error
	GetId() string
}

type Session struct {
	session
}

// Values is session item value
type Values map[string]interface{}

// session struct
type session struct {
	id         string
	CreateTime time.Time
	ExpireTime time.Time
	Values
}

// GetSession Get session data from the Request
func GetSession(w http.ResponseWriter, req *http.Request) (*Session, error) {
	var session Session

	cookie, err := req.Cookie(CK_NAME)
	if cookie == nil || err != nil {
		return createSession(w, cookie)
	}

	if len(cookie.Value) >= 73 {
		session.id = cookie.Value
		if globalStore.Read(&session) != nil {
			return createSession(w, cookie)
		}
	}

	return &session, nil
}

// ID return session id
func (s *Session) ID() string {
	return s.id
}

// Sync save data modify
func (s *Session) Sync() error {
	return globalStore.Write(s)
}

// Migrate migrate old session data to new session
func Migrate(write http.ResponseWriter, old *Session) (*Session, error) {
	var (
		ns     = NewSession()
		cookie = NewCookie()
	)

	migrateMux.Lock()
	ns.Values = old.Values
	cookie.Value = ns.id
	cookie.MaxAge = CK_AGE
	migrateMux.Unlock()

	return ns,
		func() error {
			if ns.Sync() != nil {
				return ErrMigrateSessionFail
			}
			if globalStore.Remove(old) != nil {
				return ErrRemoveSessionFail
			}
			http.SetCookie(write, cookie)
			return nil
		}()
}

// createSession return new session
func createSession(w http.ResponseWriter, cookie *http.Cookie) (*Session, error) {

	// FIX BUG:
	// https://deepsource.io/gh/auula/gws/run/5b13c99b-9101-4e4f-8197-acfd730c28a0/go/SCC-SA4009
	session := NewSession()

	if cookie == nil {
		cookie = NewCookie()
	}
	cookie.Value = session.id
	cookie.MaxAge = CK_AGE
	if err := globalStore.Write(session); err != nil {
		return nil, err
	}

	http.SetCookie(w, cookie)

	return session, nil
}

// NewCookie return default model cookie pointer
func NewCookie() *http.Cookie {
	return &http.Cookie{
		Path:     "/",
		Name:     CK_NAME,
		HttpOnly: false,
	}
}

// uuid73 generate session uuid length 73
func uuid73() string {
	return fmt.Sprintf("%s-%s", uuid.New().String(), uuid.New().String())
}

// NewSession return new session
func NewSession() *Session {
	nowTime := time.Now()
	return &Session{
		session: session{
			id:         uuid73(),
			Values:     make(Values),
			CreateTime: nowTime,
			ExpireTime: nowTime.Add(time.Duration(CK_AGE) * time.Second),
		},
	}
}

// Expired check current session whether expire
func (s *Session) Expired() bool {
	return time.Duration(s.ExpireTime.UnixNano()) <= time.Duration(time.Now().UnixNano())
}

// Invalidate remove the session
func Invalidate(s *Session) error {
	return globalStore.Remove(s)
}

// Malloc reallocation of memory
func Malloc(v *Values) {
	*v = make(Values)
}
