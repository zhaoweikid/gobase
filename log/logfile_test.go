package gobase

import "testing"

func TestLog(t *testing.T) {
	log := New("stdout", DEBUG)

	log.Debug("debug haha...")
	log.Note("info haha...")
	log.Warn("warn haha...")
	log.Error("error haha...")
	//log.Fatal("fatal haha...")

}
