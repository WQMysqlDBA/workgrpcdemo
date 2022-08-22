package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"time"
	"workgrpc/pb"
)
func NewShowProcessListTask()*pb.ShowProcesslistRequest{
	return &pb.ShowProcesslistRequest{
		ShowMsg: "",
	}
}
func main()  {
	conn, err := grpc.Dial("127.0.0.1:3000", grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
	}
	client := pb.NewMySQLShowProcessListServiceClient(conn)
	ctx, cancle := context.WithTimeout(context.Background(), time.Second)
	defer cancle()

	res,err:=client.NewShowProcesslist(ctx,NewShowProcessListTask())
	fmt.Println(res.ProcessListInfo)
}
