package conn_pool

import (
	"net"
	"time"
)

var (
	MaxConnNum = 10 // 最大连接数
	MinConnNum = 2  // 空闲最小连接数
	MillTime   = 10 // 每隔10s清理一下空闲连接

	Protof = "tcp"
	Port   = "2000"
)

var (
	ConnPoolSize    = 0 // 连接池数目
	UseConnPoolSize = 0 // 使用中的连接数目
)

// ConnPool 连捷池对象，是一个双向链表
type ConnPool struct {
	Conn *net.Listener // 链接
	Pre  *ConnPool     // 前继
	Next *ConnPool     // 后驱
	Time int64         // 放入连接池时时间戳
}

var ConnPoolHead *ConnPool // 连接池队首
var ConnPoolTail *ConnPool // 连接池队尾

// GetConn 获取链接
func GetConn() *ConnPool {
	// 使用中的连接数和连接池的空闲连接数大于配置中的连接数
	if ConnPoolSize+UseConnPoolSize > MaxConnNum {
		return nil
	}
	// 使用连接数目加一
	UseConnPoolSize++
	// 如果连接池内有空闲连接，取连接池队首连接
	if ConnPoolHead != nil {
		temp := ConnPoolHead
		ConnPoolHead = ConnPoolHead.Next
		return temp
	}
	// 如果连接池空闲连接为空，创建一个连接返回
	return AddConn()
}

// AddConn 创建一个链接
func AddConn() *ConnPool {
	conn, _ := net.Listen(Protof, Port)
	ConPoolTemp := &ConnPool{
		Conn: &conn,
		Pre:  nil,
		Next: nil,
	}

	return ConPoolTemp
}

// PutConn 使用完链接放入连接池中，放入队首
func PutConn(connPool *ConnPool) {
	connPool.Time = time.Now().Unix()
	if ConnPoolHead == nil {
		ConnPoolHead = connPool
		ConnPoolTail = connPool
	}
	connPool.Next = ConnPoolHead
	ConnPoolHead.Pre = connPool
	ConnPoolHead = connPool
	ConnPoolSize++
}

// DeleteConn 每隔MillTime删除过期空闲链接，从队尾开始遍历，删除过期空闲连接
func DeleteConn() {
	timer := time.NewTimer(time.Duration(MillTime) * time.Second)
	for {
		timer.Reset(time.Duration(MillTime) * time.Second) // 这里复用了 timer
		select {
		case <-timer.C:
			nowTime := time.Now().Unix()
			//   每隔10秒执行一次，当空闲连接数达到配置最小值，结束
			for ConnPoolSize > MinConnNum {
				// 如果空闲链接超时
				if ConnPoolTail.Time+int64(MillTime) < nowTime {
					// 删除空闲链接
					ConnPoolTail = ConnPoolTail.Pre
					ConnPoolTail.Next = nil
					// 空闲链接数目减一
					ConnPoolSize--
				} else {
					break
				}
			}
		}
	}
}
