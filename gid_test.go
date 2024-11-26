package gid

import (
	"log"
	"testing"
)

var gid *Gid

func init() {
	var err error
	gid, err = NewId(1)
	if err != nil {
		log.Panicf("Error initializing Id: %v\n", err)
		return
	}
}

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
				timestamp: 1732503415000,
				machineID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := gid.Deterministic(tt.args.timestamp)
			if (err != nil) != tt.wantErr {
				t.Errorf("Deterministic() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
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
				id: 2501138846842884096,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTimestamp, gotMachineID, gotSequence := gid.Extract(tt.args.id)
			t.Logf("Extract() gotTimestamp = %v, gotMachineID %v, gotSequence %v", gotTimestamp, gotMachineID, gotSequence)
		})
	}
}

func TestSnowflake_Generate(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "生成唯一 ID",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := gid.Generate()
			t.Logf("Generate() = %v", got)
		})
	}
}
