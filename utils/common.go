package utils

import (
	"os"
	"strings"
)

func StringJoin(sep string, strs ...string) string{
	var strList []string
	for _, str := range strs {
		strList = append(strList, str)
	}
	return strings.Join(strList, sep)
}

func FileExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}