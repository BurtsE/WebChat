package main

import (
	"errors"
	"testing"
)

func TestAutahAvatar(t *testing.T) {
	var authAvatar AuthAvatar
	testClient := new(client)
	url, err := authAvatar.GetAvatarURL(testClient) // Calling a method on a nil object
	if !errors.Is(err, ErrNoAvatarURL) {
		t.Error("AuthAvatar.GetAvatarURL should return ErrNoAvatarURL when no value present")
	}
	// set a value
	testUrl := "http://url-to-gravatar/"
	testClient.userData = map[string]interface{}{"avatar_url": testUrl}
	url, err = authAvatar.GetAvatarURL(testClient)
	if err != nil {
		t.Error("AuthAvatar.GetAvatarURL should return no error when value present")
	}
	if url != testUrl {
		t.Error("AuthAvatar.GetAvatarURL should return correct URL")
	}
}

func TestGravatarAvatar(t *testing.T) {
	var gravatarAvatar GravatarAvatar
	testClient := new(client)
	testClient.userData = map[string]interface{}{"userId": "0bc83cb571cd1c50ba6f3e8a78ef1346"}

	url, err := gravatarAvatar.GetAvatarURL(testClient)
	if err != nil {
		t.Error("GravatarAvatar.GetAvatarURL should not return an error")
	}
	if url != "//www.gravatar.com/avatar/0bc83cb571cd1c50ba6f3e8a78ef1346" {
		t.Errorf("GravatarAvatar.GetAvatarURL wrongly returned %s", url)
	}

}
