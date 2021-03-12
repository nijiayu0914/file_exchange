package utils

// RequestCheckUuid   UUid检查请求实体
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

// RequestUpdateUsage 更新网盘用量请求实体，定义了文件uuid, 更新用量值和更新方式
type RequestUpdateUsage struct {
	FileUuid string `json:"file_uuid"` // 文件uuid
	UsageCapacity float64 `json:"usage_capacity"` //更新用量值
	How string `json:"how"` //更新方式，increase, decrease, overwrite
}

// RequestCapacity 更新网盘允许用量请求实体，定义了文件uuid, 允许用量
type RequestCapacity struct {
	FileUuid string `json:"file_uuid"` // 文件uuid
	Capacity float64 `json:"capacity"` // 允许用量
}

type RequestListFiles struct {
	FileUuid string `json:"file_uuid"` // 文件uuid
	Path     string `json:"path"`
	Delimiter string `json:"delimiter"`
}

type RequestDeleteHistoryFile struct {
	FileUuid string `json:"file_uuid"` // 文件uuid
	Path     string `json:"path"`
	VersionId string `json:"version_id"`
}

type RequestRenameObject struct {
	ObjectName string `json:"object_name"`
	NewName string `json:"new_name"`
}

type RequestDeleteFiles struct {
	FileUuid string `json:"file_uuid"` // 文件uuid
	FileNames []string `json:"file_names"`
}

type RequestCopy struct {
	OriginFile string `json:"origin_file"`
	DestFile string `json:"dest_file"`
	VersionId string `json:"versionid"`
}

type RequestMultipleCopy struct {
	CopyList []RequestCopy `json:"copy_list"`
}

type RequestReadFileSize struct {
	FileUuid string `json:"file_uuid"` // 文件uuid
	FileName string `json:"file_name"`
	VersionId string `json:"version_id"`
}

type RequestReadFilesSize struct {
	FileUuid string `json:"file_uuid"` // 文件uuid
	Files []RequestReadFileSize `json:"files"`
}