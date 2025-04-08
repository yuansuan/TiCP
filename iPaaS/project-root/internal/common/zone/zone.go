package zone

// Zone 表示一个可用区
type Zone string

// IsValid 判断可用区是否有效
func (z Zone) IsValid() bool {
	return z != AZInvalid
}

// IsEmpty 是否为空可用区
func (z Zone) IsEmpty() bool {
	return z == AZEmpty
}

// String 工具函数，支持语言的默认行为
func (z Zone) String() string {
	return string(z)
}

// Desc 返回可用区的描述
func (z Zone) Desc() string {
	switch z {
	case AZYuansuan:
		return "华东1"
	case AZShangHai:
		return "华东1"
	case AZWuXi:
		return "华东2"
	case AZJiNan:
		return "华东3"
	case AZGuangZhou:
		return "华南1"
	case AZShenZhen:
		return "华南2"
	case AZTianJin:
		return "华北1"
	case AZShanXi:
		return "华北2"
	case AZZhigu:
		return "杭州智谷大厦"
	default:
		return "UNKNOWN"
	}
}

const (
	// AZYuansuan 远算云
	AZYuansuan Zone = "az-yuansuan"

	// AZShangHai 上海
	AZShangHai Zone = "az-shanghai"
	// AZWuXi 无锡
	AZWuXi Zone = "az-wuxi"
	// AZJiNan 济南
	AZJiNan Zone = "az-jinan"
	// AZGuangZhou 广州
	AZGuangZhou Zone = "az-guangzhou"
	// AZShenZhen 深圳
	AZShenZhen Zone = "az-shenzhen"
	// AZTianJin 天津
	AZTianJin Zone = "az-tianjin"
	// AZShanXi 山西
	AZShanXi Zone = "az-shanxi"
	// AZZhigu 杭州智谷大厦
	AZZhigu Zone = "az-zhigu"

	// AZGanSu 甘肃 注：目前只有开发环境
	AZGanSu Zone = "az-gansu"

	// AZInvalid 无效的区域
	AZInvalid Zone = "az-invalid"
	// AZEmpty 空区域
	AZEmpty Zone = "az-empty"
)

// Parse 解析一个可用区配置
func Parse(s string) Zone {
	switch Zone(s) {
	case AZYuansuan:
		return AZYuansuan
	case AZShangHai:
		return AZShangHai
	case AZWuXi:
		return AZWuXi
	case AZJiNan:
		return AZJiNan
	case AZGuangZhou:
		return AZGuangZhou
	case AZShenZhen:
		return AZShenZhen
	case AZTianJin:
		return AZTianJin
	case AZShanXi:
		return AZShanXi
	case AZGanSu:
		return AZGanSu
	case AZZhigu:
		return AZZhigu
	case "", AZEmpty:
		return AZEmpty
	default:
		return AZInvalid
	}
}

// ParseWithDefault 解析一个可用区配置 如果为空或者无效的 则使用默认的zone
func ParseWithDefault(s string, defaultZone Zone) Zone {
	z := Parse(s)
	switch z {
	case AZEmpty, AZInvalid:
		return defaultZone
	default:
		return z
	}
}

// MustParse 解析并检查可用区是否有效
func MustParse(s string) Zone {
	z := Parse(s)
	if !z.IsValid() {
		panic("invalid zone")
	}
	return z
}
