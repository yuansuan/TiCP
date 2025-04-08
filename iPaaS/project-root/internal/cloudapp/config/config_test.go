package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v2"
)

type Suite struct {
	suite.Suite

	mockConfig CustomT
}

func (s *Suite) SetupSuite() {

	data, err := os.ReadFile("dev_custom.yml")

	if err != nil {
		panic(err)
	}

	s.mockConfig = CustomT{}

	err = yaml.Unmarshal(data, &s.mockConfig)

	if err != nil {
		panic(err)
	}

	custom = s.mockConfig
}

func TestConfig(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestParse() {
	s.Equal(AZInvalid, Parse("az-invalid"))
	s.Equal(AZEmpty, Parse("az-empty"))
	s.Equal(AZEmpty, Parse(""))
	s.Equal(AZEmpty, Parse("az-empty"))

	s.Equal(Zone("az-yuansuan"), Parse("az-yuansuan"))
}
