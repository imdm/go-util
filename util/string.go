package util

import (
	"fmt"
	"net/url"
	"strconv"
)

func StringToInt32(arg string) int32 {
	i, err := strconv.ParseInt(arg, 10, 32)
	if err != nil {
		return -1
	}
	return int32(i)
}

func StringToInt64(arg string) int64 {
	i, err := strconv.ParseInt(arg, 10, 64)
	if err != nil {
		return -1
	}
	return i
}

func StringToFloat64(arg string) float64 {
	i, err := strconv.ParseFloat(arg, 64)
	if err != nil {
		return 0
	}
	return i
}

func URLDecode(arg string) string {
	re, err := url.QueryUnescape(arg)
	if err != nil {
		return ""
	}
	return re
}

func toStr(value interface{}) (s string) {
	switch v := value.(type) {
	case bool:
		s = strconv.FormatBool(v)
	case float32:
		s = strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		s = strconv.FormatFloat(v, 'f', -1, 64)
	case int:
		s = strconv.FormatInt(int64(v), 10)
	case int8:
		s = strconv.FormatInt(int64(v), 10)
	case int16:
		s = strconv.FormatInt(int64(v), 10)
	case int32:
		s = strconv.FormatInt(int64(v), 10)
	case int64:
		s = strconv.FormatInt(v, 10)
	case uint:
		s = strconv.FormatUint(uint64(v), 10)
	case uint8:
		s = strconv.FormatUint(uint64(v), 10)
	case uint16:
		s = strconv.FormatUint(uint64(v), 10)
	case uint32:
		s = strconv.FormatUint(uint64(v), 10)
	case uint64:
		s = strconv.FormatUint(v, 10)
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		s = fmt.Sprintf("%v", v)
	}
	return s
}

//截取字符串 start 起点下标 length 需要截取的长度
func Substr(str string, start int, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0
	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length
	if start > end {
		start, end = end, start
	}
	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}
	return string(rs[start:end])
}

var digitStringMap = map[int32]string{
	1: "一",
	2: "二",
	3: "三",
	4: "四",
	5: "五",
	6: "六",
	7: "七",
	8: "八",
	9: "九",
	0: "零",
}

func DigitToString(digit int32) string {
	if str, ok := digitStringMap[digit]; ok {
		return str
	}
	return strconv.Itoa(int(digit))
}
