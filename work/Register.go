package worker

import (
	"crontab/common"
	"go.etcd.io/etcd/clientv3"
	"golang.org/x/net/context"
	"net"
	"time"
)

type Register struct {
	client  *clientv3.Client
	kv      clientv3.KV
	lease   clientv3.Lease
	localIp string
}

var (
	G_register *Register
)

func getLocalIp() (ipv4 string, err error) {
	var (
		addrs   []net.Addr
		addr    net.Addr
		ipNet   *net.IPNet //ip地址
		isIpNet bool
	)
	if addrs, err = net.InterfaceAddrs(); err != nil {
		return
	}
	for _, addr = range addrs {
		if ipNet, isIpNet = addr.(*net.IPNet); isIpNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ipv4 = ipNet.IP.String()
				return
			}
		}
	}
	err = common.ERR_NOT_LOCAL_IP_FOUND
	return
}

func (register *Register) keepOnline() {
	var (
		regKey          string
		leaseGrantResp  *clientv3.LeaseGrantResponse
		err             error
		keepAliveChan   <-chan *clientv3.LeaseKeepAliveResponse
		keepAliveRespon *clientv3.LeaseKeepAliveResponse
		cancelCtx       context.Context
		cancelFunc      context.CancelFunc
	)
	regKey = common.JOB_WORKER_DIR + register.localIp
	cancelFunc = nil
	if leaseGrantResp, err = register.lease.Grant(context.TODO(), 10); err != nil {
		goto RETRY
	}

	if keepAliveChan, err = register.lease.KeepAlive(context.TODO(), leaseGrantResp.ID); err != nil {
		goto RETRY
	}
	cancelCtx, cancelFunc = context.WithCancel(context.TODO())
	if _, err = register.kv.Put(cancelCtx, regKey, "", clientv3.WithLease(leaseGrantResp.ID)); err != nil {
		goto RETRY
	}
	for {
		select {
		case keepAliveRespon = <-keepAliveChan:
			if keepAliveRespon == nil {
				goto RETRY
			}
		}
	}

RETRY:
	time.Sleep(1 * time.Second)
	if cancelFunc != nil {
		cancelFunc()
	}
}

func InitRegister() (err error) {
	var (
		config  clientv3.Config
		client  *clientv3.Client
		kv      clientv3.KV
		lease   clientv3.Lease
		localIp string
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
	if localIp, err = getLocalIp(); err != nil {
		return
	}
	//得到kv
	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)
	G_register = &Register{
		client:  client,
		kv:      kv,
		lease:   lease,
		localIp: localIp,
	}

	go G_register.keepOnline()
	return
}
