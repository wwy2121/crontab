package main

import (
	"crontab/master"
	"flag"
	"fmt"
	"runtime"
	"time"
)

var (
	conFile string //配置文件路劲
)

func initArgs() {
	flag.StringVar(&conFile, "config", "master.json", "配置文件")
	flag.Parse()
}

//初始化线程
func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {

	var (
		err error
	)
	//初始化命令行
	initArgs()

	//初始化线程
	initEnv()

	//加载配置
	if err = master.InitConfig(conFile); err != nil {
		goto ERR
	}

	//任务管理器
	if err = master.InitJobMgr(); err != nil {
		goto ERR
	}

	//启动api http服务
	if err := master.InitServer(); err != nil {
		goto ERR
	}

	//正常退出
	for {
		time.Sleep(100 * time.Second)
	}
	return
ERR:
	fmt.Print(err)
}
