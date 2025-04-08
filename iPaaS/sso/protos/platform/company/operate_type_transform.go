package company

// OpTypeToString 操作类型文件说明
func (x OperateType) OpTypeToString() string {
	switch x {
	case OperateType_UPLOAD:
		return "上传"
	case OperateType_DOWNLOAD:
		return "下载"
	case OperateType_DELETE:
		return "删除"
	case OperateType_RENAME:
		return "重命名"
	case OperateType_ADD_FOLDER:
		return "添加文件夹"
	default:
		return "unknown"
	}
	return ""
}

// FileTypeToString 文件类型说明
func (x FileType) FileTypeToString() string {
	switch x {
	case FileType_FILE:
		return "文件"
	case FileType_FOLDER:
		return "文件夹"
	case FileType_BATCH:
		return "批量操作"
	default:
		return "unknown"
	}
	return ""
}
