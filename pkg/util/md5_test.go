package util

import "testing"

func TestMD5(t *testing.T) {
	println(MD5(MD5("123456")))
}
