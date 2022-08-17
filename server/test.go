package main

import (
	"context"
	"fmt"
	hello_grpc "workgrpc/pb"
)

type server struct {
	hello_grpc.UnimplementedHelloGrpcServer
}

func (s *server) SayHi(ctx context.Context, req *hello_grpc.Req) (res *hello_grpc.Res, err error) {
	fmt.Println(req.GetMessage())
	return &hello_grpc.Res{Message: "我是从服务端返回的grpc的内容"}, err
}

//func main(){
//	listen ,err := net.Listen("tcp",":3000")
//   if err!=nil{
//   	log.Panic("xxxxxx")
//	}
//	s:=grpc.NewServer()
//	hello_grpc.RegisterHelloGrpcServer(s,&server{})
//	s.Serve(listen)
//}