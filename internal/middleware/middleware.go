package middleware

import (
	"crypto/subtle"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type BasicAuth struct {
	Username []byte
	Password []byte
}

func (b *BasicAuth) SetCredentials(username string, password string) {
	b.Username = []byte(username)
	b.Password = []byte(password)
}

func (b *BasicAuth) Middleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		presentedUsername, presentedPassword, ok := r.BasicAuth()
		if ok {

			usernameMatch := subtle.ConstantTimeCompare([]byte(presentedUsername)[:], b.Username[:]) == 1
			passwordMatch := bcrypt.CompareHashAndPassword(b.Password, []byte(presentedPassword)[:])

			if usernameMatch == true && passwordMatch == nil {
				next.ServeHTTP(w, r)
				return
			}
		}
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}
