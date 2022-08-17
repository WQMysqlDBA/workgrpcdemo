package main

import (
	"bytes"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os/exec"
	"time"
	"workgrpc/pb"
)

type BackUpServer struct {
	pb.UnimplementedMySQLBackupServiceServer
}

func NewBackupServer() *BackUpServer {
	return &BackUpServer{}
}

type backjobmetadata pb.BackupTaskRequest

func (server *BackUpServer) NewBackup(ctx context.Context, req *pb.BackupTaskRequest) (*pb.BackupTaskResponse, error) {
	// TODO 留作后续讨论  备份任务 时间一般都是比较久都 ，是不是 可以把备份的任务，作为异步的任务，不需要通过context来做上下文的传递 ？这里先传递ctx 看后续是应用还是去掉
	backupType := req.GetBackUpType()
	mysqlConn := req.GetMySQLConn()
	saasdbConn := req.GetSaasDBMySQLConn()
	workVm := req.GetWorkVm()
	// 在这里可以实现相关的 备份具体流程了
	// 可以导入其他包的备份恢复的方法了
	backup_cmd := fmt.Sprintf("mysqldump -h %v -u %v -p%v -P %v -A > %v_Full.sql", mysqlConn.MySQLHost, mysqlConn.MySQLUser, mysqlConn.MySQLUserpasswd, mysqlConn.MySQLPort, time.Now().Format("2006-01-02"))
	log.Printf("receive a NewBackup task with %v,%v,%v\nwill run command \"%v\"", backupType, mysqlConn, workVm, backup_cmd)

	job := &backjobmetadata{
		BackUpType:      backupType,
		MySQLConn:       mysqlConn,
		SaasDBMySQLConn: saasdbConn,
		WorkVm:          workVm,
	}
	err := server.BackupJob(ctx, job)
	return &pb.BackupTaskResponse{MessageInfo: "Demo Test I get the Message ",
		MessageWarn: "Demo Test here is a Warn Message",
	}, err
}

func (server *BackUpServer) BackupJob(ctx context.Context, j *backjobmetadata) error {
	if j.BackUpType.Type == pb.BackUpType_FullBackUpWithXtra {
	} else if j.BackUpType.Type == pb.BackUpType_IncrBackUpWithXtra {

	} else if j.BackUpType.Type == pb.BackUpType_FullBackUpWithMydumper {

	} else if j.BackUpType.Type == pb.BackUpType_SingleTableBackUpWithMydumper {

	} else if j.BackUpType.Type == pb.BackUpType_FullBackUpWithMySQLDump {
		return server.UseMysqlDump(ctx, j)
	} else if j.BackUpType.Type == pb.BackUpType_SingleTableBackUpWithMySQLDump {
	}
	return nil
}

func (server *BackUpServer) UseMysqlDump(ctx context.Context, j *backjobmetadata) error {
	// firstly ,make sure mysqldump command is available is import
	var a []string
	a = append(a, "-V")
	if err := PubRunCmd("mysqldump", a); err != nil {
		return err
	}
	// 记录备份任务开始 -> saas 数据库 通过gorm去insert 数据

	// 组装mysqldump的命令 并完成备份

	// 更新saas数据库状态

	// 这里需要使用两个go routine 和一个chan 来做通信，获得 备份任务的go routine的状态 并更新到saas数据库中

	return nil
}

func PubRunCmd(c string, a []string) error {
	cmd := exec.Command(c, a...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	if err != nil {
		return fmt.Errorf("The mysqldump command has some error when running it. Commands out is \n%v \nand the std err is %v\n", outStr, errStr)
	}
	return nil
}

func main() {
	listen, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Panic("xxxxxx")
	}
	s := grpc.NewServer()
	pb.RegisterMySQLBackupServiceServer(s, &BackUpServer{})
	if err := s.Serve(listen); err != nil {
		log.Println(fmt.Errorf("run serve err :%w", err))
	}

}
