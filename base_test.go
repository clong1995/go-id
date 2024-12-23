package gid

import "testing"

func TestEncode(t *testing.T) {
	type args struct {
		num int64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "编码",
			args: args{
				num: 156866482610417664,
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
				encoded: "Eo∑_⊝`rcF",
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

func TestEncodeNoXor(t *testing.T) {
	type args struct {
		num int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "编码",
			args: args{
				num: 156867655266209796,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EncodeNoXor(tt.args.num)
			t.Logf("EncodeXor() = %v", got)
		})
	}
}
