package utils

// RequestCheckUuid   UUid检查请求实体，定义了文件uuid
type RequestCheckUuid struct {
	FileUuid string `json:"file_uuid"` // 文件uuid
}

// RequestChangePassword   修改用户密码请求实体，定义了用户名，原始密码和新密码
type RequestChangePassword struct {
	UserName string `json:"user_name"` // 用户名
	OldPassword string `json:"old_password"` // 旧密码
	NewPassword string `json:"new_password"` // 新密码
}

// RequestResetPassword   重置用户密码请求实体，定义了用户名和新密码
type RequestResetPassword struct {
	UserName string `json:"user_name"` // 用户名
	NewPassword string `json:"new_password"` // 新密码
}

// RequestUpdateUsage 更新网盘用量请求实体，定义了文件uuid，更新用量值和更新方式
type RequestUpdateUsage struct {
	FileUuid string `json:"file_uuid"` // 文件uuid
	UsageCapacity float64 `json:"usage_capacity"` //更新用量值
	How string `json:"how"` //更新方式，increase, decrease, overwrite
}

// RequestCapacity 更新网盘允许用量请求实体，定义了文件uuid，允许用量
type RequestCapacity struct {
	FileUuid string `json:"file_uuid"` // 文件uuid
	Capacity float64 `json:"capacity"` // 允许用量
}

// RequestListFiles 列举文件请求实体，定义了文件uuid，文件路径，路径终止符号
type RequestListFiles struct {
	FileUuid string `json:"file_uuid"` // 文件uuid
	Path     string `json:"path"` // 文件路径
	Delimiter string `json:"delimiter"` // OSS路径终止符
}

// RequestDeleteHistoryFile 删除历史文件请求实体，定义了文件uuid，文件路径，文件版本号
type RequestDeleteHistoryFile struct {
	FileUuid string `json:"file_uuid"` // 文件uuid
	Path     string `json:"path"` // 文件路径
	VersionId string `json:"version_id"` // 文件版本号
}

// RequestRenameObject 文件重命名请求实体，仅更改文件元信息，
// 定义了OSS文件对象名称，重命名名称
type RequestRenameObject struct {
	FileUuid string `json:"file_uuid"` // 文件uuid
	ObjectName string `json:"object_name"` // OSS对象文件名
	NewName string `json:"new_name"` // 重命名名称
}

// RequestDeleteFiles 删除多个文件请求实体，定义了文件uuid，删除文件的文件名数组
type RequestDeleteFiles struct {
	FileUuid string `json:"file_uuid"` // 文件uuid
	FileNames []string `json:"file_names"` // 删除文件的文件名数组
}

// RequestCopy 拷贝文件请求实体，定义了源文件名称，目标文件名称，文件版本号
type RequestCopy struct {
	FileUuid string `json:"file_uuid"` // 文件uuid
	OriginFile string `json:"origin_file"` // 源文件名称
	DestFile string `json:"dest_file"` // 目标文件名称
	VersionId string `json:"versionid"` // 文件版本号
}

// RequestMultipleCopy 拷贝多个文件请求实体，定义了拷贝对象数组
type RequestMultipleCopy struct {
	FileUuid string `json:"file_uuid"` // 文件uuid
	CopyList []RequestCopy `json:"copy_list"` // RequestCopy 数组
}

// RequestReadFileSize 查询文件容量请求实体，定义了文件uuid，文件名，文件版本号
type RequestReadFileSize struct {
	FileUuid string `json:"file_uuid"` // 文件uuid
	FileName string `json:"file_name"` // 文件名
	VersionId string `json:"version_id"` // 文件版本号
}

// RequestReadFilesSize 查询多个文件容量请求实体，
// 定义了文件uuid，RequestReadFileSize 数组
type RequestReadFilesSize struct {
	FileUuid string `json:"file_uuid"` // 文件uuid
	Files []RequestReadFileSize `json:"files"` // RequestReadFileSize 数组
}