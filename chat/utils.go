package main

import (
	"net/http"
)

// getProviderName is a custom function for gothic to use
func getProviderName(req *http.Request) (string, error) {
	providerName := req.PathValue("provider")
	return providerName, nil
}
