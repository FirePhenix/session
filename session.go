package session

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"time"
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
	sid, err := newGuid()
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

func newGuid() (string, error) {
	f, err := os.Open("/dev/urandom")
	if err != nil {
		if runtime.GOOS == "windows" {
			//no /dev/urandom on windows
			//TODO: use some windows api instead ?
			guid := getMd5Hex(time.Now().UTC().Format(time.ANSIC) + strconv.FormatInt(rand.Int63(), 10))
			return guid, nil
		}
		return "", err
	}
	defer f.Close()

	b := make([]byte, 16)
	_, err = f.Read(b)
	if err != nil {
		return "", err
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid, nil
}

func getMd5Hex(input string) string {
	hasher := md5.New()
	hasher.Write([]byte(input))
	return hex.EncodeToString(hasher.Sum(nil))
}
