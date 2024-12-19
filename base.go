package gid

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
	var result [8]rune
	index := 0
	for num > 0 {
		result[index] = chars[num&0xFF]
		num >>= 8
		index++
	}
	return string(result[:index])
}

// Decode 解码函数
func Decode(encoded string) (result int64) {
	runes := []rune(encoded)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	for _, char := range runes {
		result = result<<8 | maps[char]
	}
	return result
}
