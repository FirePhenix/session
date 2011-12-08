package session

import (
    "time"
    "encoding/hex"
    "crypto/md5"
    "os"
    "syscall"
    "fmt"
    "rand"
    "strconv"
)

type Session string

var (
    sessions       = make(map[string]map[string]interface{})
    sessionTimeout = int64(3600) //By default expire session after 1 hour
)

func New() Session {
    return NewTimeout(sessionTimeout)
}

func NewTimeout(timeout int64) Session {
    sid, err := newGuid()
    if err != nil {
        panic(err)
    }
    sessions[sid] = make(map[string]interface{})
    t := time.NewTimer(timeout * 1e9)
    go startSession(t.C, sid)
    return Session(sid)
}

func SessionFromId(sid string) Session {
    return Session(sid)
}

func (s Session) RemoveSession() {
    sessions[string(s)] = nil, false
}

func (s Session) Exists() bool {
    _, ok := sessions[string(s)]
    return ok
}

func (s Session) Put(name string, value interface{}) os.Error {
    sobj, ok := sessions[string(s)]
    if !ok {
        return os.NewError("Session does not exist.")
    }

    sobj[name] = value
    return nil
}

func (s Session) Get(name string) (interface{}, os.Error) {
    sobj, ok := sessions[string(s)]
    if !ok {
        return nil, os.NewError("Session does not exist.")
    }

    var v interface{}
    v, ok = sobj[name]
    if !ok {
        return nil, os.NewError("No such key.")
    }

    return v, nil
}

func startSession(c <-chan int64, sid string) {
    <-c
    sessions[sid] = nil, false
}

func newGuid() (string, os.Error) {
    if syscall.OS == "windows" {
        //no /dev/urandom on windows
        //TODO: use some windows api instead ?
        guid := getMd5Hex(time.UTC().Format(time.ANSIC) + strconv.Itoa64(rand.Int63()))
        return guid, nil
    }

    f, err := os.Open("/dev/urandom")
    if err != nil {
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
    return hex.EncodeToString(hasher.Sum())
}
