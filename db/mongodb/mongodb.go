package mongodb

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/yu-yagishita/senryu-post/users"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	name     string
	password string
	host     string
	db       = "posts"
	// ErrInvalidHexID は有効なbsonObjectIDではないエンティティIDを表している
	ErrInvalidHexID = errors.New("Invalid Id Hex")
)

func init() {
	fmt.Println("name:" + os.Getenv("MONGO_USER"))
	fmt.Println("password:" + os.Getenv("MONGO_PASS"))
	fmt.Println("host:" + os.Getenv("MONGO_HOST"))
	flag.StringVar(&name, "mongo-user", os.Getenv("MONGO_USER"), "Mongo user")
	flag.StringVar(&password, "mongo-password", os.Getenv("MONGO_PASS"), "Mongo password")
	flag.StringVar(&host, "mongo-host", os.Getenv("MONGO_HOST"), "Mongo host")
}

// Mongo はデータベースインターフェイスの要件を満たしている
type Mongo struct {
	//Session はMongoDBのセッション
	Session *mgo.Session
}

// Init はMongoDBのInit処理
func (m *Mongo) Init() error {
	fmt.Println("MongoDB: Init")
	u := getURL()
	fmt.Println("u: " + u.String())
	var err error
	m.Session, err = mgo.DialWithTimeout(u.String(), time.Duration(5)*time.Second)
	if err != nil {
		return err
	}
	return m.EnsureIndexes()
}

// MongoUser はUserのラッパー
type MongoUser struct {
	users.User `bson:",inline"`
	ID         bson.ObjectId `bson:"_id"`
}

// New 新しいMongoUserを返す
func New() MongoUser {
	u := users.New()
	return MongoUser{
		User: u,
	}
}

// CreateUser はユーザーを作成してMongoに保存する
func (m *Mongo) CreateUser(u *users.User) error {
	s := m.Session.Copy()
	defer s.Close()
	id := bson.NewObjectId()
	mu := New()
	mu.User = *u
	mu.ID = id
	c := s.DB("").C("users")
	_, err := c.UpsertId(mu.ID, mu)
	if err != nil {
		return err
	}
	mu.User.UserID = mu.ID.Hex()
	*u = mu.User
	return nil
}

// GetUserByName はusernameでmongoのユーザーデータを取得する
func (m *Mongo) GetUserByName(name string) (users.User, error) {
	fmt.Println("mongodb: GetUserByName")
	s := m.Session.Copy()
	defer s.Close()
	c := s.DB("").C("users")
	fmt.Println("c.Name: " + c.Name)
	mu := New()
	err := c.Find(bson.M{"username": name}).One(&mu)
	return mu.User, err
}

func getURL() url.URL {
	ur := url.URL{
		Scheme: "mongodb",
		Host:   host,
		Path:   db,
	}
	if name != "" {
		u := url.UserPassword(name, password)
		ur.User = u
	}
	return ur
}

// EnsureIndexes はusernameが一意であることを確認する
func (m *Mongo) EnsureIndexes() error {
	s := m.Session.Copy()
	defer s.Close()
	i := mgo.Index{
		Key:        []string{"username"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     false,
	}
	c := s.DB("").C("users")
	return c.EnsureIndex(i)
}
