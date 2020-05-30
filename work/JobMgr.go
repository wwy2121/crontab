package worker

import (
	"go.etcd.io/etcd/clientv3"
	"time"
)

//任务管理器
type JobMgr struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

var (
	G_jobMgr *JobMgr
)

//监听任务的变化
func (jobMgr *JobMgr) watchJobs() (err error) {
	return
}

func InitJobMgr() (err error) {
	var (
		config clientv3.Config
		client *clientv3.Client
		kv     clientv3.KV
		lease  clientv3.Lease
	)
	//初始化配置
	config = clientv3.Config{
		Endpoints:   G_config.EtcdEndpoints, //集群地址
		DialTimeout: time.Duration(G_config.EtcdDialTimeout) * time.Millisecond,
	}
	//建立连接
	if client, err = clientv3.New(config); err != nil {
		return
	}
	//得到kv
	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)
	G_jobMgr = &JobMgr{
		client: client,
		kv:     kv,
		lease:  lease,
	}
	return
}
