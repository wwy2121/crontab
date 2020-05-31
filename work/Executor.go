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
			cmd     *exec.Cmd
			err     error
			ouput   []byte
			result  *common.JobExecuteResult
			jobLock *JobLock
		)
		//任务结果
		result = &common.JobExecuteResult{
			ExecuteInfo: info,
			Output:      make([]byte, 0),
		}
		//初始化分布式锁
		jobLock = G_jobMgr.CreateJobLock(info.Job.Name)

		//任务开始时间
		result.StartTime = time.Now()

		//上锁
		err = jobLock.TryLock()
		//释放锁
		defer jobLock.Unlock()

		if err != nil {
			result.Err = err
			result.EndTime = time.Now()
		} else {
			//上锁成功后 重置任务开始时间
			result.StartTime = time.Now()
			//执行shell命令
			cmd = exec.CommandContext(context.TODO(), "/bin/bash", "-c", info.Job.Command)
			//捕获输出
			ouput, err = cmd.CombinedOutput()
			//任务结束时间
			result.EndTime = time.Now()
			result.Output = ouput
			result.Err = err
		}
		//回传
		G_scheduler.PushJobResult(result)
	}()
}

func InitExecutor() (err error) {
	G_executor = &Executor{}
	return
}
