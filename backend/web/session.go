package web

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type (
	// Session implements.
	Session struct {
		sync.Mutex

		expire time.Time
		data   map[string]interface{}
	}
	// SessionContext interface.
	SessionContext interface {
		RemoteIP() string
		Cookie(key string) (*http.Cookie, error)
		SetCookie(cookie *http.Cookie)
	}
	// SessionManager is holds all sessions in runtime.
	SessionManager struct {
		sync.Mutex

		serial int64
		all    map[string]*Session
	}
)

var (
	// SessionIDKey configures name for session ID in cookie
	SessionIDKey = "session_id"
	// SessionExpire configures how long time the session will lasts.
	SessionExpire = time.Minute * 30
	// Sessions is default singleton instance of SessionManager.
	Sessions = &SessionManager{
		Mutex:  sync.Mutex{},
		serial: 0,
		all:    make(map[string]*Session),
	}
)

// Start a session.
func (mgr *SessionManager) Start(ctx SessionContext) *Session {
	mgr.Lock()
	defer mgr.Unlock()

	cookie, err := ctx.Cookie(SessionIDKey)
	now := time.Now()

	if err == nil {
		s, ok := mgr.all[cookie.Value]
		if ok {
			if s.expire.After(now) {
				s.expire = now.Add(SessionExpire)
				return s
			}

			delete(mgr.all, cookie.Value)
		}
	}

	sn := mgr.serial
	seed := rand.Intn(99999)
	ip := ctx.RemoteIP()
	timestamp := now.Unix()
	code := fmt.Sprintf("%s_%d_%d_%d", ip, timestamp, sn, seed)

	hasher := md5.New()
	hasher.Write([]byte(code))
	id := fmt.Sprintf("%x", hasher.Sum(nil))

	session := &Session{
		Mutex:  sync.Mutex{},
		expire: now.Add(SessionExpire),
		data:   make(map[string]interface{}),
	}

	ctx.SetCookie(&http.Cookie{
		Name:    SessionIDKey,
		Value:   id,
		Expires: session.expire,
		Path:    "/",
	})

	mgr.all[id] = session
	return session
}

// End session.
func (mgr *SessionManager) End(ctx SessionContext) {
	mgr.Lock()
	defer mgr.Unlock()

	cookie, err := ctx.Cookie(SessionIDKey)
	now := time.Now()
	if err == nil {
		delete(mgr.all, cookie.Value)
	}

	ctx.SetCookie(&http.Cookie{
		Name:    SessionIDKey,
		Value:   "",
		Path:    "/",
		MaxAge:  -1,
		Expires: now,
	})
}

// Has tests if there is a data named given key in this session
func (s *Session) Has(key string) bool {
	s.Lock()
	defer s.Unlock()

	_, ok := s.data[key]
	return ok
}

// Get named session value
func (s *Session) Get(key string) interface{} {
	s.Lock()
	defer s.Unlock()

	return s.data[key]
}

// Set named value into session.
func (s *Session) Set(key string, val interface{}) {
	s.Lock()
	defer s.Unlock()

	s.data[key] = val
}
