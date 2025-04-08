package zone

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	t.Run("correct input", func(t *testing.T) {
		assert.Equal(t, AZShangHai, Parse("az-shanghai"))
		assert.Equal(t, AZShanXi, Parse("az-shanxi"))
		assert.Equal(t, AZShenZhen, Parse("az-shenzhen"))
		assert.Equal(t, AZGuangZhou, Parse("az-guangzhou"))
		assert.Equal(t, AZJiNan, Parse("az-jinan"))
		assert.Equal(t, AZWuXi, Parse("az-wuxi"))
		assert.Equal(t, AZTianJin, Parse("az-tianjin"))
	})

	t.Run("bad", func(t *testing.T) {
		t.Run("empty input", func(t *testing.T) {
			assert.Equal(t, AZEmpty, Parse(""))
		})

		t.Run("invalid input", func(t *testing.T) {
			assert.Equal(t, AZInvalid, Parse("123"))
		})
		assert.Equal(t, AZYuansuan, Parse("az-yuansuan"), "parse zone error")
	})
}

func TestZone_IsValid(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		assert.Equal(t, true, AZShangHai.IsValid())
		assert.Equal(t, true, AZGuangZhou.IsValid())
		assert.Equal(t, true, AZJiNan.IsValid())
		assert.Equal(t, true, AZShanXi.IsValid())
		assert.Equal(t, true, AZTianJin.IsValid())
		assert.Equal(t, true, AZWuXi.IsValid())
		assert.Equal(t, true, AZShenZhen.IsValid())
		assert.Equal(t, true, AZEmpty.IsValid())
	})

	t.Run("invalid", func(t *testing.T) {
		assert.Equal(t, false, AZInvalid.IsValid())
	})
}

func TestZone_IsEmpty(t *testing.T) {
	t.Run("nonempty", func(t *testing.T) {
		assert.Equal(t, false, AZShangHai.IsEmpty())
		assert.Equal(t, false, AZGuangZhou.IsEmpty())
		assert.Equal(t, false, AZJiNan.IsEmpty())
		assert.Equal(t, false, AZShanXi.IsEmpty())
		assert.Equal(t, false, AZTianJin.IsEmpty())
		assert.Equal(t, false, AZWuXi.IsEmpty())
		assert.Equal(t, false, AZShenZhen.IsEmpty())
		assert.Equal(t, false, AZInvalid.IsEmpty())
	})

	t.Run("empty", func(t *testing.T) {
		assert.Equal(t, true, AZEmpty.IsEmpty())
	})
}

func TestZone_String(t *testing.T) {
	assert.Equal(t, "az-shanghai", AZShangHai.String())
	assert.Equal(t, "az-wuxi", AZWuXi.String())
	assert.Equal(t, "az-jinan", AZJiNan.String())
	assert.Equal(t, "az-guangzhou", AZGuangZhou.String())
	assert.Equal(t, "az-shenzhen", AZShenZhen.String())
	assert.Equal(t, "az-tianjin", AZTianJin.String())
	assert.Equal(t, "az-shanxi", AZShanXi.String())
	assert.Equal(t, "az-invalid", AZInvalid.String())
	assert.Equal(t, "az-empty", AZEmpty.String())
}

func TestZone_Desc(t *testing.T) {
	t.Run("correct", func(t *testing.T) {
		assert.Equal(t, "华东1", AZShangHai.Desc())
		assert.Equal(t, "华东2", AZWuXi.Desc())
		assert.Equal(t, "华东3", AZJiNan.Desc())
		assert.Equal(t, "华南1", AZGuangZhou.Desc())
		assert.Equal(t, "华南2", AZShenZhen.Desc())
		assert.Equal(t, "华北1", AZTianJin.Desc())
		assert.Equal(t, "华北2", AZShanXi.Desc())
	})

	t.Run("bad", func(t *testing.T) {
		assert.Equal(t, "UNKNOWN", Zone("123").Desc())
	})
}

func TestParseWithDefault(t *testing.T) {
	t.Run("correct", func(t *testing.T) {
		assert.Equal(t, AZShangHai, ParseWithDefault("az-shanghai", AZShanXi))
		assert.Equal(t, AZTianJin, ParseWithDefault("az-tianjin", AZShangHai))
		assert.Equal(t, AZShenZhen, ParseWithDefault("az-shenzhen", AZWuXi))
	})

	t.Run("bad", func(t *testing.T) {
		t.Run("invalid", func(t *testing.T) {
			assert.Equal(t, AZShangHai, ParseWithDefault("", AZShangHai))
		})

		t.Run("empty", func(t *testing.T) {
			assert.Equal(t, AZShangHai, ParseWithDefault("123", AZShangHai))
		})
	})
}

func TestMustParse(t *testing.T) {
	t.Run("correct", func(t *testing.T) {
		assert.Equal(t, AZShenZhen, MustParse("az-shenzhen"))
		assert.Equal(t, AZShanXi, MustParse("az-shanxi"))
		assert.Equal(t, AZWuXi, MustParse("az-wuxi"))
		assert.Equal(t, AZTianJin, MustParse("az-tianjin"))
		assert.Equal(t, AZShangHai, MustParse("az-shanghai"))
		assert.Equal(t, AZGuangZhou, MustParse("az-guangzhou"))
		assert.Equal(t, AZJiNan, MustParse("az-jinan"))
		assert.Equal(t, AZEmpty, MustParse("az-empty"))
	})

	t.Run("bad", func(t *testing.T) {
		assert.PanicsWithValue(t, "invalid zone", func() {
			MustParse("123")
		})
	})
}
