/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"log"
	"os"

	"github.com/yuansuan/ticp/common/openapi-go/tools/mysql2mongo/internal/mongo"
	"github.com/yuansuan/ticp/common/openapi-go/tools/mysql2mongo/internal/mysql"

	"github.com/ory/viper"
	"github.com/spf13/cobra"
)

var _config Config

type Config struct {
	Mysql *mysql.Config `yaml:"mysql"`
	Mongo *mongo.Config `yaml:"mongo"`
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mysql2mongo",
	Short: "Data migration from Mysql to MongoDB",
	Long:  `Data migration from Mysql to MongoDB, 目前只支持残差图数据的迁移`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		residuals := mysql.FindResidualsFromMysql(ctx, _config.Mysql.Uri)
		mongo.InsertResidualsToMongoDB(ctx, _config.Mongo.URI(), _config.Mongo.Database, residuals)
		log.Println("Data migration from Mysql to MongoDB done")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	initConfig()
}

func initConfig() {
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	cobra.CheckErr(viper.Unmarshal(&_config))

	// log.Println(spew.Sdump(_config))
}
