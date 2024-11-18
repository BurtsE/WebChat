package main

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
)

// getProviderName is a custom function for gothic to use
func getProviderName(req *http.Request) (string, error) {
	providerName := req.PathValue("provider")
	return providerName, nil
}

// getUserData retrieves cookies
// and decodes them from base64
func getUserData(r *http.Request) (map[string]interface{}, error) {
	var (
		userData    = map[string]interface{}{}
		authCookie  *http.Cookie
		cookieValue []byte
		err         error
	)
	if authCookie, err = r.Cookie("auth"); err != nil {
		return nil, err
	}
	if cookieValue, err = base64.StdEncoding.DecodeString(authCookie.Value); err != nil {
		return nil, err
	}

	err = json.Unmarshal(cookieValue, &userData)
	if err != nil {
		return nil, err
	}
	return userData, nil
}
