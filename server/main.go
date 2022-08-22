package main

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"workgrpc/pb"
)

func main() {
	listen, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Panic("xxxxxx")
	}
	s := grpc.NewServer()
	/* 备份恢复服务 */
	pb.RegisterMySQLBackupServiceServer(s, &BackUpServer{})
	/* Show Processlist 服务 */
	pb.RegisterMySQLShowProcessListServiceServer(s, &ShowProcessListServer{})

	// TODO 心跳表服务 读取到集群TOPO表中到mysql信息，然后记录心跳表 参考pt-heartbeat ，获取当前节点mysql的ip port，查询instance表中，该实例的角色信息，根据角色信息做读写心跳检测。 不支持单机多实例类型
	if err := s.Serve(listen); err != nil {
		log.Println(fmt.Errorf("run serve err :%w", err))
	}
}
