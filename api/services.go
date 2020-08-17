package api

import (
	"errors"

	"github.com/yu-yagishita/senryu-post/db"
	"github.com/yu-yagishita/senryu-post/posts"
)

var (
	// ErrUnauthorized エラーになった際に返す値
	ErrUnauthorized = errors.New("Unauthorized")
)

// Service はAPIで提供している機能をまとめたインターフェース
type Service interface {
	GetAll() (posts.Post, error)
	Register(kamigo, nakashichi, shimogo, userID string) (string, error)
}

// NewFixedService はfixedService{}を返す
func NewFixedService() Service {
	return &fixedService{}
}

type fixedService struct{}

func (s *fixedService) GetAll() (posts.Post, error) {
	p, err := db.GetAll()
	if err != nil {
		return posts.New(), err
	}

	return p, nil
}

func (s *fixedService) Register(kamigo, nakashichi, shimogo, userID string) (string, error) {
	p := posts.New()
	p.Kamigo = kamigo
	p.Nakashichi = nakashichi
	p.Shimogo = shimogo
	p.UserID = userID
	err := db.CreatePost(&p)
	return p.PostID, err
}

// ErrEmpty は入力文字列が空の場合返す
var ErrEmpty = errors.New("empty string")

// ServiceMiddleware はServiceのチェイン可能な動作修飾子
type ServiceMiddleware func(Service) Service
