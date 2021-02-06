package gobase

import (
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {
	conf := NewConfig("test.conf")
	fmt.Printf("NewConfig: %p\n", conf)
	sec := conf.Section
	if conf != nil {
		for sec, sv := range conf.Section {
			for k, v := range sv {
				//var v2 *ConfigValue
				//v2 = &v
				//fmt.Printf("Value: %p %p\n", &v, v2)
				fmt.Printf("sec1:%s %s=%s\n", sec, k, v.AsString())
				//fmt.Printf("sec2:%s %s=%s\n", sec, k, v2.AsString())
			}
		}
	}

	var p map[string]ConfigValue
	p = sec["default"]

	var c1, c2, c3 ConfigValue

	c1 = p["ip"]
	c2 = p["port"]
	c3 = p["daemon"]

	fmt.Printf("ip:%s port:%d daemon:%v\n", 
		c1.AsString(),
		c2.AsInt(0),
		c3.AsBool(false))

}
