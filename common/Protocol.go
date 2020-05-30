package common

import (
	"encoding/json"
	"strings"
)

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

type JobEvent struct {
	EventType int
	job       *Job
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

func UnpackJob(value []byte) (ret *Job, err error) {
	var (
		job *Job
	)
	job = &Job{}
	if err = json.Unmarshal(value, job); err != nil {
		return
	}
	ret = job
	return
}

func ExtractJobName(jobKey string) string {
	return strings.TrimPrefix(jobKey, JOB_SAVE_DIR)
}

func BuildJobEvent(eventType int, job *Job) (jobEvent *JobEvent) {
	return &JobEvent{
		EventType: eventType,
		job:       job,
	}
}
