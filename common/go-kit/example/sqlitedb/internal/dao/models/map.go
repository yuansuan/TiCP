package models

type KeyMap struct {
	ID    *int64 `xorm:"'rowid' pk"`
	Key   string
	Value string
}

func (m *KeyMap) TableName() string {
	return "key_map"
}
