package main

import (
	"crontab/master"
	"fmt"
	"runtime"
)

func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {

	var (
		err error
	)

	//初始化线程
	initEnv()

	//启动api http服务
	if err := master.InitServer(); err != nil {
		goto ERR
	}
ERR:
	fmt.Print(err)
}
