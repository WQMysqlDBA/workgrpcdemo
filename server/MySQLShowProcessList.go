package main

import (
	"context"
	"encoding/json"
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
	pL := make([]*pb.ProcessListInfo, 0, 0)
	for _, v := range processList {
		var a *pb.ProcessListInfo
		if data, err := json.MarshalIndent(v, "", "\t"); err == nil {
			err1 := json.Unmarshal(data, &a)
			if err1 != nil {
				return nil, err1
			}
			pL = append(pL, a)
		}
	}
	processlistInfo := &pb.ShowProcesslistResponce{ProcessListInfo: pL}
	return processlistInfo, nil
}
