session
========
session is a simple session management module for web applications written in Go programming language. 

Usage
=======
A session object is just a string represents the session identifier (session id):

    type Session string

Create a session with default timeout (1 hour):

    s := session.New() 

Create a session with specified timeout( in seconds):

    s := session.NewTimeout(20 * time.Minute)  // 20 minutes

Check if a session exists :

    ok := s.Exists()

Put something into a session:

    err := s.Put("Username", "Bob")
    err  = s.Put("Logintime", time.UTC())

Retrieve value from a session:

    v, err :=  s.Get("Logintime") //return interface{}, need type assertion 
    loginTime := v.(*time.Time)

Remove a session manually:
    
    s.RemoveSession()

Use session.go with web.go:

    s := session.New()
    ctx.SetCookie("u", string(s), 3600)  //set session-cookie
    err := s.Put("user", u)

Installation
========
session.go compiles with Go latest release:

    go get github.com/cranej/session
will download, compile and install the package.  

Then use it by `import "github.com/cranej/session"`
