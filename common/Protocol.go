package common

import "encoding/json"

//定时任务
type Job struct {
	Name     string `json:"name"`     //任务名
	Command  string `json:"command"`  //shell 命令
	CronExpr string `json:"cronExpr"` //cron表达式
}

//http返回应答
type Response struct {
	Errno int    `json:"errno"`
	Msg   string `json:"msg"`
	Data  interface{}
}

func BuildResponse(errno int, msg string, data interface{}) (rep []byte, err error) {
	var (
		response Response
	)
	response.Errno = errno
	response.Msg = msg
	response.Data = data

	rep, err = json.Marshal(response)

	return
}
