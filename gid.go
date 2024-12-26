package gid

import (
	"fmt"
	"github.com/clong1995/go-config"
	"log"
	"strconv"
	"sync"
	"time"
)

//id的结构
//| 46 bits - 时间戳部分 | 6 bits - 机器ID部分 | 12 bits - 序列号部分 |

const (
	epoch int64 = 1136185445000

	//timestampBits uint8 = 42
	machineBits  uint8 = 6
	sequenceBits uint8 = 12

	maxMachineID int64 = -1 ^ (-1 << machineBits)
	maxSequence  int64 = -1 ^ (-1 << sequenceBits)

	timestampShift = machineBits + sequenceBits
	machineShift   = sequenceBits
)

var id *gid

func init() {
	machineID := config.Value("MACHINE ID")
	if machineID == "" {
		log.Fatalln("MACHINE ID not found")
	}
	mid, err := strconv.ParseInt(machineID, 10, 64)
	if err != nil {
		log.Fatalln(err)
		return
	}

	if id, err = newId(mid); err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("gid created %s success\n", machineID)
}

// Gid 结构体
type gid struct {
	mu        sync.Mutex
	lastStamp int64
	sequence  int64
	machineID int64
}

// Generate 生成唯一ID
func Generate() int64 {
	id.mu.Lock()
	defer id.mu.Unlock()

	now := currentMillis()

	if now == id.lastStamp {
		id.sequence = (id.sequence + 1) & maxSequence
		if id.sequence == 0 {
			now = nextMillis(id.lastStamp)
		}
	} else {
		id.sequence = 0
	}

	id.lastStamp = now

	i := ((now - epoch) << timestampShift) |
		(id.machineID << machineShift) |
		id.sequence

	return i
}

// Extract 提取ID的时间戳、机器ID和序列号
func Extract(id int64) (timestamp int64, machineID int, sequence int64) {
	timestamp = id>>timestampShift + epoch
	machineID = int((id >> machineShift) & maxMachineID)
	sequence = id & maxSequence
	return timestamp, machineID, sequence
}

// Deterministic 直接生成特定时间和机器ID的ID
func Deterministic(timestamp int64) (int64, error) {
	if timestamp < epoch {
		return 0, fmt.Errorf("timestamp must be greater than or equal to the epoch: %d", epoch)
	}
	i := ((timestamp - epoch) << timestampShift) | (id.machineID << machineShift)
	return i, nil
}

func newId(machineID int64) (*gid, error) {
	if machineID < 0 || machineID > maxMachineID {
		return nil, fmt.Errorf("machine ID must be between 0 and %d", maxMachineID)
	}
	return &gid{
		lastStamp: 0,
		sequence:  0,
		machineID: machineID,
	}, nil
}

// 当前毫秒时间戳
func currentMillis() int64 {
	return time.Now().UnixMilli()
}

// 获取下一个时间戳
func nextMillis(lastStamp int64) int64 {
	timestamp := currentMillis()
	for timestamp <= lastStamp {
		timestamp = currentMillis()
	}
	return timestamp
}
