package main

import (
	"errors"
	"fmt"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/stretchr/objx"
	"net/http"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("auth")
	if errors.Is(err, http.ErrNoCookie) || cookie.Value == "" {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.next.ServeHTTP(w, r)
}

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

// loginHandler handles the third-party login process.
// format: /auth/{action}/{provider}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	action := r.PathValue("action")
	//providerName := r.PathValue("provider")
	switch action {
	case "login":
		gothic.BeginAuthHandler(w, r)
	case "callback":
		callback(w, r)
		return
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Auth action %s not supported", action)
	}
}

func callback(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error when trying to get user from %s: %s",
			user, err), http.StatusInternalServerError)
		return
	}
	authCookieValue := objx.New(map[string]interface{}{
		"name":       getUserName(user),
		"avatar_url": user.AvatarURL,
		"email":      user.Email,
	}).MustBase64()
	http.SetCookie(w, &http.Cookie{
		Name:  "auth",
		Value: authCookieValue,
		Path:  "/"})
	w.Header().Set("Location", "/chat")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// TODO refactor to return string
func getUserName(user goth.User) string {
	switch {
	case user.Name != "":
		return user.Name
	case user.NickName != "":
		return user.NickName
	case user.Email != "":
		return user.Email
	default:
		return "unknown"
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "auth",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	w.Header().Set("Location", "/chat")
	w.WriteHeader(http.StatusTemporaryRedirect)
}
