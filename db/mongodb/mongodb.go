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
func (m *Mongo) GetAll() ([]posts.Post, error) {
	s := m.Session.Copy()
	defer s.Close()
	c := s.DB("").C("posts")
	// mp := New()
	var mps []MongoPost
	var pps []posts.Post
	err := c.Find(nil).All(&mps)
	for i := 0; i < len(mps); i++ {
		mps[i].Post.PostID = mps[i].ID.Hex()
		mps[i].Post.UserID = mps[i].UserID.Hex()
		mps[i].Post.SignupAt = mps[i].SignupAt.Format("2006-01-02")
		pps = append(pps, mps[i].Post)
	}
	fmt.Println(pps)

	return pps, err
}

// Get はmongoのuserIDに紐づく川柳データを取得する
func (m *Mongo) Get(userID string) ([]posts.Post, error) {
	s := m.Session.Copy()
	defer s.Close()
	c := s.DB("").C("posts")
	var mps []MongoPost
	var pps []posts.Post
	err := c.Find(bson.M{"user_id": bson.ObjectIdHex(userID)}).All(&mps)
	for i := 0; i < len(mps); i++ {
		mps[i].Post.PostID = mps[i].ID.Hex()
		mps[i].Post.UserID = mps[i].UserID.Hex()
		mps[i].Post.SignupAt = mps[i].SignupAt.Format("2006-01-02")
		pps = append(pps, mps[i].Post)
	}
	fmt.Println(pps)

	return pps, err
}

// CreatePost はユーザーを作成してMongoに保存する
func (m *Mongo) CreatePost(p *posts.Post) error {
	s := m.Session.Copy()
	defer s.Close()
	id := bson.NewObjectId()
	mp := New()
	mp.Post = *p
	mp.ID = id
	mp.Kamigo = p.Kamigo
	mp.Nakashichi = p.Nakashichi
	mp.Shimogo = p.Shimogo
	mp.UserID = bson.ObjectIdHex(p.UserID)
	mp.SignupAt = time.Now()
	fmt.Println(mp)
	c := s.DB("").C("posts")
	err := c.Insert(mp)
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
	mp.Post.PostID = mp.ID.Hex()
	*p = mp.Post
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
