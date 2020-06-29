package posts

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"
)

var (
	// ErrNoPostInResponse はレスポンスの川柳が違っていたらエラーを返す
	ErrNoPostInResponse = errors.New("Response has no matching post")
	// ErrMissingField は入力フォームに誤りがあるときに返す
	ErrMissingField = "Error missing %v"
)

// Post は川柳情報
type Post struct {
	PostID string 
	Kamigo string 
	Nakashichi string 
	Shimogo string 
	UserID string
	SignupAt string 
}

// New は川柳データを作成する
func New() Post {
	p := Post{}
	return u
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
