package gid

import (
	"sync/atomic"
	"time"

	"github.com/clong1995/go-ansi-color"
	"github.com/clong1995/go-config"
)

// id的结构
// | 46 bits - 时间戳部分 | 6 bits - 机器ID部分 | 12 bits - 序列号部分 |
const (
	//timestampBits uint8 = 46
	machineBits  uint8 = 6
	sequenceBits uint8 = 12

	maxMachineID int64 = -1 ^ (-1 << machineBits)
	maxSequence  int64 = -1 ^ (-1 << sequenceBits)

	timestampShift = machineBits + sequenceBits
	machineShift   = sequenceBits
)

var epoch int64
var id *gid

func init() {
	var prefix = "gid"
	//
	machineID, exists := config.Value[int64]("MACHINE ID")
	if !exists {
		pcolor.PrintFatal(prefix, "MACHINE ID not found")
		return
	}

	//
	epoch_, exists := config.Value[string]("EPOCH")
	if !exists {
		epoch_ = "2006-01-02 15:04:05"
	}

	t, err := time.ParseInLocation(time.DateTime, epoch_, time.UTC)
	if err != nil {
		pcolor.PrintFatal(prefix, err.Error())
		return
	}

	epoch = t.UnixMilli()

	if machineID < 0 || machineID > maxMachineID {
		pcolor.PrintFatal(prefix, "machine ID must be between 0 and %d", maxMachineID)
		return
	}

	id = newGid(machineID)

	num := ID()

	shuffleBase()
	pcolor.PrintSucc(prefix, "created %d success, %d : %s", machineID, num, Encode(num))
}

// gid 结构体
type gid struct {
	state     atomic.Uint64
	machineID int64
}

// ID 生成唯一ID
func ID() int64 {
	for {
		oldPacked := id.state.Load()
		oldStamp := int64(oldPacked >> sequenceBits)
		oldSeq := int64(oldPacked & uint64(maxSequence))
		now := time.Now().UnixMilli()

		var newStamp int64
		var newSeq int64

		if now < oldStamp {
			// 时钟回拨，等待
			time.Sleep(time.Duration(oldStamp-now) * time.Millisecond)
			continue
		}

		if now == oldStamp {
			newSeq = (oldSeq + 1) & maxSequence
			if newSeq == 0 {
				// 序列号溢出，等待下一毫秒
				time.Sleep(time.Millisecond)
				continue
			}
			newStamp = now
		} else {
			// 新的毫秒
			newStamp = now
			newSeq = 0
		}

		newPacked := (uint64(newStamp) << sequenceBits) | uint64(newSeq)

		if id.state.CompareAndSwap(oldPacked, newPacked) {
			return ((newStamp - epoch) << timestampShift) | (id.machineID << machineShift) | newSeq
		}
	}
}

// Extract 提取ID的时间戳、机器ID和序列号
func Extract(id int64) (timestamp int64, machineID int, sequence int64) {
	timestamp = id>>timestampShift + epoch
	machineID = int((id >> machineShift) & maxMachineID)
	sequence = id & maxSequence
	return timestamp, machineID, sequence
}

// Deterministic 直接生成特定时间和机器ID的ID，序列号是0，仅用于时间查询
func Deterministic(timestamp int64) int64 {
	i := ((timestamp - epoch) << timestampShift) | (id.machineID << machineShift)
	return i
}

func newGid(machineID int64) *gid {
	g := &gid{
		machineID: machineID,
	}
	g.state.Store(0)
	return g
}
