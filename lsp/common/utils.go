package common

import (
	"unsafe"
)

// 这儿的两个函数有些危险: string 是不可变的，[]byte 是可变的。

// BytesToString 实现 []byte 转换成 string, 不需要额外的内存分配
func BytesToString(bytes []byte) string {
	return unsafe.String(unsafe.SliceData(bytes), len(bytes))
}

// StringToBytes 实现string 转换成 []byte, 不用额外的内存分配
func StringToBytes(str string) (bytes []byte) {
	return unsafe.Slice(unsafe.StringData(str), len(str))
}

func IsDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func IsHexChar(c byte) bool {
	return IsDigit(c) || c >= 'a' && c <= 'f' || c >= 'A' && c <= 'F'
}

func IsLetterChar(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

func IsNameChar(c byte) bool {
	return IsLetterChar(c) || IsDigit(c) || c == '_'
}
