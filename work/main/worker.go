package main

import (
	"crontab/work"
	"flag"
	"fmt"
	"runtime"
	"time"
)

var (
	conFile string //配置文件路劲
)

func initArgs() {
	flag.StringVar(&conFile, "config", "worker.json", "配置文件")
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
	if err = worker.InitConfig(conFile); err != nil {
		goto ERR
	}

	//初始化任务管理器
	if err = worker.InitJobMgr(); err != nil {
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
