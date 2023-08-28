package snowflake

import (
	"sync"
	"time"
)

const (
	machineBits  = int64(5)  // 机器id位数
	serviceBits  = int64(5)  // 服务id位数
	sequenceBits = int64(12) // 序列id位数

	maxMachineID  = int64(-1) ^ (int64(-1) << machineBits)  // 最大机器id 用于防止溢出
	maxServiceID  = int64(-1) ^ (int64(-1) << serviceBits)  // 最大服务id
	maxSequenceID = int64(-1) ^ (int64(-1) << sequenceBits) // 最大序列id

	timeLeft    = uint8(22) // 时间id向左移位的量
	machineLeft = uint8(17) // 机器id向左移位的量
	serviceLeft = uint8(12) // 服务id向左移位的量

	twepoch = int64(1667972427000) // 常量时间戳(毫秒),时间是: Wed Nov 9 13:40:27 CST 2022
)

type Worker struct {
	sync.Mutex       // 添加互斥锁，确保并发安全性
	lastStamp  int64 // 记录上一次生成ID的时间戳
	machineID  int64 // 机器id,0~31
	serviceID  int64 // 服务id,0~31
	sequenceID int64 // 当前毫秒已经生成的id序列号(从0开始累加) 1毫秒内最多生成4096个ID
}

// NewWorker 创建 worker 对象
// 分布式情况下,我们应通过外部配置文件或其他方式为每台机器分配独立的id
func NewWorker(machineID, serviceID int64) *Worker {
	return &Worker{
		lastStamp:  0,
		machineID:  machineID,
		serviceID:  serviceID,
		sequenceID: 0,
	}
}

// GenerateID 生成 ID
func (w *Worker) GenerateID() int64 {
	//多线程互斥
	w.Lock()
	defer w.Unlock()

	// 获取当前毫秒值
	mill := time.Now().UnixMilli()

	if mill == w.lastStamp {
		w.sequenceID = (w.sequenceID + 1) & maxSequenceID
		//当一个毫秒内分配的 id数>4096 个时，只能等待到下一毫秒去分配。
		if w.sequenceID == 0 {
			for mill > w.lastStamp {
				mill = time.Now().UnixMilli()
			}
		}
	} else {
		w.sequenceID = 0
	}

	w.lastStamp = mill

	id := (w.lastStamp-twepoch)<<timeLeft | w.machineID<<machineLeft | w.serviceID<<serviceLeft | w.sequenceID
	return id
}
