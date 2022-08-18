package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"time"
	"workgrpc/pb"
)

const DBSAAS = "saasdb"

func NewBackupTaskRequest(vmIp, host, user, passwd string, Port int) *pb.BackupTaskRequest {
	b := &pb.BackupTaskRequest{
		WorkVm:          NewWorkVm(vmIp),
		MySQLConn:       NewMySQLConn(host, user, passwd, Port),
		BackUpType:      NewBackUpType(),
		RemoteStorageS3: nil,
		SaasDBMySQLConn: NewSaasMySQLConn(host, user, passwd, Port),
		BackUpTimeout:   5,
		DomainId:        100,
	}
	return b
}

func NewWorkVm(vmIp string) *pb.WorkVm {
	return &pb.WorkVm{WorkVm: vmIp}
}

func NewMySQLConn(host, user, passwd string, Port int) *pb.MySQLConn {
	return &pb.MySQLConn{
		MySQLHost:       host,
		MySQLUser:       user,
		MySQLUserpasswd: passwd,
		MySQLPort:       uint32(Port),
	}
}
func NewSaasMySQLConn(host, user, passwd string, Port int) *pb.SaasDBMySQLConn {
	return &pb.SaasDBMySQLConn{
		MySQLUser:       user,
		MySQLUserpasswd: passwd,
		MySQLHost:       host,
		MySQLPort:       uint32(Port),
		SaaSDBName:      DBSAAS,
	}
}

func NewBackUpType() *pb.BackUpType {
	return &pb.BackUpType{Type: pb.BackUpType_FullBackUpWithMySQLDump}
}
func main() {
	conn, err := grpc.Dial("127.0.0.1:3000", grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
	}
	client := pb.NewMySQLBackupServiceClient(conn)
	ctx, cancle := context.WithTimeout(context.Background(), time.Second)
	defer cancle()

	// Demo test 调用本地的Mysqldump命令 去备份容器 mysql-cmdb1的全部的表数据
	// 后续这些参数的生成由平台查询saas数据库的信息，通过节点选择策略来生成相关的参数值
	vmip := "127.0.0.1"
	host := "127.0.0.1"
	user := "root"
	passwd := "letsg0"
	port := 3307
	res, err := client.NewBackup(ctx, NewBackupTaskRequest(vmip, host, user, passwd, port))
	if err != nil {
		fmt.Println(err)
	}
	// Demo test 打印出收到的消息
	fmt.Println(res.GetMessageInfo())
	fmt.Println(res.GetMessageWarn())
}
