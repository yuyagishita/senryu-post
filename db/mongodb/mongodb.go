package mongodb

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/yu-yagishita/senryu-post/posts"

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

// MongoPost はPostのラッパー
type MongoPost struct {
	posts.Post `bson:",inline"`
	ID         bson.ObjectId `bson:"_id"`
	UserID     bson.ObjectId `bson:"user_id"`
	SignupAt   time.Time     `bson:"signup_at"`
}

// New 新しいMongoPostを返す
func New() MongoPost {
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
	mu := New()
	err := c.Find(nil).One(&mu)
	mu.Post.PostID = mu.ID.Hex()
	mu.Post.UserID = mu.UserID.Hex()
	mu.Post.SignupAt = mu.SignupAt.Format("2006-01-02")
	return mu.Post, err
}

// CreatePost はユーザーを作成してMongoに保存する
func (m *Mongo) CreatePost(p *posts.Post) error {
	s := m.Session.Copy()
	defer s.Close()
	id := bson.NewObjectId()
	mu := New()
	fmt.Print(p)
	// mu.Post = *p
	mu.ID = id
	mu.Kamigo = p.Kamigo
	mu.Nakashichi = p.Nakashichi
	mu.Shimogo = p.Shimogo
	mu.UserID = bson.ObjectIdHex(p.UserID)
	mu.SignupAt = time.Now()
	c := s.DB("").C("posts")
	err := c.Insert(mu)
	if err != nil {
		if mgo.IsDup(err) {
			fmt.Printf("Duplicate key error \n")
		}
		if v, ok := err.(*mgo.LastError); ok {
			fmt.Printf("Code:%d N:%d UpdatedExisting:%t WTimeout:%t Waited:%d \n", v.Code, v.N, v.UpdatedExisting, v.WTimeout, v.Waited)
		} else {
			fmt.Printf("%+v \n", err)
		}
	}
	mu.Post.PostID = mu.ID.Hex()
	*p = mu.Post
	return nil
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
