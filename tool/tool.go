package tool

import (
	"fmt"
	"strings"

	"gopkg.in/ini.v1"
)

func ReadIni(sec string, key string) string {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		fmt.Println("Failed to read config.ini")
		return ""
	}
	section := cfg.Section(sec)
	if section != nil {
		return section.Key(key).String()
	} else {
		fmt.Printf("No %s section", sec)
		return ""
	}
}

// 获取元素在切片的位置 默认返回0
func GetKeyIndex(arr []string, key string) int {
	if key == "" || len(arr) == 0 {
		return 0
	}
	for i, v := range arr {
		if v == key {
			return i
		}
	}
	return 0
}

func GetLastThreeLines(str string) string {
	lines := strings.Split(str, "\n")
	if len(lines) <= 4 {
		return str
	}
	lastThreeLines := strings.Join(lines[len(lines)-3:], "\n")
	return lastThreeLines
}
