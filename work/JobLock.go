package worker

import (
	"crontab/common"
	"go.etcd.io/etcd/clientv3"
	"golang.org/x/net/context"
)

//分布式锁(TXN事务)
type JobLock struct {
	kv         clientv3.KV
	lease      clientv3.Lease
	jobName    string
	cancelFunc context.CancelFunc
	leaseId    clientv3.LeaseID
}

//初始化锁
func InitJobLock(jobName string, kv clientv3.KV, lease clientv3.Lease) (jobLock *JobLock) {
	jobLock = &JobLock{
		kv:      kv,
		lease:   lease,
		jobName: jobName,
	}
	return
}

//尝试上锁
func (jobLock *JobLock) Try() (err error) {
	var (
		leaseGrantResp *clientv3.LeaseGrantResponse
		cancelCtx      context.Context
		cancelFunc     context.CancelFunc
		leaseId        clientv3.LeaseID
		keepRespChan   <-chan *clientv3.LeaseKeepAliveResponse
		txn            clientv3.Txn
		lockKey        string
		txnRespon      *clientv3.TxnResponse
	)

	if leaseGrantResp, err = jobLock.lease.Grant(context.TODO(), 5); err != nil {
		return
	}
	cancelCtx, cancelFunc = context.WithCancel(context.TODO())
	leaseId = leaseGrantResp.ID
	if keepRespChan, err = jobLock.lease.KeepAlive(cancelCtx, leaseId); err != nil {
		goto FAIL
	}
	go func() {
		var (
			keepResp *clientv3.LeaseKeepAliveResponse
		)
		for {
			select {
			case keepResp = <-keepRespChan:
				if keepResp == nil {
					goto END
				}
			}
		}
	END:
	}()
	txn = jobLock.kv.Txn(context.TODO())
	lockKey = common.JOB_LOCK_DIR + jobLock.jobName
	txn.If(clientv3.Compare(clientv3.CreateRevision(lockKey), "=", 0)).
		Then(clientv3.OpPut(lockKey, "", clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet(lockKey))
	if txnRespon, err = txn.Commit(); err != nil {
		goto FAIL
	}
	if !txnRespon.Succeeded {
		err = common.ERR_LOCK_ALREADY_REQUIRED
		goto FAIL
	}
	//枪锁成功
	return

FAIL:
	cancelFunc()
	jobLock.lease.Revoke(context.TODO(), leaseId)
	return
}

func (jobLock *JobLock) Unlock() {

}
