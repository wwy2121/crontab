package master

import (
	"net"
	"net/http"
	"strconv"
	"time"
)

//任务的http接口
type ApiServer struct {
	httpServer *http.Server
}

var (
	G_apiserver *ApiServer
)

//保存任务接口
func handleJobSave(w http.ResponseWriter, r *http.Request) {
	//任务保存到rtcd中

}

//初始化服务
func InitServer() (err error) {
	var (
		mux        *http.ServeMux
		listener   net.Listener
		httpServer *http.Server
	)
	//配置路由
	mux = http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)

	//启动tcp监听
	if listener, err = net.Listen("tcp", ":"+strconv.Itoa(G_config.ApiPort)); err != nil {
		return
	}

	//创建http服务器
	httpServer = &http.Server{
		ReadTimeout:  time.Duration(G_config.ApiReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(G_config.ApiWriteTimeout) * time.Millisecond,
		Handler:      mux,
	}
	//单例赋值
	G_apiserver = &ApiServer{
		httpServer: httpServer,
	}
	//启动服务端
	go httpServer.Serve(listener)
	return
}
