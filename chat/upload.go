package main

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path"
)

func uploaderHandler(w http.ResponseWriter, req *http.Request) {
	userId := req.FormValue("userId")
	file, header, err := req.FormFile("avatarFile")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = os.Mkdir(path.Join("users", "avatars", userId), 0777)
	if err != nil && !errors.Is(err, os.ErrExist) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	filename := path.Join("users", "avatars", userId+path.Ext(header.Filename))
	err = os.WriteFile(filename, data, 0777)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = io.WriteString(w, "Successful")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
