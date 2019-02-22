package idGenerate

/*Twitter snowflake ID实现
 *不重复ID一共64位
 *首位保留1+类型位1(0:uid/1:mid)+时间戳位41(精确到ms，即时间戳*1e3)+工作机器id9+序列号12
 */

import (
	"errors"
	"time"
)

var Epoch int64 = 1551698559000
var timestampPos int64 = 0x3fffffffffe00000

type nodePool struct {
	nodeCnt int32
	nodes   map[int32]*idNode

	nodeType uint8 //类型位

	nodeBits  uint8 //节点id位数
	stepBits  uint8 //序列号位数
	nodeMax   int64 //节点id的最大值（512）
	stepMax   int64 //序列号的最大值
	timeShift uint8 //时间戳向左的偏移量
	nodeShift uint8 //节点id向左的偏移量
}

var midPool nodePool
var uidPool nodePool

type idNode struct {
	lockChan  chan bool //并发安全用互斥锁
	timestamp int64     //时间戳部分
	node      int32     //节点id
	step      int64     //序列id
}

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
	midPool = nodePool{
		nodeCnt:   cnt,
		nodeType:  1,
		nodes:     midNodes,
		nodeBits:  9,
		stepBits:  12,
		nodeMax:   511,
		stepMax:   2 ^ 12 - 1,
		timeShift: 21,
		nodeShift: 12,
	}
	return nil
}

func InitUidPool(cnt int32) error {
	if cnt <= 0 {
		return errors.New("invalid node count")
	}
	uidNodes := make(map[int32]*idNode, cnt)
	var i int32
	for i = 1; i <= cnt; i++ {
		uidNodes[i] = &idNode{
			lockChan:  make(chan bool, 1),
			timestamp: 0,
			node:      i,
			step:      1,
		}
	}
	for i = 1; i <= cnt; i++ {
		uidNodes[i].lockChan <- true
	}
	uidPool = nodePool{
		nodeCnt:   cnt,
		nodeType:  0,
		nodes:     uidNodes,
		nodeBits:  9,
		stepBits:  12,
		nodeMax:   511,
		stepMax:   2 ^ 12 - 1,
		timeShift: 21,
		nodeShift: 12,
	}
	return nil
}

func GenMidId() (mid int64, err error) {
	if midPool.nodeCnt == 0 {
		return 0, errors.New("no valid gen mid pool")
	}
	tryTime := 0
	var i int32
	for true {
		for i = 1; i <= midPool.nodeCnt; i++ {
			mid, err := midPool.genId(midPool.nodes[i])
			if err == nil {
				return mid, nil
			}
		}
		tryTime++
		if tryTime > 10 {
			return 0, errors.New("gen mid pool is busy")
		}
	}
	return 0, errors.New("fail to gen mid")
}

func GenUidId() (uid int64, err error) {
	if uidPool.nodeCnt == 0 {
		return 0, errors.New("no valid gen uid pool")
	}
	tryTime := 0
	var i int32
	for true {
		for i = 1; i <= uidPool.nodeCnt; i++ {
			uid, err := uidPool.genId(uidPool.nodes[i])
			if err == nil {
				return uid, nil
			}
		}
		tryTime++
		if tryTime > 100 {
			return 0, errors.New("gen uid pool is busy")
		}
	}
	return 0, errors.New("fail to gen uid")
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

func GetTimeInfoById(id int64) time.Time {
	mblogTime := (id&timestampPos)>>21
	realTime := (mblogTime + Epoch)/1e3
	return time.Unix(realTime,0)
}
