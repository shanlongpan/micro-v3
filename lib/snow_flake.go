/**
* @Author:Tristan
* @Date: 2021/12/30 8:13 下午
 */

package lib

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

/**
 *   64bit, 前41 bit 毫秒级别时间戳 最大 2199023255552， 中间 10 bit 工作机 id,最后12bit-序列号
 *   41位 2^{41}-1 个毫秒的值，转化成单位年则是 (2^{41}-1) / (1000  60  60  24  365) = 69 年
 *   10位，用来记录工作机器id。
 *		- 可以部署在 2^{10} = 1024 个节点 0- 1023
 *   后12位 序列号，用来记录同毫秒内产生的不同id。
 *   12位（bit）可以表示的最大正整数是 2^{12}-1 = 4095
 */
const (
	machineIDBits = uint64(10)
	sequenceBits  = uint64(12)

	maxMachineID = int64(-1) ^ (int64(-1) << machineIDBits) //节点ID的最大值 用于防止溢出
	maxSequence  = int64(-1) ^ (int64(-1) << sequenceBits)

	timeLeft    = uint8(machineIDBits + sequenceBits) // timeLeft = workerIDBits + sequenceBits // 时间戳向左偏移量 22
	machineLeft = uint8(sequenceBits)                 // workLeft = sequenceBits // 节点IDx向左偏移量 12
	// 2021-12-01 00:00:00 +0800 CST
	twepoch = int64(1638288000000) // 常量时间戳(毫秒)
)

const maskSequence = uint64(1<<sequenceBits - 1)
const maskMachineID = uint64((1<<machineIDBits - 1) << sequenceBits)

// Decompose returns a set of Snowflake ID parts.
func Decompose(id uint64) map[string]uint64 {

	msb := id >> 63
	timer := id >> (sequenceBits + machineIDBits)
	machineID := id & maskMachineID >> sequenceBits
	sequence := id & maskSequence
	return map[string]uint64{
		"id":         id,
		"msb":        msb,
		"time":       timer + uint64(twepoch),
		"sequence":   sequence,
		"machine-id": machineID,
	}
}

//分布式情况下,我们应通过外部配置文件或其他方式为每台机器分配独立的id
func NewSnowFlake(machineID int64) (*SnowFlake, error) {
	if machineID > maxMachineID {
		return nil, errors.New(fmt.Sprintf("%d 大于最大机器id%d", machineID, maxMachineID))
	}

	return &SnowFlake{
		MachineID: machineID,
		LastStamp: 0,
		Sequence:  0,
	}, nil
}

type SnowFlake struct {
	mu        sync.Mutex
	LastStamp int64 // 记录上一次ID的时间戳
	MachineID int64 // 该节点的ID  工作机器ID(0~31)
	Sequence  int64 // 当前毫秒已经生成的ID序列号(从0 开始累加) 1毫秒内最多生成4095个ID
}

// 毫秒
func (w *SnowFlake) getMilliSeconds() int64 {
	return time.Now().UnixNano() / 1e6
}

func (w *SnowFlake) NextID() (uint64, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	return w.nextID()
}

func (w *SnowFlake) nextID() (uint64, error) {
	timeStamp := w.getMilliSeconds()
	if timeStamp < w.LastStamp {
		return 0, errors.New("time is moving backwards,waiting until")
	}

	if w.LastStamp == timeStamp {

		w.Sequence = (w.Sequence + 1) & maxSequence

		if w.Sequence == 0 {
			for timeStamp <= w.LastStamp {
				timeStamp = w.getMilliSeconds()
			}
		}
	} else {
		w.Sequence = 0
	}

	w.LastStamp = timeStamp
	id := ((timeStamp - twepoch) << timeLeft) |
		(w.MachineID << machineLeft) |
		w.Sequence

	return uint64(id), nil
}


func privateIPv4() (net.IP, error) {
	as, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, a := range as {
		ipnet, ok := a.(*net.IPNet)
		if !ok || ipnet.IP.IsLoopback() {
			continue
		}

		ip := ipnet.IP.To4()
		if isPrivateIPv4(ip) {
			return ip, nil
		}
	}
	return nil, errors.New("no private ip address")
}

func isPrivateIPv4(ip net.IP) bool {
	return ip != nil &&
		(ip[0] == 10 || ip[0] == 172 && (ip[1] >= 16 && ip[1] < 32) || ip[0] == 192 && ip[1] == 168)
}

//获取后16位 局域网 私有地址 192.168.0.0--192.168.255.255 ,可以作为 machineId，注意不要超过1023，也可以自定义，
func lower16BitPrivateIP() (uint16, error) {
	ip, err := privateIPv4()
	if err != nil {
		return 0, err
	}
	return uint16(ip[2])<<8 + uint16(ip[3]), nil
}
