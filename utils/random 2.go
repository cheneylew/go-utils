package utils

import (
	"math/rand"
	"time"
	"fmt"
)

func JKRandInt(max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max)
}

func JKRandStr(length int) string {
	lowerLetter := "qwertyuiopasdfghjklzxcvbnm"
	upperLetter := "QWERTYUIOPASDFGHJKLZXCVBNM"
	number := "0123456789"
	pool := fmt.Sprintf("%v%v%v", lowerLetter, upperLetter, number)
	result := ""
	for i := 0; i < length; i++ {
		index := JKRandInt(len(pool))
		result = fmt.Sprintf("%s%c",result,pool[index])
	}

	return result
}