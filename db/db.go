package db

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/yu-yagishita/senryu-post/posts"
)

// Database 新しいシステムに簡単に切り替えることができるようにシンプルなインターフェースにしている
type Database interface {
	Init() error
	GetAll() (posts.Post, error)
	CreatePost(*posts.Post) error
}

var (
	database string
	// DefaultDb はマイクロサービスのデータベースセット
	DefaultDb Database
	// DBTypes はこのサービスで使用できるDBインターフェースのマップ
	DBTypes = map[string]Database{}
	// ErrNoDatabaseFound はDBTypesにデータベースインターフェースが存在しない場合にエラーを返す
	ErrNoDatabaseFound = "No database with name %v registered"
	// ErrNoDatabaseSelected はflagまたはenvにデータベースが指定されていない場合に返す
	ErrNoDatabaseSelected = errors.New("No DB selected")
)

func init() {
	flag.StringVar(&database, "database", os.Getenv("USER_DATABASE"), "Database to use, Mongodb or ...")
}

// Init は DefaultDb で選択した DB を起動する
func Init() error {
	if database == "" {
		return ErrNoDatabaseSelected
	}
	err := Set()
	if err != nil {
		return err
	}
	return DefaultDb.Init()
}

// Set はDefaultDbを設定する
func Set() error {
	if v, ok := DBTypes[database]; ok {
		DefaultDb = v
		return nil
	}
	return fmt.Errorf(ErrNoDatabaseFound, database)
}

// Register DBTypesにデータベースインターフェイスを登録する
func Register(name string, db Database) {
	DBTypes[name] = db
}

// GetAll はDefaultDbメソッドを呼び出す
func GetAll() (posts.Post, error) {
	p, err := DefaultDb.GetAll()
	if err == nil {
		// u.AddLinks()
	}
	return p, err
}

// CreatePost はDefaultDbメソッドを呼び出す
func CreatePost(p *posts.Post) error {
	return DefaultDb.CreatePost(p)
}
