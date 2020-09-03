# senryu-post
ユーザー関連のサービス

## ホットリロード
freshを使ってホットリロードができるようにしている。
`fresh -c .fresh.conf`でGoを実行できる。

## go-kitのサンプル
https://github.com/go-kit/kit/tree/master/examples  
を参考にして最初の学習を開始した。
そのときに、`stringsvc3`でロードバランスを試すときに、自分の環境にfishを使っていてたので、
サンプルでのコマンドと少し違っている。
以下に記載する。

１　APIを立てる
```
$ go get github.com/go-kit/kit/examples/stringsvc3

$ stringsvc3 -listen=:8001 &
env listen=:8001 caller=proxying.go:25 proxy_to=none
env listen=:8001 caller=main.go:72 msg=HTTP addr=:8001


$ stringsvc3 -listen=:8002 &
env listen=:8002 caller=proxying.go:25 proxy_to=none
env listen=:8002 caller=main.go:72 msg=HTTP addr=:8002

$ stringsvc3 -listen=:8080 -proxy=localhost:8001,localhost:8002,localhost:8003
env listen=:8080 caller=proxying.go:29 proxy_to="[localhost:8001 localhost:8002 localhost:8003]"
env listen=:8080 caller=main.go:72 msg=HTTP addr=:8080
```

２　テストする
```
for s in foo bar baz ; curl -d"{\"s\":\"$s\"}" localhost:8080/uppercase; end
{"v":"FOO"}
{"v":"BAR"}
{"v":"BAZ"}
```

サーバのログが以下のようになっていれば成功
```
listen=:8002 caller=logging.go:22 method=uppercase input=foo output=FOO err=null took=11.8µs
listen=:8080 caller=logging.go:22 method=uppercase input=foo output=FOO err=null took=1.8198ms
listen=:8001 caller=logging.go:22 method=uppercase input=bar output=BAR err=null took=8.4µs
listen=:8080 caller=logging.go:22 method=uppercase input=bar output=BAR err=null took=2.8662ms
listen=:8002 caller=logging.go:22 method=uppercase input=baz output=BAZ err=null took=8.3µs
listen=:8080 caller=logging.go:22 method=uppercase input=baz output=BAZ err=null took=1.156ms
```

## 作成したAPIのテスト
ローカル環境で`curl`でテストの仕方を以下に記載する。
```
// 全取得
curl -XPOST -d'{}' localhost:8081/get-all

// ユーザーに紐づくデータ取得
curl -XPOST -d'{}' localhost:8081/get/57a98d98e4b00679b4a830af

// 投稿
curl -XPOST -d'{"kamigo": "テストだよ", "nakashichi": "この投稿は",  "shimogo": "テストだよ", "userId": "57a98d98e4b00679b4a830af"}' localhost:8081/post
```
