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
	"github.com/yu-yagishita/senryu-post/users"
)

// MakeUppercaseEndpoint は渡された文字を大文字に変換して返す
func MakeUppercaseEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(uppercaseRequest)
		v, err := svc.Uppercase(req.S)
		if err != nil {
			return uppercaseResponse{v, err.Error()}, nil
		}
		return uppercaseResponse{v, ""}, nil
	}
}

// MakeCountEndpoint は渡された文字から文字数をカウントして返す
func MakeCountEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(countRequest)
		v := svc.Count(req.S)
		return countResponse{v}, nil
	}
}

// MakeLoginEndpoint はログインチェックをしてDBに該当データがある場合、ユーザー情報を返す
func MakeLoginEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(loginRequest)
		fmt.Println("req.Username: " + req.Username)
		fmt.Println("req.Password: " + req.Password)
		u, err := svc.Login(req.Username, req.Password)
		return userResponse{User: u}, err
	}
}

// MakeRegisterEndpoint は新規ユーザーを登録する
func MakeRegisterEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(registerRequest)
		fmt.Println("req.Username: " + req.Username)
		fmt.Println("req.Email: " + req.Email)
		fmt.Println("req.Password: " + req.Password)
		id, err := svc.Register(req.Username, req.Email, req.Password)
		return postResponse{ID: id}, err
	}
}

// MakeGetAllEndpoint は新規ユーザーを登録する
func MakeGetAllEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		fmt.Println("aaaaaaaaaaaaaaaaa")
		// req := request.(getAllRequest)
		p, err := svc.GetAll()
		return getAllResponse{Post: p}, err
	}
}

// DecodeUppercaseRequest はリクエストをデコードする
func DecodeUppercaseRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request uppercaseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

// DecodeCountRequest はリクエストをデコードする
func DecodeCountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request countRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

// DecodeLoginRequest はリクエストをデコードする
func DecodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request loginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

// DecodeRegisterRequest はregisterのリクエストをデコードする
func DecodeRegisterRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request registerRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

// DecodeGetAllRequest はregisterのリクエストをデコードする
func DecodeGetAllRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request getAllRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

// DecodeUppercaseResponse はレスポンスをデコードする
func DecodeUppercaseResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var response uppercaseResponse
	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
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

type uppercaseRequest struct {
	S string `json:"s"`
}

type uppercaseResponse struct {
	V   string `json:"v"`
	Err string `json:"err,omitempty"`
}

type countRequest struct {
	S string `json:"s"`
}

type countResponse struct {
	V int `json:"v"`
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type userResponse struct {
	User users.User `json:"user"`
}

type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type postResponse struct {
	ID string `json:"id"`
}

type getAllRequest struct {
}

type getAllResponse struct {
	Post posts.Post `json:"post"`
}
