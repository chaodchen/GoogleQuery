package tool
import (
	"fmt"

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