package gobase

import (
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {
	conf := NewConfig("test.conf")

	if conf != nil {
		for sec, sv := range conf.Section {
			for k, v := range sv {
				fmt.Printf("sec:%s %s=%s\n", sec, k, v.AsString())
			}
		}
	}

}
