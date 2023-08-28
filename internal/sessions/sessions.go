package sessions

import (
	"encoding/gob"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

// User holds a users account information
type User struct {
	Authenticated    bool
	AccountConnected bool
}

// store will hold all session data
var store *sessions.CookieStore

func Init() {
	authKeyOne := securecookie.GenerateRandomKey(64)
	encryptionKeyOne := securecookie.GenerateRandomKey(32)

	store = sessions.NewCookieStore(
		authKeyOne,
		encryptionKeyOne,
	)

	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   0,
		HttpOnly: true,
	}

	gob.Register(User{})
}

	session, err := store.Get(r, "cookie-name")
	if err != nil {
	}

	if r.FormValue("code") != "code" {
		if r.FormValue("code") == "" {
			session.AddFlash("Must enter a code")
		}
		session.AddFlash("The code was incorrect")
		err = session.Save(r, w)
		if err != nil {
		}

	}

	username := r.FormValue("username")

	user := &User{
		Username:      username,
		Authenticated: true,
	}

	session.Values["user"] = user

	err = session.Save(r, w)
