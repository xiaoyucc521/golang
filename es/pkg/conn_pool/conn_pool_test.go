package conn_pool

import "testing"

func Test(t *testing.T) {
	// 获取链接：从连接池那一个链接，没有就创建一个链接
	conn := GetConn()

	// TODO 执行代码
	Todo(conn)

	//该链接事件执行完毕，将该放入连接池队首
	PutConn(conn)

	//删除长时间不用的链接
	go DeleteConn()
}

// 连接后执行操作...
func Todo(conn *ConnPool) {

}
