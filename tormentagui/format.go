package tormentagui

import (
	"fmt"
	"time"
)

func formatTime(t time.Time) string {
	return t.Format("2006, 2 Jan 15:04:05")
}

func autoFormat(i interface{}) string {

	switch i.(type) {
	case time.Time:
		return formatTime(i.(time.Time))
	case bool:
		v := i.(bool)
		if v {
			return "T"
		} else {
			return "F"
		}
	default:
		return fmt.Sprint(i)
	}

}
