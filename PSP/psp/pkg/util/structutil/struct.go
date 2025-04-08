package structutil

import "github.com/jinzhu/copier"

// CopyStruct 拷贝结构体
func CopyStruct(dst interface{}, src interface{}) error {
	if err := copier.Copy(dst, src); err != nil {
		return err
	}
	return nil
}
