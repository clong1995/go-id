package gid

import (
	"math/rand/v2"
)

var (
	chars = [256]rune{}
	maps  = make(map[rune]int64, 256)
)

func init() { //
	base := []rune("0a¢⋼⋴⊘b⊛⊝JKLMNOnₔ☊pCDoEFqlmGH⊙IST;☈UVW⌂₫⋄yzA3B45ef@çgκhcRdk≎⁏ℯ1⊟2ij6789⊞PQ⦙x°⁞৳X§Y⊠Z!∔$£₹ruv®wst⊡⊢⊣∅ℵαî≏↭()*៛ₓ-ƒ¬±∉∩€β␏γδεζηθι™¥⊆⊇⊥⊤⊕␎␐¶©Ⅷ␁∪Ⅸℜ∐Ⅎ⅄ⅅⅆⅇⅈⅉ⊗πρσ∂∞τ␍␑υφξο∇∈∑∏∫√∝∧∨⊂Ⅳ␌⏛∆ⅤχψωΘΛΣ⚲Ψ⊃ΩⅡ␇⏞Ⅲλ⏦μνⅥ⏝⚶Ⅶ⏚⏟␉ⅾ⏼⏽␃ⅻⅼⅽⅿ❍✔⊜␊ⅶ␂ⅷ⅌⅍ⅸ⏜ⅹⅺⅬⅭⅫ␅Ⅾ␋⅊⅋ⅎ⅏ⅠⅩ⊚Ⅺ␄⏣␈ⅵⅯⅰ⏿␀␆ⅱⅲ⏡ⅳⅴ⊻⏠⏧Ł")
	for i, v := range base {
		chars[i] = v
		maps[v] = int64(i)
	}
}

// Encode 编码
func Encode(num int64) string {
	if num == 0 {
		return ""
	}
	xor := xorKey()
	result := [9]rune{chars[xor+1]}
	index := 1
	for num > 0 {
		result[index] = chars[num&0xFF^xor]
		num >>= 8
		index++
	}
	return string(result[:index])
}

// EncodeNoXor 非xor编码
func EncodeNoXor(num int64) string {
	if num == 0 {
		return string(chars[0])
	}
	var result [8]rune
	index := 0
	for num > 0 {
		result[index] = chars[num&0xFF]
		num >>= 8
		index++
	}
	return string(result[:index])
}

// Decode 解码
func Decode(encoded string) int64 {
	if encoded == "" {
		return 0
	}

	runes := []rune(encoded)
	isXor := len(runes) > 8
	var xor int64
	if isXor {
		xor = maps[runes[0]] - 1
		runes = runes[1:]
	}

	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	var result int64

	if isXor {
		for _, char := range runes {
			result = result<<8 | maps[char] ^ xor
		}
	} else {
		for _, char := range runes {
			result = result<<8 | maps[char]
		}
	}

	return result
}

func Union(num int64, salt ...int64) string {
	var s int64
	if len(salt) != 0 && salt[0] != 0 {
		s = salt[0]
	}
	return EncodeNoXor(num + s)
}

func xorKey() int64 {
	return rand.Int64N(255)
}
