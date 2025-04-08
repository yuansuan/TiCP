package config

import (
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/mitchellh/mapstructure"
	"github.com/ory/viper"
	"github.com/stretchr/testify/assert"
)

const cfgStr = `
change_license: true

zones:
  az-jinan:
    hpc_endpoint: https://jn_hpc_endpoint:8080
    storage_endpoint: https://jn_storage_endpoint:8899
    cloud_app_enable: true
  az-wuxi:
    hpc_endpoint: https://wx_hpc_endpoint:8080
    storage_endpoint: https://wx_storage_endpoint:8899
    cloud_app_enable: false

self_ys_id: 4VKvohbmSjC
ak: 123456
as: 999+1123
`

func initConfig() CustomT {
	viper.SetConfigType("yaml")
	err := viper.ReadConfig(strings.NewReader(cfgStr))
	if err != nil {
		panic(err)
	}

	md := mapstructure.Metadata{}
	customT := CustomT{}
	err = viper.Unmarshal(&customT, func(config *mapstructure.DecoderConfig) {
		config.TagName = "yaml"
		config.Metadata = &md
	})
	if err != nil {
		panic(err)
	}

	return customT
}

func Test_initConfig(t *testing.T) {
	customT := initConfig()
	t.Logf("customT: %+v", customT)
}

func TestGetZones(t *testing.T) {
	customT := initConfig()
	zones := customT.Zones

	spew.Dump(zones)

	jnhpc := "https://jn_hpc_endpoint:8080"
	jnstorage := "https://jn_storage_endpoint:8899"

	zone := zones.GetZoneByEndpoint(jnhpc)
	zone2 := zones.GetZoneByEndpoint(jnstorage)

	t.Logf("zone: %s, zone2: %s", zone, zone2)
	assert.Equal(t, "az-jinan", zone)
	assert.Equal(t, "az-jinan", zone2)

	wuxihpc := "https://wx_hpc_endpoint:8080"
	wuxistorage := "https://wx_storage_endpoint:8899"

	zone3 := zones.GetZoneByEndpoint(wuxihpc)
	zone4 := zones.GetZoneByEndpoint(wuxistorage)

	t.Logf("zone3: %s, zone4: %s", zone3, zone4)
	assert.Equal(t, "az-wuxi", zone3)
	assert.Equal(t, "az-wuxi", zone4)
}
