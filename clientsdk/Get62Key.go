package clientsdk

import (
	"encoding/hex"
	"feiyu.com/wx/clientsdk/baseutils"
	math_rand "math/rand"
	"strings"
)

func Get62Key(Key string) string {
	if len(Key) < 344 {
		return baseutils.MD5ToLower(RandSeq(15))
	}
	start := strings.Index(strings.ToUpper(Key), "6E756C6C5F1020") + len("6E756C6C5F1020")
	m, _ := hex.DecodeString(Key[start : start+64])
	return string(m)
}

func RandSeq(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[math_rand.Intn(len(letters))]
	}
	return string(b)
}
