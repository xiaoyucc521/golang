package main

import (
	"context"
	"etcd/config"
	"etcd/discovery"
	"fmt"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientV3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// ServiceDiscovery 服务发现
type ServiceDiscovery struct {
	cli        *clientV3.Client  // etcd client
	serverList map[string]string // 服务列表
	lock       sync.RWMutex
}

var (
	Register *discovery.Resolver
)

func main() {
	config.InitConfig()

	endpoints := config.Conf.Etcd.Address
	Register = discovery.NewResolver(endpoints, log.Default())
	resolver.Register(Register)

	serverInfo := discovery.Server{
		Project: config.Conf.Project.Name, // 项目名
		Name:    config.Conf.Server.Name,  // 服务名
		Addr:    config.Conf.Server.Host,
		Version: config.Conf.Server.Version,
	}
	keyPrefix := discovery.BuildPrefix(serverInfo)
	// Scheme://Authority/Endpoint
	addr := fmt.Sprintf("%s://%s/%s", Register.Scheme(), "", keyPrefix)

	_, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return
	}

	// 监控系统型号，等待 ctrl + c 系统信号通知关闭
	exitCh := make(chan int, 1)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
		<-c
		Register.Close()
		time.Sleep(time.Second)
		close(exitCh)
	}()
	log.Println(fmt.Sprintf("exit %v", <-exitCh))

}

func main1() {
	var endpoints = []string{"xiaoyucc521.com:2379"}
	ser := NewServiceDiscovery(endpoints)
	defer func() {
		err := ser.Close()
		if err != nil {
			log.Println("关闭出错")
		}
	}()

	err := ser.WatchService("/server/")
	if err != nil {
		log.Fatal(err)
	}

	// 监控系统信号，等待 ctrl + c 系统信号通知服务关闭
	c := make(chan os.Signal, 1)
	go func() {
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	}()
	for {
		select {
		case <-time.Tick(10 * time.Second):
			log.Println(ser.GetServices())
		case <-c:
			log.Println("server discovery exit")
			return
		}
	}
}

// NewServiceDiscovery 新建服务发现
func NewServiceDiscovery(endpoints []string) *ServiceDiscovery {
	// 初始化etcd client
	cli, err := clientV3.New(clientV3.Config{
		Endpoints:   endpoints,
		DialTimeout: time.Duration(5) * time.Second,
	})
	if err != nil {
		log.Fatalln(err)
	}

	return &ServiceDiscovery{
		cli:        cli,
		serverList: make(map[string]string),
	}
}

// WatchService 初始化服务列表和监视
func (s *ServiceDiscovery) WatchService(prefix string) error {
	// 根据前缀获取现有的key
	resp, err := s.cli.Get(context.Background(), prefix, clientV3.WithPrefix())
	if err != nil {
		return err
	}

	// 遍历获取得到的k和v
	for _, ev := range resp.Kvs {
		s.SetServiceList(string(ev.Key), string(ev.Value))
	}

	// 监视前缀，修改变更server
	go s.watcher(prefix)
	return nil
}

// watcher 监听Key的前缀
func (s *ServiceDiscovery) watcher(prefix string) {
	rch := s.cli.Watch(context.Background(), prefix, clientV3.WithPrefix())
	log.Printf("watching prefix:%s now...", prefix)
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case mvccpb.PUT: // 修改或者新增
				s.SetServiceList(string(ev.Kv.Key), string(ev.Kv.Value))
			case mvccpb.DELETE: // 删除
				s.DelServiceList(string(ev.Kv.Key))
			}
		}
	}
}

func (s *ServiceDiscovery) SetServiceList(key, val string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.serverList[key] = string(val)
	log.Println("put key:", key, "val:", val)
}

func (s *ServiceDiscovery) DelServiceList(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.serverList, key)
	log.Println("Del key :", key)
}

// GetServices 获取服务地址
func (s *ServiceDiscovery) GetServices() []string {
	s.lock.RLock()
	defer s.lock.RUnlock()
	addrs := make([]string, 0, len(s.serverList))

	for _, v := range s.serverList {
		addrs = append(addrs, v)
	}
	return addrs
}

// Close 关闭服务
func (s *ServiceDiscovery) Close() error {
	return s.cli.Close()
}
