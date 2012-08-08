package session

import (
	"errors"
	"time"
    "crypto/rand"
)

type Session string

var (
	sessions       = make(map[string]map[string]interface{})
	sessionTimeout = time.Hour //default expire session after 1 hour
)

func New() Session {
	return NewTimeout(sessionTimeout)
}

func NewTimeout(timeout time.Duration) Session {
	sid, err := generateSessionId()
	if err != nil {
		panic(err)
	}
	sessions[sid] = make(map[string]interface{})
	t := time.NewTimer(timeout)
	go startSession(t.C, sid)
	return Session(sid)
}

func SessionFromId(sid string) Session {
	return Session(sid)
}

func (s Session) RemoveSession() {
	delete(sessions, string(s))
}

func (s Session) Exists() bool {
	_, ok := sessions[string(s)]
	return ok
}

func (s Session) Put(name string, value interface{}) error {
	sobj, ok := sessions[string(s)]
	if !ok {
		return errors.New("Session does not exist.")
	}

	sobj[name] = value
	return nil
}

func (s Session) Get(name string) (interface{}, error) {
	sobj, ok := sessions[string(s)]
	if !ok {
		return nil, errors.New("Session does not exist.")
	}

	var v interface{}
	v, ok = sobj[name]
	if !ok {
		return nil, errors.New("No such key.")
	}

	return v, nil
}

func startSession(c <-chan time.Time, sid string) {
	<-c
	delete(sessions, sid)
}

func generateSessionId() (string, error) {
    bytes := make([]byte, 32)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }

    return string(bytes), nil
}
