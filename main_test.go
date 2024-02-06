package main

import (
	"testing"
)

func Test001(t *testing.T) {
	t.Log(123)
}

func Test002(t *testing.T) {
	t.Logf("%q", ParseCMD("RDS DM1001.H 2\r\n"))
}
func Test003(t *testing.T) {
	t.Logf("%q", ParseCMD("RDS DM2001.H 7\r\n"))
}
