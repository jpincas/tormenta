package tormentarest

import "strings"

func entityNameFromPath_List(path string) string {
	return strings.TrimPrefix(path, "/"+App.Root+"/")
}

func entityNameFromPath_ID(path string) string {
	return strings.Split(path, "/")[2]
}
