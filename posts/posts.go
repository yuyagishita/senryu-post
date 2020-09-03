package posts

import (
	"errors"
	"fmt"
)

var (
	// ErrNoPostInResponse はレスポンスの川柳が違っていたらエラーを返す
	ErrNoPostInResponse = errors.New("Response has no matching post")
	// ErrMissingField は入力フォームに誤りがあるときに返す
	ErrMissingField = "Error missing %v"
)

// Post は川柳情報
type Post struct {
	PostID     string `json:"postId" bson:"-"`
	Kamigo     string `json:"kamigo" bson:"kamigo"`
	Nakashichi string `json:"nakashichi" bson:"nakashichi"`
	Shimogo    string `json:"shimogo" bson:"shimogo"`
	UserID     string `json:"userId" bson:"-"`
	SignupAt   string `json:"signupAt" bson:"-"`
}

// New は川柳データを作成する
func New() Post {
	p := Post{}
	// p.SignupAt = time.Now()
	return p
}

// News は川柳データの配列を作成する
func News() []Post {
	p := []Post{}
	return p
}

// Validate は入力フォームのバリデーションをする
func (p *Post) Validate() error {
	if p.Kamigo == "" {
		return fmt.Errorf(ErrMissingField, "Kamigo")
	}
	if p.Nakashichi == "" {
		return fmt.Errorf(ErrMissingField, "Nakashichi")
	}
	if p.Shimogo == "" {
		return fmt.Errorf(ErrMissingField, "Shimogo")
	}
	return nil
}
