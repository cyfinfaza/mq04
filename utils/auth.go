package utils

import "github.com/tg123/go-htpasswd"

type FileAuth struct {
	Checker *htpasswd.File
}

func (a *FileAuth) Authenticate(username []byte, password []byte) bool {
	return a.Checker.Match(string(username), string(password))
}

func (a *FileAuth) ACL(user []byte, topic string, write bool) bool {
	return true
}
