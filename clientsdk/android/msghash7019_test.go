package android

import (
	"github.com/lunny/log"
	"testing"
)

func TestCalcMsgCrcForData_7019(t *testing.T) {
	string_7019 := CalcMsgCrcForString_7019("1452087f5e5b0ca9064fdf6aed88e9bb")
	log.Println(string_7019)
}
