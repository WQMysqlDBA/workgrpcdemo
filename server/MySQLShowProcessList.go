package main

import (
	"context"
	"encoding/json"
	"fmt"
	"workgrpc/model"
	"workgrpc/pb"
)

type ShowProcessListServer struct {
	pb.UnimplementedMySQLShowProcessListServiceServer
}

func (server *ShowProcessListServer) NewShowProcesslist(ctx context.Context, req *pb.ShowProcesslistRequest) (*pb.ShowProcesslistResponce, error) {

	db, err := model.GormMysql("root", "letsg0", "127.0.0.1", "information_schema", 3307)
	if err != nil {
	}
	var r model.InformationSchemaProcesslist
	processList := r.GetAllProcesslist(db)

	for _, v := range processList {
		if data, err := json.MarshalIndent(v, "", "\t"); err == nil {
			fmt.Println(string(data))

		}
	}
	processlistInfo := &pb.ShowProcesslistResponce{}
	// TODO 这里先把数据拿到 ， 然后结构体copy的方法 到 return 这里
	return processlistInfo, nil
}
