package tormentagui

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// prettyPrintStruct just outputs a struct as indented JSON
func printJSONWithoutModel(s interface{}) string {
	res, _ := json.MarshalIndent(s, "", " ")
	return trimModel(string(res))
}

func printJSON(s interface{}) string {
	res, _ := json.MarshalIndent(s, "", " ")
	return string(res)
}

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

// trimModel takes the marshalled tormentable JSON string,
// and removes the embedded Model (ID, Created, LastUpdated)
// This is useful to show just editable fields
func trimModel(s string) string {
	withoutParens := strings.Trim(s, "{}")
	splitOnCommas := strings.Split(withoutParens, ",")
	rejoinedWithoutModel := strings.Join(splitOnCommas[3:], ",")
	return fmt.Sprintf("{%s}", rejoinedWithoutModel)
}
