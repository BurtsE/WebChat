package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"strings"
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

// TODO prevent hashing every time we need avatar URL
func (GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	var (
		email    interface{}
		emailStr string
		ok       bool
	)
	if email, ok = c.userData["email"]; !ok {
		return "", ErrNoAvatarURL
	}
	if emailStr, ok = email.(string); !ok {
		return "", ErrNoAvatarURL
	}

	m := md5.New()
	_, err := io.WriteString(m, strings.ToLower(emailStr))
	if err != nil {
		return "", ErrNoAvatarURL
	}
	return fmt.Sprintf("//www.gravatar.com/avatar/%x", m.Sum(nil)), nil
}
