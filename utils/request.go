package utils

type RequestChangePassword struct {
	UserName string `json:"user_name"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type RequestResetPassword struct {
	UserName string `json:"user_name"`
	NewPassword string `json:"new_password"`
}

type RequestUpdateUsage struct {
	FileUuid string `json:"file_uuid"`
	UsageCapacity float64 `json:"usage_capacity"`
	How string `json:"how"`
}

type RequestCapacity struct {
	FileUuid string `json:"file_uuid"`
	Capacity float64 `json:"capacity"`
}

type RequestListFiles struct {
	FileUuid string `json:"file_uuid"`
	Path     string `json:"path"`
	Delimiter string `json:"delimiter"`
}

type RequestDeleteHistoryFile struct {
	FileUuid string `json:"file_uuid"`
	Path     string `json:"path"`
	VersionId string `json:"version_id"`
}

type RequestRenameObject struct {
	ObjectName string `json:"object_name"`
	NewName string `json:"new_name"`
}

type RequestDeleteFiles struct {
	FileUuid string `json:"file_uuid"`
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
	FileName string `json:"file_name"`
	VersionId string `json:"version_id"`
}

type RequestReadFilesSize struct {
	Files []RequestReadFileSize `json:"files"`
}