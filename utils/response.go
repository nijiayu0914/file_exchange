package utils

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"net/http"
)

// Response 标准响应实体的数据结构，定义了响应码，响应消息，响应数据
type Response struct {
	Code    int `json:"code"` // http code
	Message string `json:"message"` // 响应消息
	Data    interface{} `json:"data"` // 响应数据
}

// ObjectInfoCollection OSS文件对象信息对象映射，定义了OSS基本信息与OSS元信息
type ObjectInfoCollection struct {
	Basic oss.ObjectProperties `json:"basic"` // OSS对象基本信息
	Meta  http.Header `json:"meta"` // OSS对象元信息
}

// DirInfoCollection OSS模拟文件夹对象信息对象映射，定义了文件夹名称和OSS元信息
type DirInfoCollection struct {
	Basic string `json:"basic"` // 文件夹名称
	Meta  http.Header `json:"meta"` // OSS对象元信息
}

// ResponseListFiles 列举文件夹对象映射，定义了文件数量，文件夹数量，文件数组，文件夹数组
type ResponseListFiles struct {
	FilesCount int `json:"file_count"` // 文件数量
	DirsCount int `json:"dis_count"` // 文件夹数量
	Files []ObjectInfoCollection `json:"files"` // 文件数组
	Dirs  []DirInfoCollection    `json:"dirs"` // 文件夹数组
}
