package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"log"
	"net"
	"os"
	"os/exec"
	"sync"
	"time"
	"workgrpc/model"
	"workgrpc/pb"
)

var waitBackupWork sync.WaitGroup
var onlyOneBackJobRun sync.Mutex

type BackUpServer struct {
	pb.UnimplementedMySQLBackupServiceServer
}

func NewBackupServer() *BackUpServer {
	return &BackUpServer{}
}

type backjobmetadata pb.BackupTaskRequest

func (server *BackUpServer) NewBackup(con context.Context, req *pb.BackupTaskRequest) (*pb.BackupTaskResponse, error) {
	// TODO 留作后续讨论  备份任务 时间一般都是比较久都 ，是不是 可以把备份的任务，作为异步的任务，不需要通过context来做上下文的传递 ？这里先传递ctx 看后续是应用还是去掉
	backupType := req.GetBackUpType()
	mysqlConn := req.GetMySQLConn()
	saasdbConn := req.GetSaasDBMySQLConn()
	workVm := req.GetWorkVm()
	backuptime := req.GetBackUpTimeout()
	domainid := req.GetDomainId()
	// 在这里可以实现相关的 备份具体流程了
	// 可以导入其他包的备份恢复的方法了
	backup_cmd := fmt.Sprintf("mysqldump -h %v -u %v -p%v -P %v -A > %v_Full.sql", mysqlConn.MySQLHost, mysqlConn.MySQLUser, mysqlConn.MySQLUserpasswd, mysqlConn.MySQLPort, time.Now().Format("2006-01-02"))
	log.Printf("receive a NewBackup task with %v,%v,%v\nwill run command \"%v\"", backupType, mysqlConn, workVm, backup_cmd)

	job := &backjobmetadata{
		BackUpType:      backupType,
		MySQLConn:       mysqlConn,
		SaasDBMySQLConn: saasdbConn,
		WorkVm:          workVm,
		BackUpTimeout:   backuptime,
		DomainId:        domainid,
	}
	// 设置任务的超时时间
	//tt := time.Hour * time.Duration(job.BackUpTimeout)
	ctx := context.Background()
	db := SaasDB(job.SaasDBMySQLConn.MySQLUser, job.SaasDBMySQLConn.MySQLUserpasswd, job.SaasDBMySQLConn.MySQLHost, job.SaasDBMySQLConn.SaaSDBName, int(job.SaasDBMySQLConn.MySQLPort))

	go server.BackupJob(ctx, job, db)

	return &pb.BackupTaskResponse{MessageInfo: "Demo Test I get this BackupJob Message",
		MessageWarn: "Demo Test I get this BackupJob Message",
	}, nil
}

func (server *BackUpServer) BackupJob(ctx context.Context, j *backjobmetadata, db *gorm.DB) error {
	// TODO 加锁 defer unlock 这里必须加一个lock，否则可能出现多个备份进程存在的情况 不允许 同时有不同类型的备份任务
	// TODO 发起备份的时候，先检查 实例的状态是否是 available ，再决定是否可以进行备份
	onlyOneBackJobRun.Lock()
	defer onlyOneBackJobRun.Unlock()

	// below to test 调用BackupJob 是通过go命令调用的
	// server 运行着，这里其实不需要 `waitgroup`

	//array := [...]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	//array := [...]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	//for k, v := range array {
	//	fmt.Println(k, v)
	//	time.Sleep(time.Second)
	//}

	if j.BackUpType.Type == pb.BackUpType_FullBackUpWithXtra {
	} else if j.BackUpType.Type == pb.BackUpType_IncrBackUpWithXtra {

	} else if j.BackUpType.Type == pb.BackUpType_FullBackUpWithMydumper {

	} else if j.BackUpType.Type == pb.BackUpType_SingleTableBackUpWithMydumper {

	} else if j.BackUpType.Type == pb.BackUpType_FullBackUpWithMySQLDump {
		return server.UseMysqlDump(ctx, j, db)
	} else if j.BackUpType.Type == pb.BackUpType_SingleTableBackUpWithMySQLDump {
	}
	return nil
}

func (server *BackUpServer) UseMysqlDump(ctx context.Context, j *backjobmetadata, db *gorm.DB) error {
	// 记录备份任务开始 -> saas 数据库 通过gorm去insert 数据
	backLogUuid, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("生成备份日志的uuid记录失败 %w", err)
	}
	fmt.Println(backLogUuid)
	backlog := &model.BackLog{
		GvaModel:      model.GvaModel{},
		DomainId:      int(j.DomainId),
		BackupType:    "mysqldump",
		DataSize:      0,
		Status:        "backup",
		BackUpFeature: nil,
		BackUpUuid:    backLogUuid,
	}
	err = model.CreateBackLog(db, backlog)
	// TODO 如何保证这里一定是success ？或是把备份日志的创建 在平台调度的时候就先写好  这几个函数中的错误不会反馈给平台了，最好是统一做到日志里面
	if err != nil {
		return fmt.Errorf("创建备份日志到saas数据库失败 err:%w", err)
	}

	// 记录当前备份日志的log的uuid ,完成备份任务之后，更新备份是否成功 成功的话 还有feature的信息

	// 组装mysqldump的命令 并完成备份
	sqlFile := FileNameFormat() + "_full.sql"
	backupCmd := fmt.Sprintf("mysqldump -h %v -u %v -p%v -P %v -A > %v", j.MySQLConn.MySQLHost, j.MySQLConn.MySQLUser, j.MySQLConn.MySQLUserpasswd, j.MySQLConn.MySQLPort, sqlFile)
	if output, err := PubCmd(backupCmd, true); err != nil {
		setvars := make(map[string]interface{})
		setvars["status"] = "failed"
		model.UpdateBackLogByUuid(db, backLogUuid, setvars)
		return fmt.Errorf("run mydumper cmd error, output is  %v \n And the err is :%v ", output, err.Error())
	}

	// 到这里备份完成了，获取备份文件的大小 并更新saas数据库状态
	if backSqlFileInfo, err := os.Stat(sqlFile); err == nil {
		size := backSqlFileInfo.Size() // bytes
		status := "success"
		model.UpdateBackLogJsonByUuid(db, backLogUuid, status, backSqlFileInfo.Name(), size)
	}

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
		return fmt.Errorf("The command %v has some error when running it. Commands out is \n%v \nand the std err is %v\n", c, outStr, errStr)
	}
	return nil
}

func PubCmd(cmd string, shell bool) (string, error) {
	if shell {
		output, err := exec.Command("bash", "-c", cmd).Output()
		if err != nil {
			log.Println("cmd: ", cmd, " ", err, err.Error())
			return "", err
		}
		return string(output), err
	} else {
		output, err := exec.Command(cmd).Output()
		if err != nil {
			log.Println("cmd: ", cmd, " ", err, err.Error())
			return "", err
		}
		return string(output), err
	}
}

func FileNameFormat() string {
	nt := time.Now()
	return fmt.Sprintf("%v%v%v-%v%v%v", nt.Year(), nt.Month(), nt.Day(), nt.Hour(), nt.Minute(), nt.Second())
}

func SaasDB(user, passwd, host, db string, port int) *gorm.DB {
	return model.GormMysql(user, passwd, host, db, port)
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
