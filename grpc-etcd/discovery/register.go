package discovery

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	clientV3 "go.etcd.io/etcd/client/v3"
	"strings"
	"time"
)

type Register struct {
	endpoints   []string // etcd 节点地址
	dialTimeout int      // 建立连接失败的超时时

	keepAliveChan <-chan *clientV3.LeaseKeepAliveResponse // 协议消息

	leaseID   clientV3.LeaseID // 租约ID
	serverTTL int64            // 租约时间

	closeCh chan struct{} // 定义一个手动关闭通道

	serverInfo Server // 服务详细信息

	cli    *clientV3.Client // etcd 客户端会话
	logger *logrus.Logger   // 日志
}

func NewRegister(endpoints []string, logger *logrus.Logger) *Register {
	return &Register{
		endpoints:   endpoints,
		dialTimeout: 3,
		logger:      logger,
	}
}

func (r *Register) Register(serverInfo Server, ttl int64) error {
	var err error

	// 判断服务地址是否正确
	if strings.Split(serverInfo.Addr, ":")[0] == "" {
		return errors.New("无效的服务地址")
	}

	r.serverInfo = serverInfo
	r.serverTTL = ttl

	// 创建 etcd 会话
	if r.cli, err = clientV3.New(clientV3.Config{
		Endpoints:   r.endpoints,
		DialTimeout: time.Duration(r.dialTimeout) * time.Second,
	}); err != nil {
		return err
	}

	if err = r.register(); err != nil {
		return err
	}

	r.closeCh = make(chan struct{})

	go r.keepAlive()

	return nil
}

func (r *Register) register() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.dialTimeout)*time.Second)
	defer cancel()

	// 创建一个租约，每隔 10s 需要向 etcd 汇报一次心跳，证明当前节点仍然存活
	grant, err := r.cli.Grant(ctx, r.serverTTL)
	if err != nil {
		return err
	}
	r.leaseID = grant.ID

	// 设置续租 定期发送需求请求
	// KeepAlive 使给定的租约永远有效。如果发布到通道的keepalive响应没有立即使用，
	// 则租约客户端至少每秒钟向etcd服务器发送保持活动请求，直到获取最新响应为止
	// etcd client 会自动发送ttl到etcd server， 从而保证租约一直有效
	if r.keepAliveChan, err = r.cli.KeepAlive(context.Background(), r.leaseID); err != nil {
		return err
	}

	// 转化成 json
	data, err := json.Marshal(r.serverInfo)
	if err != nil {
		return err
	}
	key := BuildRegisterPath(r.serverInfo)

	// 存住并且绑定租约
	_, err = r.cli.Put(context.Background(), key, string(data), clientV3.WithLease(r.leaseID))
	return err
}

func (r *Register) keepAlive() {
	// 周期性定时器
	ticker := time.NewTicker(time.Duration(r.dialTimeout) * time.Second)

	for {
		select {
		case <-r.closeCh:
			if err := r.Unregister(); err != nil {
				r.logger.Fatal(fmt.Sprintf("keepAlive r.Unregister, ERROR: %v", err))
			}
			return
		case res := <-r.keepAliveChan: // 如果 keepAliveChan 为空 则重新注册
			if res == nil {
				if err := r.register(); err != nil {
					r.logger.Fatal(fmt.Sprintf("keepAlive r.register, ERROR: %v", err))
				}
			}
		case <-ticker.C:
			if r.keepAliveChan == nil {
				if err := r.register(); err != nil {
					r.logger.Fatal(fmt.Sprintf("keepAlive r.register, ERROR: %v", err))
				}
			}
		}
	}
}

// Stop 手动关闭
func (r *Register) Stop() {
	r.closeCh <- struct{}{}
}

// Unregister 注销服务
func (r *Register) Unregister() error {
	// 撤销租约
	if _, err := r.cli.Revoke(context.Background(), r.leaseID); err != nil {
		return err
	}
	r.logger.Println("撤销租约")

	// 注销服务
	if _, err := r.cli.Delete(context.Background(), BuildRegisterPath(r.serverInfo)); err != nil {
		return err
	}
	r.logger.Println("注销服务")
	// 关闭 etcd 客户端链接
	return r.cli.Close()
}
