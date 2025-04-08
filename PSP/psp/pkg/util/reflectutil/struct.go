package reflectutil

import "reflect"

func GetStructName(v interface{}) string {
	vType := reflect.TypeOf(v)

	if vType.Kind() == reflect.Ptr {
		vType = vType.Elem()
	}
	if vType.Kind() != reflect.Struct {
		return ""
	}

	return vType.Name()
}
