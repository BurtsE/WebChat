package main

import (
	"errors"
	"os"
	"path"
	"path/filepath"
)

// ErrNoAvatarURL is the error that is returned when the
// Avatar instance is unable to provide an avatar URL.
var ErrNoAvatarURL = errors.New("chat: Unable to get an avatar URL")

// Avatar represents types capable of representing
// user profile pictures.
type Avatar interface {
	// GetAvatarURL gets the avatar URL for the specified client,
	// or returns an error if something goes wrong.
	GetAvatarURL(c *client) (string, error)
}

type AuthAvatar struct{}

var UseAuthAvatar AuthAvatar

func (AuthAvatar) GetAvatarURL(user chatUser) (string, error) {
	url := user.AvatarURL
	if len(url) == 0 {
		return "", ErrNoAvatarURL
	}
	return url, nil
}

type GravatarAvatar struct{}

var UseGravatarAvatar GravatarAvatar

func (GravatarAvatar) GetAvatarURL(user chatUser) (string, error) {
	return "//www.gravatar.com/avatar/" + user.UniqueID(), nil
}

type FileSystemAvatar struct{}

var UseFileSystemAvatar FileSystemAvatar

func (FileSystemAvatar) GetAvatarURL(user ChatUser) (string, error) {
	dirname := filepath.Join("users", "avatars")
	files, err := os.ReadDir(dirname)
	if err != nil {
		return "", ErrNoAvatarURL
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if match, _ := path.Match(user.UniqueID()+"*", file.Name()); match {
			return filepath.Join(dirname, file.Name()), nil
		}
	}
	return "", ErrNoAvatarURL
}
