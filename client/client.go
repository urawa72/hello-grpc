package main

import (
  "context"
  "log"
  "os"
  "time"

  pb "github.com/urawa72/hello-grpc"
  "google.golang.org/genproto/googleapis/rpc/errdetails"
  "google.golang.org/grpc"
  // "google.golang.org/grpc/credentials"
  "google.golang.org/grpc/metadata"
  "google.golang.org/grpc/status"
)

func main() {
  addr := "localhost:50051"
  // TLSで通信する
  // creds, err := credentials.NewClientTLSFromFile("server.crt", "")
  // if err != nil {
  //   log.Fatal(err)
  // }
  // conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))

  // intercepter
  conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithUnaryInterceptor(unaryIntercepter))
  if err != nil {
    log.Fatalf("did not connect: %v", err)
  }
  defer conn.Close()
  c := pb.NewGreeterClient(conn)

  name := os.Args[1]

  // matadataを追加する
  ctx := context.Background()
  md := metadata.Pairs("timestamp", time.Now().Format(time.Stamp))
  ctx = metadata.NewOutgoingContext(ctx, md)
  r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name}, grpc.Trailer(&md))

  // キャンセル
  // ctx, cancel := context.WithCancel(context.Background())
  // defer cancel()
  // go func() {
  //   time.Sleep(1 * time.Second)
  //   cancel()
  // }()
  // r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})

  // エラーハンドリング
  if err != nil {
    s, ok := status.FromError(err)
    if ok {
      log.Printf("gRPC Error (message: %s)", s.Message())
      for _, d := range s.Details() {
        switch info := d.(type) {
        case *errdetails.RetryInfo:
          log.Printf("  RetryInfo: %v", info)
        }
      }
      os.Exit(1)
    } else {
      log.Fatalf("could not greet: %v", err)
    }
  }
  log.Printf("Greeting: %s", r.Message)
}

// interceptor
func unaryIntercepter(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
  log.Printf("before call: %s, request: %+v", method, req)
  err := invoker(ctx, method, req, reply, cc, opts...)
  log.Printf("after call: %s, response: %+v", method, reply)
  return err
}
