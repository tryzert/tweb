package tool

import (
	"math/rand"
	"time"
)

//用于产生一个给定长度的随机字符串 userPassKey
func CreateRandPassKey(length int) string {
	return randString(length)
}

//给定长度的随机字符串生成器
func randString(length int) string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^&*-+"
	res := make([]byte, length)
	for i := 0; i < length; i++ {
		res[i] = str[r.Intn(len(str))]
	}
	return string(res)
}
