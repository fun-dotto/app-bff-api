package domain

import "errors"

// ErrUserNotFound ユーザーが見つからない
var ErrUserNotFound = errors.New("user not found")

// User ユーザー
type User struct {
	ID     string
	Email  string
	Grade  *Grade
	Course *Course
	Class  *Class
}

// UserRequest ユーザー作成・更新リクエスト
type UserRequest struct {
	Email  string
	Grade  *Grade
	Course *Course
	Class  *Class
}
