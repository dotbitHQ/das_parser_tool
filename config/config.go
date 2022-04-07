package config

import (
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Cfg CfgServer
var Env core.Env

func InitCfg(cfgFile string) {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("./config")
		viper.AddConfigPath("../config")
	}
	if err := viper.ReadInConfig(); err != nil {
		cobra.CheckErr(fmt.Errorf("ReadInConfig err: %v", err.Error()))
	}
	if err := viper.Unmarshal(&Cfg); err != nil {
		cobra.CheckErr(fmt.Errorf("Unmarshal err: %v ", err.Error()))
	}
}

type CfgServer struct {
	Chain struct {
		Net      common.DasNetType `mapstructure:"net"`
		CkbUrl   string            `mapstructure:"ckb_url"`
		IndexUrl string            `mapstructure:"index_url"`
	} `mapstructure:"chain"`
}
