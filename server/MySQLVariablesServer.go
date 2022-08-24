package main

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/any"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"workgrpc/model"
	"workgrpc/pb"
)

type VariablesServer struct {
	pb.UnimplementedMySQLVariablesServiceServer
}

func (server *VariablesServer) VariablesHandler(ctx context.Context, req *pb.HandleVariablesRequest) (*pb.HandleVariablesResponse, error) {
	if !req.GetMethod() {
		// show variables
		return server.ShowVariables(req.GetShowVariablesUseArray())
	} else {
		// set variables
		isetok, err := server.SetVariables(req.GetSetVariablesUseMap())
		if !isetok || err != nil {
			return &pb.HandleVariablesResponse{
				SetOK: false,
			}, err
		} else {
			return &pb.HandleVariablesResponse{
				SetOK: true,
			}, nil
		}
	}
}

func (server *VariablesServer) ShowVariables(arr []*pb.ShowVariablesUseArray) (res *pb.HandleVariablesResponse, err error) {
	variablesValue := make(map[string]string)
	for _, v := range arr {
		variablesValue[v.Var], err = server.GetVariableRunningValue(v.Var)
		if err != nil {
			return res, err
		}
	}
	res = &pb.HandleVariablesResponse{
		VariableName: variablesValue,
	}
	return res, nil
}

// 精确匹配参数
func (server *VariablesServer) GetVariableRunningValue(variable string) (string, error) {
	// use gorm to get value and defer close session
	db, err := model.GormMysql("root", "letsg0", "127.0.0.1", "information_schema", 3307)
	if err != nil {
		return "", err
	}
	defer func() {
		if sqlDB, err := db.DB(); err == nil {
			err = sqlDB.Close()
		}
	}()

	type Variables struct {
		VariableName string `gorm:"column:Variable_name"`
		Value        string `gorm:"column:Value"`
	}
	var v Variables
	sql := fmt.Sprintf("SHOW VARIABLES LIKE \"%v\"", variable)
	db.Debug().Raw(sql).Scan(&v)
	if v.Value != "" {
		return v.Value, nil
	}
	return "NONE_VALUE_RETURNED", nil
}

// 模糊匹配参数
func (server *VariablesServer) GetVariableRunningValueFuzzyMatching(variable string) (string, error) {
	// use gorm to get value and defer close session
	db, err := model.GormMysql("root", "letsg0", "127.0.0.1", "information_schema", 3307)
	if err != nil {
		return "", err
	}
	defer func() {
		if sqlDB, err := db.DB(); err == nil {
			err = sqlDB.Close()
		}
	}()

	type Variables struct {
		VariableName string `gorm:"column:Variable_name"`
		Value        string `gorm:"column:Value"`
	}
	var v Variables
	sql := "SHOW VARIABLES LIKE " + "\"%" + variable + "%\""
	db.Debug().Raw(sql).Scan(&v)
	if v.Value != "" {
		return v.Value, nil
	}
	return "NONE_VALUE_RETURNED", nil
}

// 修改参数 每次只支持单个参数修改
func (server *VariablesServer) SetVariables(setvariablemap map[string]*pb.SetVariablesUseMap) (bool, error) {
	// use gorm to get value and defer close session
	db, err := model.GormMysql("root", "letsg0", "127.0.0.1", "information_schema", 3307)
	if err != nil {
		return false, err
	}
	defer func() {
		if sqlDB, err := db.DB(); err == nil {
			err = sqlDB.Close()
		}
	}()

	if len(setvariablemap) != 1 {
		return false, fmt.Errorf("当前版本不支持同时修改多个参数")
	} else {
		sql := ""
		for k, v := range setvariablemap {
			// 判断value消息
			switch v.VariableValue.(type) {
			case *pb.SetVariablesUseMap_VariableValueString:
				fmt.Println("string",k,v)
				sql = fmt.Sprintf("SET GLOBAL %v = '%v'", k, v.GetVariableValueString())
			case *pb.SetVariablesUseMap_VariableValueInt32:
				fmt.Println("int",k,v)
				sql = fmt.Sprintf("SET GLOBAL %v = %v", k, v.GetVariableValueInt32())
			}
		}
		err=db.Debug().Raw(sql).Exec(sql).Error
		if err!=nil{
			return false, fmt.Errorf("MySQL Server Returns Error when exec sql,sql is %v ,err is %v ",sql,err)
		}
		return true, nil
	}
}

func Unmarshal(data *any.Any) (*pb.SetVariablesUseMap, error) {
	var m pb.SetVariablesUseMap
	err := anypb.UnmarshalTo(data, &m, proto.UnmarshalOptions{})
	return &m, err
}
