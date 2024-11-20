package main

import (
	"errors"
	"fmt"
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

var useAuthAvatar AuthAvatar

func (AuthAvatar) GetAvatarURL(c *client) (string, error) {
	var (
		url    interface{}
		urlStr string
		ok     bool
	)
	if url, ok = c.userData["avatar_url"]; !ok {
		return "", ErrNoAvatarURL
	}
	if urlStr, ok = url.(string); !ok {
		return "", ErrNoAvatarURL
	}

	return urlStr, nil
}

type GravatarAvatar struct{}

var useGravatarAvatar GravatarAvatar

func (GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	var (
		userId   interface{}
		userHash string
		ok       bool
	)
	if userId, ok = c.userData["userId"]; !ok {
		return "", ErrNoAvatarURL
	}
	if userHash, ok = userId.(string); !ok {
		return "", ErrNoAvatarURL
	}

	return fmt.Sprintf("//www.gravatar.com/avatar/%s", userHash), nil
}
