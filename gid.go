package gid

import (
	"strconv"
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
	machineID := config.Value("MACHINE ID")
	if machineID == "" {
		pcolor.PrintFatal(prefix, "MACHINE ID not found")
		return
	}
	mid, err := strconv.ParseInt(machineID, 10, 64)
	if err != nil {
		pcolor.PrintFatal(prefix, err.Error())
		return
	}

	//
	epoch_ := config.Value("EPOCH")
	if epoch_ == "" {
		epoch_ = "2006-01-02 15:04:05"
	}

	t, err := time.ParseInLocation(time.DateTime, epoch_, time.Local)
	if err != nil {
		pcolor.PrintFatal(prefix, err.Error())
		return
	}

	epoch = t.UnixMilli()

	if mid < 0 || mid > maxMachineID {
		pcolor.PrintFatal(prefix, "machine ID must be between 0 and %d", maxMachineID)
		return
	}

	id = newGid(mid)

	num := ID()

	shuffleBase()
	pcolor.PrintSucc(prefix, "created %d success, %d : %s", mid, num, Encode(num))
}

// state holds the atomic state of the generator.
type state struct {
	lastStamp int64
	sequence  int64
}

// gid 结构体
type gid struct {
	state     atomic.Pointer[state]
	machineID int64
}

// ID 生成唯一ID
func ID() int64 {
	for {
		oldState := id.state.Load()
		now := time.Now().UnixMilli()

		var newState state

		if now < oldState.lastStamp {
			// 时钟回拨，等待
			time.Sleep(time.Duration(oldState.lastStamp-now) * time.Millisecond)
			continue
		}

		if now == oldState.lastStamp {
			newSequence := (oldState.sequence + 1) & maxSequence
			if newSequence == 0 {
				// 序列号溢出，等待下一毫秒
				time.Sleep(time.Millisecond)
				continue
			}
			newState = state{lastStamp: now, sequence: newSequence}
		} else {
			// 新的毫秒
			newState = state{lastStamp: now, sequence: 0}
		}

		if id.state.CompareAndSwap(oldState, &newState) {
			return ((newState.lastStamp - epoch) << timestampShift) | (id.machineID << machineShift) | newState.sequence
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
	g.state.Store(&state{lastStamp: 0, sequence: 0})
	return g
}
