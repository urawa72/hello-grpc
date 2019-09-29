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
  "google.golang.org/grpc/resolver"
  "google.golang.org/grpc/status"
)

func main() {
  // addr := "localhost:50051"
  // TLSで通信する
  // creds, err := credentials.NewClientTLSFromFile("server.crt", "")
  // if err != nil {
  //   log.Fatal(err)
  // }
  // conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))

  // intercepter
  // conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithUnaryInterceptor(unaryIntercepter))

  // load balancer
  resolver.Register(&exampleResolverBuilder{})
  addr := "testScheme:///example"
  conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBalancerName("round_robin"))

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

type exampleResolverBuilder struct{}

func (*exampleResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error) {
  r := &exampleResolver {
    target: target,
    cc: cc,
    addrsStore: map[string][]string {
      "example": {"localhost:50051", "localhost:50052"},
    },
  }
  r.start()
  return r, nil
}
func (*exampleResolverBuilder) Scheme() string { return "testScheme" }

type exampleResolver struct {
  target      resolver.Target
  cc          resolver.ClientConn
  addrsStore  map[string][]string
}

func (r *exampleResolver) start() {
  addrStrs := r.addrsStore[r.target.Endpoint]
  addrs := make([]resolver.Address, len(addrStrs))
  for i, s := range addrStrs {
    addrs[i] = resolver.Address{Addr: s}
  }
  r.cc.UpdateState(resolver.State{Addresses: addrs})
}
func (*exampleResolver) ResolveNow(o resolver.ResolveNowOption) {}
func (*exampleResolver) Close() {}

