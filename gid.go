package gid

import (
	"fmt"
	"sync"
	"time"
)

//id的结构
//| 42 bits - 时间戳部分 | 10 bits - 机器 ID 部分 | 12 bits - 序列号部分 |

const (
	epoch int64 = 1136185445000

	//timestampBits uint8 = 42
	machineBits  uint8 = 10
	sequenceBits uint8 = 12

	maxMachineID int64 = -1 ^ (-1 << machineBits)
	maxSequence  int64 = -1 ^ (-1 << sequenceBits)

	timestampShift = machineBits + sequenceBits
	machineShift   = sequenceBits
)

// Gid 结构体
type Gid struct {
	mu        sync.Mutex
	lastStamp int64
	sequence  int64
	machineID int64
}

// NewId 初始化
func NewId(machineID int64) (*Gid, error) {
	if machineID < 0 || machineID > maxMachineID {
		return nil, fmt.Errorf("machine ID must be between 0 and %d", maxMachineID)
	}
	return &Gid{
		lastStamp: 0,
		sequence:  0,
		machineID: machineID,
	}, nil
}

// Generate 生成唯一 ID
func (s *Gid) Generate() uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := currentMillis()

	if now == s.lastStamp {
		s.sequence = (s.sequence + 1) & maxSequence
		if s.sequence == 0 {
			now = nextMillis(s.lastStamp)
		}
	} else {
		s.sequence = 0
	}

	s.lastStamp = now

	id := ((now - epoch) << timestampShift) |
		(s.machineID << machineShift) |
		s.sequence

	return uint64(id)
}

// Extract 提取ID的时间戳、机器ID和序列号
func (s *Gid) Extract(id int64) (timestamp int64, machineID int64, sequence int64) {
	timestamp = (id >> timestampShift) + epoch
	machineID = (id >> machineShift) & maxMachineID
	sequence = id & maxSequence
	return timestamp, machineID, sequence
}

// Deterministic 直接生成特定时间和机器ID的ID
func (s *Gid) Deterministic(timestamp int64) (uint64, error) {
	if timestamp < epoch {
		return 0, fmt.Errorf("timestamp must be greater than or equal to the epoch: %d", epoch)
	}
	id := ((timestamp - epoch) << timestampShift) | (s.machineID << machineShift)
	return uint64(id), nil
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
