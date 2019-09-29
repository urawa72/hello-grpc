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
