package tool

import (
	"math/rand"
	"time"
)

func CreateRandPassKey(length int) string {
	return randString(length)
}


func randString(length int) string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^&*-+"
	res := make([]byte, length)
	for i := 0; i < length; i++ {
		res[i] = str[r.Intn(len(str))]
	}
	return string(res)
}