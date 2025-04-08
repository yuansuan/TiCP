package _type

import "fmt"

type Type int

const (
	Default   = Type(0) // 默认值
	Internal  = Type(1) // 自有的
	External  = Type(2) // 外部的
	Consigned = Type(3) // 寄售的
)

type UnsupportedLicenseTypeError struct {
	Value int
}

func (s UnsupportedLicenseTypeError) Error() string {
	return fmt.Sprintf("Unsupported license type: %d", s.Value)
}

func ToType(value int) (Type, error) {
	switch Type(value) {
	case Default:
		return Default, nil
	case Internal:
		return Internal, nil
	case External:
		return External, nil
	case Consigned:
		return Consigned, nil
	default:
		return -1, UnsupportedLicenseTypeError{Value: value}
	}
}

func (t *Type) GetValue() int {
	return int(*t)
}

func (t *Type) IsSelfLic() bool {
	return *t == Internal
}

func (t *Type) IsConsignedLic() bool {
	return *t == Consigned
}

func (t *Type) IsOthersLic() bool {
	return *t == External
}
