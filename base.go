package gid

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

	for i, j := 0, index-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return string(result[:index])
}

// Decode 解码函数
func Decode(encoded string) uint64 {
	var result uint64
	for _, char := range encoded {
		result = result<<8 | maps[char]
	}
	return result
}
