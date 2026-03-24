package gid

import (
	"math/rand/v2"
)

var (
	chars = [256]rune{}
	maps  = make(map[rune]int64, 256)
)

func shuffleBase() {
	base := []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyzАБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯабвгдеёжзийклмнопрстуфхцчшщъыьэюяΑΒΓΔΕΖΗΘΙΚΛΜΝΞΟΠΡΣΤΥΦΧΨΩαβγδεζηθικλμνξοπρστυφχψωàáâãäåāăąćĉċčçēĕėěèéêëęĝğġģĩīĭıìíîïįĺľŀłļńňŉñņōŏőòóôõöŕřŗśŝšşťŧţũūŭůűùúûüųŷÿýźżž")
	// Fisher-Yates 洗牌算法
	seed := uint64(epoch)
	r := rand.New(rand.NewPCG(seed, seed))
	for i := len(base) - 1; i > 0; i-- {
		j := r.IntN(i + 1)
		base[i], base[j] = base[j], base[i]
	}
	//
	for i, v := range base {
		chars[i] = v
		maps[v] = int64(i)
	}

	return
}

// Encode 编码
func Encode(num int64) string {
	xor := xorKey()
	result := [9]rune{chars[xor+1]}
	index := 1
	unum := uint64(num)
	for i := 0; i < 8; i++ {
		result[index] = chars[(unum&0xFF)^uint64(xor)]
		unum >>= 8
		index++
	}
	return string(result[:index])
}

// EncodeNoXor 非xor编码
func EncodeNoXor(num int64) string {
	var result [8]rune
	index := 0
	unum := uint64(num)
	if unum == 0 {
		return string([]rune{chars[0]})
	}
	for unum > 0 {
		result[index] = chars[unum&0xFF]
		unum >>= 8
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

	var result uint64

	if isXor {
		for _, char := range runes {
			result = result<<8 | uint64(maps[char]^xor)
		}
	} else {
		for _, char := range runes {
			result = result<<8 | uint64(maps[char])
		}
	}

	return int64(result)
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
