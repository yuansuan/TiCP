package linker

type Linker interface {
	// Link 链接文件
	Link(sourcePath string, destPath string) error
}
