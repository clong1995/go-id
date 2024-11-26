package gid

import "strings"

var (
	chars = [256]rune{}
	maps  = make(map[rune]uint64, 256)
)

func init() {
	base := []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!#$%&'()*+,-./:;<=>?@[]^_`{|}~¬±§¶©®™¥€£₹¢ªº∆∇∈∉∩∪∑∏∫√∂∞∝∧∨⊂⊃⊆⊇⊥⊤⊕⊗⊘⊙⊛⊝⊞⊟⊠⊡⊢⊣∅ℵαβγδεζηθικλμνξοπρστυφχψωΘΛΣΨΩℜ∐Ⅎ⅄ⅅⅆⅇⅈⅉ⅊⅋⅌⅍ⅎ⅏ⅠⅡⅢⅣⅤⅥⅦⅧⅨⅩⅪⅫⅬⅭⅮⅯⅰⅱⅲⅳⅴⅵⅶⅷⅸⅹⅺⅻⅼⅽⅾⅿ❍✔⚲⚶⊚⊜⊻⏠⏡⏦⏧⏼⏽⏿␀␁␂␃␄␅␆␇⏚⏛⏜⏝⏞⏟⏣␈␉␊␋␌␍␎␏␐␑")
	for i, v := range base {
		chars[i] = v
		maps[v] = uint64(i)
	}
}

// Encode 编码
func Encode(num uint64) string {
	var builder strings.Builder
	for i := 7; i >= 0; i-- {
		byteValue := byte((num >> (i * 8)) & 0xFF)
		builder.WriteRune(chars[byteValue])
	}
	return builder.String()
}

// Decode 解码函数
func Decode(encoded string) uint64 {
	var result uint64
	runes := []rune(encoded)
	for i, char := range runes {
		byteValue := maps[char]
		result |= byteValue << (uint64(7-i) * 8)
	}
	return result
}
