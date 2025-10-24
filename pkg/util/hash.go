package util

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
	"hash/crc32"
)

// HashCrc32 通过字符串获取32位数字
func HashCrc32(str string) uint32 {

	return crc32.ChecksumIEEE([]byte(str))
}

func Md5Hmac(key, seg string) string {
	mac := hmac.New(md5.New, []byte(key))
	mac.Write([]byte(seg))
	return hex.EncodeToString(mac.Sum(nil))
}

func Md5(plain string) string {
	h := md5.New()
	h.Write([]byte(plain))
	return hex.EncodeToString(h.Sum(nil))
}
