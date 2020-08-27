package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/yu-yagishita/senryu-post/posts"
)

// MakeGetAllEndpoint は全データを取得する
func MakeGetAllEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		p, err := svc.GetAll()
		return getAllResponse{Post: p}, err
	}
}

// MakePostEndpoint はデータを投稿する
func MakePostEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postRequest)
		fmt.Println("req.Kamigo: " + req.Kamigo)
		fmt.Println("req.Nakashichi: " + req.Nakashichi)
		fmt.Println("req.Nakashichi: " + req.Shimogo)
		fmt.Println("req.UserID: " + req.UserID)
		id, err := svc.Register(req.Kamigo, req.Nakashichi, req.Shimogo, req.UserID)
		return postResponse{ID: id}, err
	}
}

// DecodeGetAllRequest はregisterのリクエストをデコードする
func DecodeGetAllRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request getAllRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

// DecodePostRequest はpostのリクエストをデコードする
func DecodePostRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request postRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

// EncodeResponse はレスポンスをエンコードする
func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

// EncodeRequest はレスポンスをエンコードする
func EncodeRequest(_ context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

type getAllRequest struct {
}

type getAllResponse struct {
	Post []posts.Post `json:"posts"`
}

type postRequest struct {
	Kamigo     string `json:"kamigo"`
	Nakashichi string `json:"nakashichi"`
	Shimogo    string `json:"shimogo"`
	UserID     string `json:"user_id"`
}

type postResponse struct {
	ID string `json:"id"`
}
