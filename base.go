package gid

import (
	"math/rand"
	"time"
)

var (
	chars = [256]rune{}
	maps  = make(map[rune]int64, 256)
)

func init() { //
	base := []rune("0a¢⋼⋴⊘b⊛⊝JKLMNOn=☊pCDoEFqlmGH⊙IST;☈UVW⌂₫⋄yzA3B45ef@çgκhcRdk≎⁏ℯ1⊟2ij6789⊞PQ⦙x°⁞৳X§Y⊠Z!∔$£₹ruv®wst⊡⊢⊣∅ℵαî≏↭()*៛,-ƒ¬±∉∩€β␏γδεζηθι™¥⊆⊇⊥⊤⊕␎␐¶©Ⅷ␁∪Ⅸℜ∐Ⅎ⅄ⅅⅆⅇⅈⅉ⊗πρσ∂∞τ␍␑υφξο∇∈∑∏∫√∝∧∨⊂Ⅳ␌⏛∆ⅤχψωΘΛΣ⚲Ψ⊃ΩⅡ␇⏞Ⅲλ⏦μνⅥ⏝⚶Ⅶ⏚⏟␉ⅾ⏼⏽␃ⅻⅼⅽⅿ❍✔⊜␊ⅶ␂ⅷ⅌⅍ⅸ⏜ⅹⅺⅬⅭⅫ␅Ⅾ␋⅊⅋ⅎ⅏ⅠⅩ⊚Ⅺ␄⏣␈ⅵⅯⅰ⏿␀␆ⅱⅲ⏡ⅳⅴ⊻⏠⏧Ł")
	for i, v := range base {
		chars[i] = v
		maps[v] = int64(i)
	}
}

// Encode 编码
func Encode(num int64) string {
	if num == 0 {
		return string(chars[0])
	}
	xor := xorKey()
	result := [9]rune{chars[xor+1]}
	index := 1
	for num > 0 {
		encodedByte := num&0xFF ^ xor
		result[index] = chars[encodedByte]
		num >>= 8
		index++
	}
	return string(result[:index])
}

// EncodeId 使用id加密
func EncodeId(num, uid int64) string {
	if num == 0 && uid == 0 {
		return string(chars[0])
	}
	return Encode(num + uid)
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
func Decode(encoded string) (result int64) {
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

	if isXor {
		for _, char := range runes {
			result = result<<8 | maps[char] ^ xor
		}
	} else {
		for _, char := range runes {
			result = result<<8 | maps[char]
		}
	}
	return
}

// DecodeId 使用id解密
func DecodeId(encoded string, uid int64) (result int64) {
	dId := Decode(encoded)
	//提取最高位
	//highestDigit := int64(float64(dId) / math.Pow10(int(math.Log10(float64(dId))))) // 用 log10 计算位数，再除以 10 的次幂
	result = dId - uid
	return
}

func xorKey() (key int64) {
	randSource := rand.NewSource(time.Now().UnixNano())
	randGen := rand.New(randSource)
	key = int64(randGen.Intn(255))
	return
}
