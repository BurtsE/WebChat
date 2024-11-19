package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/markbates/goth/gothic"
	"net/http"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("auth")
	if err == http.ErrNoCookie {
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

	var userName interface{}
	provider, _ := getProviderName(r)
	switch provider {
	case "github":
		userName = user.RawData["login"]
	case "google":
		userName = user.RawData["email"]
	}
	if userName == nil {
		userName = "unknown"
	}

	authCookie, _ := json.Marshal(map[string]interface{}{
		"name": userName,
	})
	authCookieValue := base64.StdEncoding.EncodeToString(authCookie)
	http.SetCookie(w, &http.Cookie{
		Name:  "auth",
		Value: authCookieValue,
		Path:  "/"})
	w.Header().Set("Location", "/chat")
	w.WriteHeader(http.StatusTemporaryRedirect)
}
