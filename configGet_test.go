package main

import (
	"fmt"
	"testing"
)

func TestGetHostInfo(t *testing.T) {
	t.Log("测试开始")
	data, _ := GetHostInfo()
	for _, v := range data {
		fmt.Println(v.hostname, v.ip)
	}
}
