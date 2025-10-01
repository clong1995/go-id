package gid

import (
	"sync"
	"testing"
	"time"
)

func TestDeterministic(t *testing.T) {
	type args struct {
		timestamp int64
		machineID int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "直接生成特定时间和机器ID的ID",
			args: args{
				timestamp: 1734506945677,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Deterministic(tt.args.timestamp)
			t.Logf("Deterministic() got = %v", got)
		})
	}
}

func TestExtract(t *testing.T) {
	type args struct {
		id int64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "提取ID的时间戳、机器ID和序列号",
			args: args{
				id: 156846391473475584,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTimestamp, gotMachineID, gotSequence := Extract(tt.args.id)
			t.Logf("Extract() gotTimestamp = %v, gotMachineID %v, gotSequence %v", gotTimestamp, gotMachineID, gotSequence)
		})
	}
}

func TestID(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "生成唯一 ID",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ID()
			t.Logf("Generate() = %d", got)
		})
	}
}

func TestIDs(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "生成唯一 ID",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wg sync.WaitGroup
			for i := 0; i < 10; i++ {
				wg.Add(1)
				go func(i int, wg *sync.WaitGroup) {
					defer wg.Done()
					b := i == 7
					if b {
						time.Sleep(1 * time.Millisecond)
					}
					got := ID()
					t.Logf("ID() = %v => %v => %v", i, got, b)
				}(i, &wg)
			}
			wg.Wait()
		})
	}
}
