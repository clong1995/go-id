package gid

import "testing"

func TestEncode(t *testing.T) {
	type args struct {
		num uint64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "编码",
			args: args{
				num: 156846391473475584,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Encode(tt.args.num)
			t.Logf("Encode() = %v", got)
		})
	}
}

func TestDecode(t *testing.T) {
	type args struct {
		encoded string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "解码",
			args: args{
				encoded: "2JW␈,kg0",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Decode(tt.args.encoded)
			t.Logf("Decode() = %v", got)
		})
	}
}
