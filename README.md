# senryu-post

## 概要
川柳を共有することができるSNSアプリの川柳データを管理しているバックエンド

## 技術
|名前|備考|
|--|--|
|Go kit(Golang)|バックエンド|
|MongoDB|データベース|
|Docker|開発環境|

## Run
### Docker Compose
ルートディレクトリで`docker-compose up`  

## Use
```
// 全取得
curl -XPOST -d'{}' localhost:8081/get-all

// ユーザーに紐づくデータ取得
curl -XPOST -d'{}' localhost:8081/get/57a98d98e4b00679b4a830af

// 投稿
curl -XPOST -d'{"kamigo": "テストだよ", "nakashichi": "この投稿は",  "shimogo": "テストだよ", "userId": "57a98d98e4b00679b4a830af"}' localhost:8081/post
```

## 開発手法（VSCode Remote Container）
`docker-compose.yml`の`command: fresh -c .fresh.conf`をコメントアウトして、
VSCode Remote Containerで実行する。  
コンテナ内で`fresh -c .fresh.conf`を実行するとホットリロードで開発ができる。
