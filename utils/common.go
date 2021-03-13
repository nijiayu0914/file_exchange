// utils 开发中各类常用工具整合
package utils

import (
	"crypto/md5"
	"fmt"
	"os"
	"strings"
)

// StringJoin 用于字符串连接
//    StringJoin("/", "a", "b")
//    return "a/b"
func StringJoin(sep string, strs ...string) string{
	var strList []string
	for _, str := range strs {
		strList = append(strList, str)
	}
	return strings.Join(strList, sep)
}

// FileExist 判断目录下文件是否存在
//    FileExist("./config.yml")
//    return true
func FileExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// Str2md5 字符串进行md5转换
func Str2md5(str string) (md5Str string){
	data:=[]byte(str)
	h := md5.New()
	h.Write(data)
	md5Str = fmt.Sprintf("%x", h.Sum(nil))
	return md5Str
}