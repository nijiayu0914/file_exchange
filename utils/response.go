package utils

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"net/http"
)

type Response struct {
	Code    int
	Message string
	Data    interface{}
}

type ObjectInfoCollection struct {
	Basic oss.ObjectProperties
	Meta  http.Header
}

type DirInfoCollection struct {
	Basic string
	Meta  http.Header
}

type ResponseListFiles struct {
	FilesCount int
	DirsCount int
	Files []ObjectInfoCollection `json:"files"`
	Dirs  []DirInfoCollection    `json:"dirs"`
}