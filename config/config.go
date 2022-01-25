package config

import (
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Cfg CfgServer

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
	DasCore struct {
		THQCodeHash         string                            `mapstructure:"thq_code_hash"`
		DasContractArgs     string                            `mapstructure:"das_contract_args"`
		DasContractCodeHash string                            `mapstructure:"das_contract_code_hash"`
		DasConfigCodeHash   string                            `mapstructure:"das_config_code_hash"`
		MapDasContract      map[common.DasContractName]string `mapstructure:"map_das_contract"`
		CellDeps            map[string]string                 `mapstructure:"cell_deps"`
	} `mapstructure:"das_core"`
}
