package gobase

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"
)

//type ConfigSection map[string]map[string]string

type ConfigValue struct {
	value string
}
type ConfigSection struct {
	Section map[string]map[string]ConfigValue
}

func NewConfig(filename string) *ConfigSection {
	f, err := os.Open(filename)
	if err != nil {
		return nil
	}

	defer f.Close()

	conf := new(ConfigSection)
	conf.Section = make(map[string]map[string]ConfigValue)
	c := *conf
	reader := bufio.NewReader(f)
	var sec string
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return conf
		}
		if line[0] == '#' {
			continue
		}
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		if line[0] == '[' {
			sec = line[1 : len(line)-1]
			c.Section[sec] = make(map[string]ConfigValue)
		}else{
			key, value := parseLine(line)
			c.Section[sec][key] = *NewConfigValue(value)
		}
	}

	return conf
}

func parseLine(line string) (string, string) {
	ln := strings.TrimSpace(line)
	pair := strings.Split(ln, "=")
	if len(pair) != 2 {
		Warn("config parse error at: %s\n", line)
		os.Exit(1)
	}

	key := strings.TrimSpace(pair[0])
	value := strings.TrimSpace(pair[1])

	return key, value
}

func NewConfigValue(value string) *ConfigValue {
	return &ConfigValue{value: value}
}

func (v ConfigValue) AsInt(defv int64) int64 {
	ret, err := strconv.ParseInt(v.value, 0, 64)
	if err != nil {
		Warn("config parse int error: %s", v.value)
		return defv
	}
	return ret
}

func (v ConfigValue) AsFloat(defv float64) float64 {
	ret, err := strconv.ParseFloat(v.value, 64)
	if err != nil {
		Warn("config parse float error: %s", v.value)
		return defv
	}
	return ret
}

func (v ConfigValue) AsBool(defv bool) bool {
	ret, err := strconv.ParseBool(v.value)
	if err != nil {
		Warn("config parse bool error: %s", v.value)
		return defv
	}
	return ret
}

func (v ConfigValue) AsString() string {
	//Warn("ConfigValue: %p\n", &v)
	return v.value
}

func (v ConfigValue) AsStrArray() []string {
	s := strings.Split(v.value, ",")

	for i := 0; i < len(s); i++ {
		s[i] = strings.TrimSpace(s[i])
	}
	return s
}

func (v ConfigValue) AsIntArray(defv int64) []int64 {
	s := strings.Split(v.value, ",")

	var ret []int64
	//var err error
	for i := 0; i < len(s); i++ {
		a, err := strconv.ParseInt(strings.TrimSpace(s[i]), 0, 64)
		if err != nil {
			Warn("config parse int error: %s", v.value)
			ret = append(ret, defv)
		} else {
			ret = append(ret, a)
		}
	}
	return ret
}

func (v ConfigValue) AsFloatArray(defv float64) []float64 {
	s := strings.Split(v.value, ",")

	var ret []float64
	//var err error
	for i := 0; i < len(s); i++ {
		a, err := strconv.ParseFloat(strings.TrimSpace(s[i]), 64)
		if err != nil {
			ret = append(ret, defv)
		} else {
			ret = append(ret, a)
		}
	}
	return ret
}

func (v ConfigValue) AsAnyArray(restr string) []string {
	s := strings.Split(v.value, ",")

	re := regexp.MustCompile(restr)
	//groups := re.SubexpNames()
	var ret []string
	for i := 0; i < len(s); i++ {
		match := re.FindStringSubmatch(s[i])
		if len(match) == 0 {
			Warn("config parser value error: %s", v.value)
			return nil
		} else {
			for j := 1; j < len(match)-1; j++ {
				ret = append(ret, match[j])
			}
		}
	}
	return ret
}
