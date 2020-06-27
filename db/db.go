package db

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/yu-yagishita/senryu-post/users"
)

// Database 新しいシステムに簡単に切り替えることができるようにシンプルなインターフェースにしている
type Database interface {
	Init() error
	GetUserByName(string) (users.User, error)
	CreateUser(*users.User) error
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
	fmt.Println("db: init")
	flag.StringVar(&database, "database", os.Getenv("USER_DATABASE"), "Database to use, Mongodb or ...")

}

// Init は DefaultDb で選択した DB を起動する
func Init() error {
	fmt.Println("db: Init")
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
	fmt.Println("Mongo: Register")
	DBTypes[name] = db
}

// CreateUser はDefaultDbメソッドを呼び出す
func CreateUser(u *users.User) error {
	return DefaultDb.CreateUser(u)
}

// GetUserByName はDefaultDbメソッドを呼び出す
func GetUserByName(n string) (users.User, error) {
	fmt.Println("start GetUserByName")
	u, err := DefaultDb.GetUserByName(n)
	if err == nil {
		// u.AddLinks()
	}
	fmt.Println("end GetUserByName")
	return u, err
}
