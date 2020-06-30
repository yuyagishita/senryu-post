package mongodb

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/yu-yagishita/senryu-post/posts"
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
	u := getURL()
	var err error
	m.Session, err = mgo.DialWithTimeout(u.String(), time.Duration(5)*time.Second)
	if err != nil {
		return err
	}
	return err
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

// MongoPost はPostのラッパー
type MongoPost struct {
	posts.Post `bson:",inline"`
	ID         bson.ObjectId `bson:"_id"`
	UserID     bson.ObjectId `bson:"user_id"`
}

// NewPost 新しいMongoPostを返す
func NewPost() MongoPost {
	p := posts.New()
	return MongoPost{
		Post: p,
	}
}

// GetAll はmongoの全川柳データを取得する
func (m *Mongo) GetAll() (posts.Post, error) {
	s := m.Session.Copy()
	defer s.Close()
	c := s.DB("").C("posts")
	mu := NewPost()
	err := c.Find(nil).One(&mu)
	mu.Post.PostID = mu.ID.Hex()
	mu.Post.UserID = mu.UserID.Hex()
	fmt.Println(mu.Post)
	return mu.Post, err
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
