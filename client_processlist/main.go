package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/songzhibin97/gkit/tools/pretty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"time"
	"workgrpc/pb"
)

func NewShowProcessListTask() *pb.ShowProcesslistRequest {
	return &pb.ShowProcesslistRequest{
		ShowMsg: "",
	}
}

func main() {
	//creditsServeKey := "/Users/anderalex/go/src/workgrpc/certify/server.key"
	creditsServePem := "/Users/anderalex/go/src/workgrpc/certify/server.crt"
	commandName := "example.server.com"
	creds, _ := credentials.NewClientTLSFromFile(creditsServePem, commandName)

	conn, err := grpc.Dial("127.0.0.1:3000", grpc.WithTransportCredentials(creds))
	if err != nil {
		fmt.Println(err)
	}
	client := pb.NewMySQLShowProcessListServiceClient(conn)
	ctx, cancle := context.WithTimeout(context.Background(), time.Second)
	defer cancle()

	res, err := client.NewShowProcesslist(ctx, NewShowProcessListTask())
	if err != nil {
		fmt.Println(err)
	}
	if res != nil {
		strbyte, e := json.Marshal(res)
		if e != nil {
			fmt.Println(e)
		}
		// strbyte 是 []byte 类型，可以直接通过接口 func()gin.H{} 返回给前端json数组
		fmt.Printf("%s\n", pretty.Pretty(strbyte))
	}
}
