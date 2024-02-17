package PoliteDog

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

// 字符串拼接
func joinStrings(n int, strs ...string) string {
	arr := make([]string, n)
	for i, str := range strs {
		arr[i] = str
	}

	return strings.Join(arr, "")
}

// MD5编码
func md5Encode(data []byte) string {
	h := md5.New()
	h.Write(data)
	r := h.Sum(nil)
	
	return hex.EncodeToString(r)
}
