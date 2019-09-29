## 準備
* protobufインストール
```
brew install protobuf
```
* 追加プラグイン
```
go get -u github.com/golang/protobuf/protoc-gen-go
```
* PATHに追加
```
export PATH=$PATH:$GOPATH/bin
```

## 実行
* GoのgRPC用コード出力
```
protoc xxxx.proto --go_out=plugins=grpc:.
```
* server起動してclient.go実行
```
go run server.go
go run client.go World
```
* デバッグをONにする
```
export GODEBUG=http2debug=2
```

## オレオレ証明書
* 以下の3つを`./server`に配置
* `server.crt`は`./client`にも
```
# 秘密鍵生成
openssl genrsa 2048 > private.key
# 署名要求作成 Common Nameにはlocalhostを設定
openssl req -new -key private.key > server.csr
# 証明書作成
openssl x509 -days 367 -req -signkey private.key < server.csr > server.crt
```
