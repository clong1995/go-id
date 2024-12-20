package gid

import (
	"math/rand"
	"time"
)

var (
	chars = [256]rune{}
	maps  = make(map[rune]int64, 256)
)

func init() {
	base := []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!#$%&'()*+,-./:;<=>?@[]^_`{|}~¬±§¶©®™¥€£₹¢ªº∆∇∈∉∩∪∑∏∫√∂∞∝∧∨⊂⊃⊆⊇⊥⊤⊕⊗⊘⊙⊛⊝⊞⊟⊠⊡⊢⊣∅ℵαβγδεζηθικλμνξοπρστυφχψωΘΛΣΨΩℜ∐Ⅎ⅄ⅅⅆⅇⅈⅉ⅊⅋⅌⅍ⅎ⅏ⅠⅡⅢⅣⅤⅥⅦⅧⅨⅩⅪⅫⅬⅭⅮⅯⅰⅱⅲⅳⅴⅵⅶⅷⅸⅹⅺⅻⅼⅽⅾⅿ❍✔⚲⚶⊚⊜⊻⏠⏡⏦⏧⏼⏽⏿␀␁␂␃␄␅␆␇⏚⏛⏜⏝⏞⏟⏣␈␉␊␋␌␍␎␏␐␑")
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
	result := [9]rune{chars[xor]}
	index := 1
	for num > 0 {
		encodedByte := num&0xFF ^ xor
		result[index] = chars[encodedByte]
		num >>= 8
		index++
	}
	return string(result[:index])
}

// Decode 解码函数
func Decode(encoded string) (result int64) {
	runes := []rune(encoded)
	isXor := len(runes) > 8
	var xor int64
	if isXor {
		xor = maps[runes[0]]
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
	return result
}

// EncodeNoXor 编码
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

func xorKey() (key int64) {
	randSource := rand.NewSource(time.Now().UnixNano())
	randGen := rand.New(randSource)
	key = int64(randGen.Intn(256))
	return
}
