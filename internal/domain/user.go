package domain

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
