package newIdGen

import (
	"errors"
	"fmt"
	"gopkg.in/go-playground/pool.v3"
	"time"
)

var Epoch int64 = 1551698559000
var timestampPos int64 = 0x3fffffffffe00000

type idNode struct {
	lockChan  chan bool //并发安全用互斥锁
	timestamp int64     //时间戳部分
	node      int32     //节点id
	step      int64     //序列id
}

type nodePool struct {
	nodeCnt int32
	nodes   map[int32]*idNode

	workQueue chan int

	nodeType uint8 //类型位

	nodeBits  uint8 //节点id位数
	stepBits  uint8 //序列号位数
	nodeMax   int64 //节点id的最大值（512）
	stepMax   int64 //序列号的最大值
	timeShift uint8 //时间戳向左的偏移量
	nodeShift uint8 //节点id向左的偏移量

	workerPool pool.Pool
}

var midPool nodePool

func InitMidPool(cnt int32) error {
	if cnt <= 0 {
		return errors.New("invalid node count")
	}
	midNodes := make(map[int32]*idNode, cnt)
	var i int32
	for i = 1; i <= cnt; i++ {
		midNodes[i] = &idNode{
			lockChan:  make(chan bool, 1),
			timestamp: 0,
			node:      i,
			step:      1,
		}
	}
	for i = 1; i <= cnt; i++ {
		midNodes[i].lockChan <- true
	}
	fmt.Println("cnt", uint(cnt))
	midPool = nodePool{
		nodeCnt:    cnt,
		nodeType:   1,
		nodes:      midNodes,
		nodeBits:   9,
		stepBits:   12,
		nodeMax:    511,
		stepMax:    2 ^ 12 - 1,
		timeShift:  21,
		nodeShift:  12,
		workerPool: pool.NewLimited(uint(cnt)),
	}
	return nil
}

func GenMidId() (mid int64, err error) {

	tempGen := midPool.workerPool.Queue(genId())
	tempGen.Wait()
	if err = tempGen.Error(); err != nil {
		return 0, errors.New("gen mid fail")
	}

	mid = tempGen.Value().(int64)
	return mid, nil

	//if midPool.nodeCnt == 0 {
	//	return 0, errors.New("no valid gen mid pool")
	//}
	//tryTime := 0
	//var i int32
	//for true {
	//	for i = 1; i <= midPool.nodeCnt; i++ {
	//		mid, err := midPool.genId(midPool.nodes[i])
	//		if err == nil {
	//			return mid, nil
	//		}
	//	}
	//	tryTime++
	//	if tryTime > 1000 {
	//		return 0, errors.New("gen mid pool is busy")
	//	}
	//}
	//return 0, errors.New("fail to gen mid")
}

func genId() pool.WorkFunc {
	return func(wu pool.WorkUnit) (interface{}, error) {
		if wu.IsCancelled() {
			// return values not used
			return nil, nil
		}
		//找到当前空闲的节点
		var i int32
		for i = 1; i <= midPool.nodeCnt; i++ {
			n := midPool.nodes[i]
			select {
			case <-midPool.nodes[i].lockChan:
				fmt.Println(i)
				defer func() {
					midPool.nodes[i].lockChan <- true
				}()
				now := time.Now().UnixNano() / 1e6

				if n.timestamp == now {
					//step 步进1
					n.step++

					//当前step用完
					if n.step > midPool.stepMax {
						//等待本毫秒结束
						for now <= n.timestamp {
							now = time.Now().UnixNano() / 1e6
						}
						n.step = 0
					}
				} else {
					// 本毫秒内step用完
					n.step = 0
				}

				n.timestamp = now
				//移位运算
				result := int64(int64(midPool.nodeType)<<(midPool.timeShift+41) | (now-Epoch)<<midPool.timeShift | (int64(n.node) << midPool.nodeShift) | int64(n.step))
				return result, nil
			}
		}
		return 0, errors.New("no valid node")
	}
}

func (p *nodePool) genId(n *idNode) (id int64, err error) {
	select {
	case <-n.lockChan:
		defer func() {
			n.lockChan <- true
		}()
		now := time.Now().UnixNano() / 1e6

		if n.timestamp == now {
			//step 步进1
			n.step++

			//当前step用完
			if n.step > p.stepMax {
				//等待本毫秒结束
				for now <= n.timestamp {
					now = time.Now().UnixNano() / 1e6
				}
				n.step = 0
			}
		} else {
			// 本毫秒内step用完
			n.step = 0
		}

		n.timestamp = now
		//移位运算
		result := int64(int64(p.nodeType)<<(p.timeShift+41) | (now-Epoch)<<p.timeShift | (int64(n.node) << p.nodeShift) | int64(n.step))
		return result, nil
	default:
		return 0, errors.New("this node is busy")
	}
}
