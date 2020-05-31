package worker

import (
	"crontab/common"
	"golang.org/x/net/context"
	"os/exec"
	"time"
)

//任务执行器
type Executor struct {
}

var (
	G_executor *Executor
)

//执行任务
func (executor *Executor) ExecueJob(info *common.JobExecuteInfo) {
	go func() {
		var (
			cmd    *exec.Cmd
			err    error
			ouput  []byte
			result *common.JobExecuteResult
		)
		//任务结果
		result = &common.JobExecuteResult{
			ExecuteInfo: info,
			Output:      make([]byte, 0),
		}
		//任务开始事件
		result.StartTime = time.Now()
		//执行shell命令
		cmd = exec.CommandContext(context.TODO(), "/bin/bash", "-c", info.Job.Command)
		//捕获输出
		ouput, err = cmd.CombinedOutput()
		//任务结束时间
		result.EndTime = time.Now()
		result.Output = ouput
		result.Err = err
		//回传
		G_scheduler.PushJobResult(result)

	}()
}

func InitExecutor() (err error) {
	G_executor = &Executor{}
	return
}
