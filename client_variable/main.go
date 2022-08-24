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

/*      show vars         */
func ShowVariablesTask(s []string) *pb.HandleVariablesRequest {
	return &pb.HandleVariablesRequest{
		Method:                false, // show
		ShowVariablesUseArray: NewVariables(s),
	}
}

func NewVariables(s []string) (vv []*pb.ShowVariablesUseArray) {

	for _, v := range s {
		b := &pb.ShowVariablesUseArray{Var: v}
		vv = append(vv, b)
	}
	return vv
}

/*      set vars         */
func SetVariablesStringTask(v map[string]string) *pb.HandleVariablesRequest {
	a := make(map[string]*pb.SetVariablesUseMap)
	for k, v := range v {
		a[k] = &pb.SetVariablesUseMap{VariableValue: SetVariablesUseMap_VariableValueString(v)}
	}
	return &pb.HandleVariablesRequest{
		Method:             true, // set
		SetVariablesUseMap: a,
	}
}

func SetVariablesInt32Task(v map[string]string) *pb.HandleVariablesRequest {
	a := make(map[string]*pb.SetVariablesUseMap)
	for k, v := range v {
		a[k] = &pb.SetVariablesUseMap{VariableValue: SetVariablesUseMap_VariableValueInt32(v)}
	}
	return &pb.HandleVariablesRequest{
		Method:             true, // set
		SetVariablesUseMap: a,
	}
}

func SetVariablesUseMap_VariableValueString(v string) *pb.SetVariablesUseMap_VariableValueString {
	return &pb.SetVariablesUseMap_VariableValueString{VariableValueString: v}
}

func SetVariablesUseMap_VariableValueInt32(v string) *pb.SetVariablesUseMap_VariableValueInt32 {
	return &pb.SetVariablesUseMap_VariableValueInt32{VariableValueInt32: v}
}

func main() {
	creditsServePem := "/Users/anderalex/go/src/workgrpc/certify/server.crt"
	commandName := "example.server.com"
	creds, _ := credentials.NewClientTLSFromFile(creditsServePem, commandName)

	conn, err := grpc.Dial("127.0.0.1:3000", grpc.WithTransportCredentials(creds))
	if err != nil {
		fmt.Println(err)
	}
	client := pb.NewMySQLVariablesServiceClient(conn)
	ctx, cancle := context.WithTimeout(context.Background(), time.Second)
	defer cancle()
	var testvariable = []string{"innodb_buffer_pool_size", "sql_mode", "version", "transaction_isolation"}
	res, err := client.VariablesHandler(ctx, ShowVariablesTask(testvariable))
	if err != nil {
		fmt.Println(err)
	}
	for k, v := range res.VariableName {
		fmt.Printf("variable %v  value : %v\n", k, v)
	}
	if res != nil {
		strbyte, e := json.Marshal(res)
		if e != nil {
			fmt.Println(e)
		}
		// strbyte 是 []byte 类型，可以直接通过接口 func()gin.H{} 返回给前端json数组
		fmt.Printf("%s\n", pretty.Pretty(strbyte))
	}

	setvars := make(map[string]string)
	//setvars["sql_mode"] = ""
	setvars["innodb_buffer_pool_size"]="1024*1024*1024"
	//_, err = client.VariablesHandler(ctx, SetVariablesStringTask(setvars))
	_, err = client.VariablesHandler(ctx, SetVariablesInt32Task(setvars))
	if err != nil {

		fmt.Println(err.Error())
	}
}
