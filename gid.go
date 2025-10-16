package gid

import (
	"strconv"
	"sync/atomic"
	"time"

	pcolor "github.com/clong1995/go-ansi-color"
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

	maxBackoff = 5 * time.Second
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
		pcolor.PrintFatal(prefix, "EPOCH not found")
		return
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

	id = newId(mid)

	num := ID()

	shuffleBase()
	pcolor.PrintSucc(prefix, "created %d success, %d : %s", mid, num, Encode(num))
}

// Gid 结构体
type gid struct {
	lastStamp int64
	sequence  int64
	machineID int64
}

// ID 生成唯一ID
func ID() int64 {
	for {
		now := currentMillis()
		lastStamp := atomic.LoadInt64(&id.lastStamp)

		if now < lastStamp {
			// 时间回拨
			now = nextMillis(lastStamp)
			continue
		}

		seq := atomic.LoadInt64(&id.sequence)
		var newSeq int64
		if now == lastStamp {
			newSeq = (seq + 1) & maxSequence
			if newSeq == 0 {
				// 序列号溢出，等待下一毫秒
				now = nextMillis(lastStamp)
				continue
			}
		}

		if atomic.CompareAndSwapInt64(&id.lastStamp, lastStamp, now) &&
			atomic.CompareAndSwapInt64(&id.sequence, seq, newSeq) {
			return ((now - epoch) << timestampShift) | (id.machineID << machineShift) | newSeq
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

// Deterministic 直接生成特定时间和机器ID的ID
func Deterministic(timestamp int64) int64 {
	i := ((timestamp) << timestampShift) | (id.machineID << machineShift)
	return i
}

func newId(machineID int64) *gid {
	return &gid{
		lastStamp: 0,
		sequence:  0,
		machineID: machineID,
	}
}

// 当前毫秒时间戳
func currentMillis() int64 {
	return time.Now().UnixMilli()
}

// 获取下一个时间戳
func nextMillis(lastStamp int64) int64 {
	for {
		now := currentMillis()
		if now > lastStamp {
			return now
		}
		// 计算需要等待的时间
		waitTime := time.Duration(lastStamp-now+1) * time.Millisecond
		if waitTime > maxBackoff {
			time.Sleep(maxBackoff)
		} else {
			time.Sleep(waitTime)
		}
	}
}
