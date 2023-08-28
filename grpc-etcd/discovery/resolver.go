package discovery

import (
	"context"
	"github.com/sirupsen/logrus"
	clientV3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
	"grpc-etcd/discovery/balancer/weight"
	"time"
)

const (
	scheme = "etcd"
)

type Resolver struct {
	scheme      string
	endpoints   []string // etcd 节点地址
	dialTimeout int      // 建立连接失败的超时时间

	closeCh    chan struct{} // 定义一个手动关闭通道
	watchCh    clientV3.WatchChan
	serverList []resolver.Address // 服务列表
	keyPrefix  string             // key前缀

	cli    *clientV3.Client // etcd 客户端会话
	cc     resolver.ClientConn
	logger *logrus.Logger // 日志
}

// Scheme 返回此解析器支持的模式
func (r *Resolver) Scheme() string {
	return r.scheme
}

func NewResolver(endpoints []string, logger *logrus.Logger) *Resolver {
	return &Resolver{
		scheme:      scheme,
		endpoints:   endpoints,
		logger:      logger,
		dialTimeout: 3,
	}
}

// ResolveNow watch有变化以后会调用
func (r *Resolver) ResolveNow(options resolver.ResolveNowOptions) {
	r.logger.Println("ResolveNow")
	r.logger.Println(options)
}

// Build 构建解析器 grpc.Dial()同步调用
func (r *Resolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r.cc = cc

	r.keyPrefix = target.URL.Path

	if err := r.start(); err != nil {
		return nil, err
	}

	return r, nil
}

// start 开始执行
func (r *Resolver) start() error {
	var err error
	// 创建 etcd 会话
	if r.cli, err = clientV3.New(clientV3.Config{
		Endpoints:   r.endpoints,
		DialTimeout: time.Duration(r.dialTimeout) * time.Second,
	}); err != nil {
		return err
	}

	if err := r.sync(); err != nil {
		return err
	}

	go r.watch()

	return nil
}

// sync 同步获取所有地址信息
func (r *Resolver) sync() error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := r.cli.Get(ctx, r.keyPrefix, clientV3.WithPrefix())
	if err != nil {
		return err
	}

	// 初始化
	if res.Kvs != nil {
		for _, v := range res.Kvs {
			info, err := ParseValue(v.Value)
			if err != nil {
				continue
			}
			if !Exist(r.serverList, info.Addr) {

				addr := weight.SetAddrInfo(resolver.Address{
					Addr:       info.Addr,
					ServerName: info.Name,
				}, weight.AddrInfo{Weight: info.Weight})

				r.serverList = append(r.serverList, addr)
			}
		}

		return r.cc.UpdateState(resolver.State{Addresses: r.serverList})
	}
	return nil
}

// watch 监听etcd中某个key前缀的服务地址列表的变化
func (r *Resolver) watch() {
	// 周期性定时器
	ticker := time.NewTicker(time.Duration(r.dialTimeout) * time.Second)
	r.watchCh = r.cli.Watch(context.Background(), r.keyPrefix, clientV3.WithPrefix())

	for {
		select {
		case <-r.closeCh:
			return
		case res, ok := <-r.watchCh:
			if ok {
				r.update(res.Events)
			}
		case <-ticker.C:
			if err := r.sync(); err != nil {
				r.logger.Fatal("sync failed ", err)
			}
			r.logger.Println(r.serverList)
		}
	}
}

func (r *Resolver) update(events []*clientV3.Event) {
	for _, ev := range events {
		var info Server
		var err error

		switch ev.Type {
		case clientV3.EventTypePut:
			if info, err = ParseValue(ev.Kv.Value); err != nil {
				continue
			}
			r.logger.Println("PUT: ", info)
			if !Exist(r.serverList, info.Addr) {

				addr := weight.SetAddrInfo(resolver.Address{
					Addr:       info.Addr,
					ServerName: info.Name,
				}, weight.AddrInfo{Weight: info.Weight})

				r.serverList = append(r.serverList, addr)
				_ = r.cc.UpdateState(resolver.State{Addresses: r.serverList})
			}

		case clientV3.EventTypeDelete:
			if info, err = SplitPath(string(ev.Kv.Key)); err != nil {
				continue
			}

			r.logger.Println("DELETE: ", info)

			if serverList, ok := Remove(r.serverList, info.Addr); ok {
				r.serverList = serverList
				_ = r.cc.UpdateState(resolver.State{Addresses: r.serverList})
			}
		}
	}
}

// Exist 判断这个服务地址是否已经存在，防止服务访问冲突
func Exist(l []resolver.Address, addr string) bool {
	for i := range l {
		if l[i].Addr == addr {
			return true
		}
	}

	return false
}

// Remove 从服务列表中移除服务
func Remove(s []resolver.Address, addr string) ([]resolver.Address, bool) {
	for i := range s {
		if s[i].Addr == addr {
			s[i] = s[len(s)-1]
			return s[:len(s)-1], true
		}
	}
	return nil, false
}

// Close 手动关闭
func (r *Resolver) Close() {
	r.closeCh <- struct{}{}
}
